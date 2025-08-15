package server

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"moon/internal/pkg/communication"
	"moon/internal/pkg/server/authent"
	"moon/internal/pkg/server/config"
	"moon/internal/pkg/server/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrInvalidToken = fmt.Errorf("Token is not valid.")

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

		client := communication.NewClient(conn.(*tls.Conn))
		go handleClient(client, channelsDomains, db)
	}
}

func handleClient(client *communication.Client, channelsDomains *ChannelsDomains, db *gorm.DB) {
	defer client.Connection.Close()
	defer log.Printf("Connection closed with %v", client.Connection.RemoteAddr())

	// create random domain name
	channelsName, err := createOrSelectChannelForUser(client, channelsDomains, db)
	if err != nil {
		if errors.Is(err, ErrInvalidToken) {
			// drop connection
			log.Printf("Invalid token from %v", client.Connection.RemoteAddr())
			return
		}
		log.Printf("Error while creating channels : %v", err)
		return
	}

	channels := channelsDomains.Get(channelsName)
	if channels == nil {
		log.Fatalf("Error while retrieving channel.")
	}

	defer channelsDomains.Delete(channelsName)
	defer func() {
		// close connection
		record, _ := database.FindRecordByDomainFQDN(channelsName, db)
		record.ConnectionOpen = false
		db.Save(&record)
	}()

	log.Printf("Connection started with %v", client.Connection.RemoteAddr())
	var reply *http.Request
	readChan := make(chan *communication.Packet)

	go func() {
		for {
			responsePacket, err := client.Read()
			if err != nil {
				if err != io.EOF {
					log.Fatalf("Error while reading response from %v : %v", client.Connection.RemoteAddr(), err)
				}
				continue
			}
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
			switch response.Header.Type {
			case communication.ConnectionClose:
				return
			case communication.Ping:
				err = client.SendPong()
				if err != nil {
					log.Fatalf("Error while responding to ping %v", err)
					return
				}
			case communication.HttpResponse:
				reader := bytes.NewReader(response.Payload.Data)
				respBufio := bufio.NewReader(reader)
				resp, err := http.ReadResponse(respBufio, reply)
				if err != nil {
					log.Fatalf("Error while converting bytes to HTTP response %v", err)
					return
				}
				reply = nil
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

	// Ensure we have a auth message
	if packet.Header.Type != communication.ConnectionStart {
		return "", fmt.Errorf("Can't start a connection")
	}

	msg, err := packet.Message()
	if err != nil {
		return "", err
	}

	authMsg := msg.(*communication.AuthMessage)
	accessToken, err := authent.VerifyJwt(authMsg.Token)
	if err != nil {
		err = client.SendUnauthorized()
		if err != nil {
			return "", err
		}
		return "", ErrInvalidToken
	}

	// tell client he is connected
	client.SendAuthorized()

	sub, err := accessToken.Claims.GetSubject()
	if err != nil {
		return "", err
	}
	fmt.Println("sub: ", sub)

	user, res := database.FindUserByKCUID(sub, db)
	if res.RowsAffected == 0 {
		return "", fmt.Errorf("User doesn't exist")
	}

	record, res := database.FindRecordByUserID(user.ID, db)
	var fqdn string
	if res.RowsAffected == 0 {
		// record doesn't exist
		fqdn = uuid.NewString() + "." + config.GlobalConfig.App.GlobalDomainName
		domainName := database.DomainName{
			FQDN: fqdn,
		}
		db.Create(&domainName)

		newRecord := database.Record{
			UserID:         user.ID,
			User:           *user,
			DomainNameID:   domainName.ID,
			DomainName:     domainName,
			ConnectionOpen: true,
		}
		db.Create(&newRecord)
	} else {
		// record exist
		record.ConnectionOpen = true
		db.Save(&record)
		domainName, _ := database.FindDomainNameById(record.DomainNameID, db)
		fqdn = domainName.FQDN
	}

	log.Printf("Connection open : http://%v", fqdn)

	channels.Add(fqdn)

	return fqdn, nil
}
