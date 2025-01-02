package server

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/bmehdi777/moon/internal/pkg/communication"
	"github.com/bmehdi777/moon/internal/pkg/server/config"
	"github.com/bmehdi777/moon/internal/pkg/server/database"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func tcpServe(channelsDomains *ChannelsDomains, db *gorm.DB) {
	cert, err := tls.LoadX509KeyPair(config.GlobalConfig.App.CertPemPath, config.GlobalConfig.App.CertKeyPath)
	if err != nil {
		log.Fatalf("Can't load TLS certificates : %v", err)
		return
	}
	configTls := tls.Config{Certificates: []tls.Certificate{cert}}
	configTls.Rand = rand.Reader

	fullAddrFmt := fmt.Sprintf("%v:%v", config.GlobalConfig.App.TcpAddr, config.GlobalConfig.App.TcpPort)
	listener, err := tls.Listen("tcp", fullAddrFmt, &configTls)
	if err != nil {
		log.Fatalf("Can't setup port :  %v", err)
		return
	}
	defer listener.Close()
	log.Printf("TCP Server is up at %v", fullAddrFmt)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Error while accepting a connection : %v", err)
			continue
		}

		if tcpConn, ok := conn.(*net.TCPConn); ok {
			tcpConn.SetKeepAlive(true)
			tcpConn.SetKeepAlivePeriod(1 * time.Hour)
		}

		client := communication.NewClient(conn.(*tls.Conn), nil)
		go handleClient(client, channelsDomains, db)
	}
}

func handleClient(client *communication.Client, channelsDomains *ChannelsDomains, db *gorm.DB) {
	// defer aren't working, search why
	defer client.Connection.Close()
	defer log.Printf("Connection closed with %v", client.Connection.RemoteAddr())
	// create random domain name
	channelsName, err := createOrSelectChannelForUser(client, channelsDomains, db)
	if err != nil {
		// maybe shouldn't crash but skip this connection ?
		log.Fatalf("Error while creating channels : %v", err)
	}

	channels := channelsDomains.Get(channelsName)
	if channels == nil {
		log.Fatalf("Error while retrieving channel.")
	}

	defer channelsDomains.Delete(channelsName)
	defer db.Model(&database.DomainRecord{}).Where("dns_record = ?", channelsName).Update("connection_open", false)

	log.Printf("Connection started with %v", client.Connection.RemoteAddr())
	var reply *http.Request
	readChan := make(chan *communication.Packet)

	go func() {
		for {
			fmt.Println("called here twice ?")
			responsePacket, err := client.Read()
			if err != nil {
				if err != io.EOF {
					log.Fatalf("Error while reading response from %v : %v", client.Connection.RemoteAddr(), err)
				}
				return
			}
			fmt.Println("responsePacket ", responsePacket)
			readChan <- responsePacket
		}
	}()
	for {

		select {
		case reply = <-channels.RequestChannel:
			var buf bytes.Buffer
			err := reply.Write(&buf)
			if err != nil {
				log.Fatalf("Error while writing request to wire : %v", err)
				return
			}
			bufBytes := buf.Bytes()

			// redirecting HTTP request to TCP connection
			err = client.SendHttpRequest(bufBytes)
			if err != nil {
				log.Fatalf("Error while sending bytes to %v : %v", client.Connection.RemoteAddr(), err)
				return
			}
		case response := <-readChan:

			// move in a function
			switch response.Header.Type {
			case communication.ConnectionClose:
				fmt.Println("closing conn")
				// it will automically conn.close
				return
			case communication.HttpResponse:
				reader := bytes.NewReader(response.Payload.Data)
				respBufio := bufio.NewReader(reader)
				resp, err := http.ReadResponse(respBufio, reply)
				if err != nil {
					log.Fatalf("Error while converting bytes to HTTP response %v", err)
					return
				}
				channels.ResponseChannel <- resp
			default:
				log.Fatalf("Weird packet received. Skipping it.")
			}
		}
	}
}

func createOrSelectChannelForUser(client *communication.Client, channels *ChannelsDomains, db *gorm.DB) (string, error) {
	// received auth message
	packet, err := client.Read()
	if err != nil {
		return "", err
	}

	if packet.Header.Type != communication.ConnectionStart {
		return "", fmt.Errorf("Can't start a connection")
	}

	log.Printf("Packet received : %#v", packet)

	accessToken, err := verifyJwt(packet.Payload.Token)
	if err != nil {
		// TODO: currently segfault bc connection is close
		err = client.SendInvalidToken()
		if err != nil {
			return "", err
		}
	}
	sub, err := accessToken.Claims.GetSubject()
	if err != nil {
		return "", err
	}
	fmt.Println("sub: ", sub)

	// create domain record in db
	var user database.User
	db.First(&user, "kc_user_id = ? ", sub)

	// no domain record registered
	var dnsRecord string
	if user.DomainRecordID == 0 {
		dnsRecord = uuid.NewString() + "." + config.GlobalConfig.App.GlobalDomainName
		record := database.DomainRecord{
			DNSRecord:      dnsRecord,
			ConnectionOpen: true,
		}
		db.Model(&user).Update("DomainRecord", record)
	} else {
		// refacto one day
		var domainRecord database.DomainRecord
		db.Table("domain_records").Select("domain_records.dns_record").Joins("INNER JOIN users ON domain_records.ID = users.domain_record_id").Scan(&domainRecord)
		dnsRecord = domainRecord.DNSRecord
		log.Printf("dns record : %v", user.DomainRecord.DNSRecord)
		db.Model(&database.DomainRecord{}).Where("dns_record = ?", dnsRecord).Update("connection_open", true)
	}

	log.Printf("Connection open : http://%v", dnsRecord)

	// create channel
	channels.Add(dnsRecord)

	return dnsRecord, nil
}

func verifyJwt(tokenStr string) (*jwt.Token, error) {
	// local cert for testing purpose
	spkiPem := `
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0GGJjxtCXGQgKxwcFZwd2AkNdaPSMN76A5bJyk6Dve8gMi8sbypzKngzhkziqofVe9g5H9kWRyZNIVzKiK4OnFhTRvRtAXoeWj98EINRMmvmWGv5BKwGmfr7g/mVvr+viyROUrrPUWx6TslyVD7VxLFrSchLiAdV6pZdMrKD1tlSXNQ78N3Q2Nw/SmuYd07wBIbtDCTwG9XaCJFaw0jgbKs6wdpTSqkfTNnYE2ekOlI8nAtTwAthjJeIfuPuScG4wVvbTTMx+Hd3z4kU2ripynSOVOWioyWUw6uerJqt1sgclNdQkFwdXgCzcOmJYIt8cOvCm8jEkNPmL3jJMN/eVQIDAQAB
-----END PUBLIC KEY-----
	`

	spkiBlock, _ := pem.Decode([]byte(spkiPem))
	var spkiKey *rsa.PublicKey
	pubInterface, _ := x509.ParsePKIXPublicKey(spkiBlock.Bytes)
	spkiKey = pubInterface.(*rsa.PublicKey)

	token, err := jwt.Parse(tokenStr, func(tok *jwt.Token) (interface{}, error) {
		if _, ok := tok.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", tok.Header["alg"])
		}

		return spkiKey, nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}
