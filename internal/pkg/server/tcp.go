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
		log.Fatalf("TLS - %v", err)
		return
	}
	configTls := tls.Config{Certificates: []tls.Certificate{cert}}
	configTls.Rand = rand.Reader

	fullAddrFmt := fmt.Sprintf("%v:%v", config.GlobalConfig.TcpAddr, config.GlobalConfig.TcpPort)
	listener, err := tls.Listen("tcp", fullAddrFmt, &configTls)
	if err != nil {
		log.Fatalf("TLS - TCP - %v", err)
		return
	}
	defer listener.Close()
	log.Printf("TCP Server is up at %v", fullAddrFmt)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("TCP - CONN - %v", err)
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
	defer conn.Close()
	// create random domain name
	channelsName, err := createChannelForUser(conn, channelsDomains, db)
	if err != nil {
		log.Fatalf("Error while creating channels : %v", err)
	}

	channels := channelsDomains.Get(channelsName)
	if channels == nil {
		log.Fatalf("Error while retrieving channel.")
	}

	log.Printf("Connection started with %v", conn.RemoteAddr())
	respBytes := make([]byte, 1024)

	for {
		reply := <-channels.RequestChannel
		// temp

		var buf bytes.Buffer
		err := reply.Write(&buf)
		if err != nil {
			log.Fatalf("TCP - HTTP - REQ - %v", err)
			return
		}

		_, err = conn.Write(buf.Bytes())
		if err != nil {
			log.Fatalf("TCP - CONN - WRITE - %v", err)
			return
		}

		_, err = conn.Read(respBytes)
		if err != nil {
			log.Fatalf("TCP - CONN - READ - %v", err)
			return
		}

		reader := bytes.NewReader(respBytes)
		respBufio := bufio.NewReader(reader)
		resp, err := http.ReadResponse(respBufio, reply)
		if err != nil {
			log.Fatalf("TCP - READER - %v", err)
			return
		}

		channels.ResponseChannel <- resp
	}
}

func createChannelForUser(conn net.Conn, channels *ChannelsDomains, db *gorm.DB) (string, error) {
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

	log.Printf("User get : %v", user)

	dnsRecord := uuid.NewString() + ".m00n.fr"
	record := database.DomainRecord{
		DNSRecord:      dnsRecord,
		ConnectionOpen: true,
	}
	db.Model(&user).Update("DomainRecord", record)

	log.Printf("Record created : http://%v", dnsRecord)

	// create channel
	channels.Add(dnsRecord)

	return dnsRecord, nil
}
