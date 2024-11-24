package server

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/bmehdi777/moon/internal/pkg/messages"
	"github.com/bmehdi777/moon/internal/pkg/server/config"
	"github.com/bmehdi777/moon/internal/pkg/server/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func tcpServe(channelsDomains *ChannelsDomains, db *gorm.DB) {
	cert, err := tls.LoadX509KeyPair(config.GlobalConfig.CertPemPath, config.GlobalConfig.CertKeyPath)
	if err != nil {
		log.Fatalf("Can't load TLS certificates : %v", err)
		return
	}
	configTls := tls.Config{Certificates: []tls.Certificate{cert}}
	configTls.Rand = rand.Reader

	fullAddrFmt := fmt.Sprintf("%v:%v", config.GlobalConfig.TcpAddr, config.GlobalConfig.TcpPort)
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

		go handleClient(conn, channelsDomains, db)
	}
}

func handleClient(conn net.Conn, channelsDomains *ChannelsDomains, db *gorm.DB) {
	// defer aren't working, search why
	defer conn.Close()
	defer log.Printf("Connection closed with %v", conn.RemoteAddr())
	// create random domain name
	channelsName, err := createOrSelectChannelForUser(conn, channelsDomains, db)
	if err != nil {
		log.Fatalf("Error while creating channels : %v", err)
	}

	channels := channelsDomains.Get(channelsName)
	if channels == nil {
		log.Fatalf("Error while retrieving channel.")
	}

	defer channelsDomains.Delete(channelsName)
	defer db.Model(&database.DomainRecord{}).Where("dns_record = ?", channelsName).Update("connection_open", false)

	log.Printf("Connection started with %v", conn.RemoteAddr())
	respBytes := make([]byte, 1024)

	for {
		// getting request from HTTP thread
		reply := <-channels.RequestChannel

		var buf bytes.Buffer
		err := reply.Write(&buf)
		if err != nil {
			log.Fatalf("Error while writing request to wire : %v", err)
			return
		}

		// redirecting HTTP request to TCP connection
		_, err = conn.Write(buf.Bytes())
		if err != nil {
			log.Fatalf("Error while sending bytes to %v : %v", conn.RemoteAddr(), err)
			return
		}

		// reading response
		_, err = conn.Read(respBytes)
		if err != nil {
			log.Fatalf("Error while reading response from %v : %v", conn.RemoteAddr(), err)
			return
		}

		reader := bytes.NewReader(respBytes)
		respBufio := bufio.NewReader(reader)
		resp, err := http.ReadResponse(respBufio, reply)
		if err != nil {
			log.Fatalf("Error while converting bytes to HTTP response %v", err)
			return
		}

		// sending response to HTTP thread
		channels.ResponseChannel <- resp
	}
}

func createOrSelectChannelForUser(conn net.Conn, channels *ChannelsDomains, db *gorm.DB) (string, error) {
	networkBytes := make([]byte, 1024)

	// received auth message
	_, err := conn.Read(networkBytes)
	if err != nil {
		return "", err
	}

	buf := bytes.NewBuffer(networkBytes)
	dec := gob.NewDecoder(buf)

	var authRequest messages.AuthRequest
	err = dec.Decode(&authRequest)
	if err != nil {
		return "", err
	}

	log.Printf("Email received : %v", authRequest.Email)

	// create domain record in db
	var user database.User
	db.First(&user, "Email = ? ", authRequest.Email)

	// no domain record registered
	var dnsRecord string
	if user.DomainRecordID == 0 {
		dnsRecord = uuid.NewString() + "."+ config.GlobalConfig.GlobalDomainName
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
