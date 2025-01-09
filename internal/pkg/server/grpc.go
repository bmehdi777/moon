package server

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

	"github.com/bmehdi777/moon/internal/pkg/server/authent"
	"github.com/bmehdi777/moon/internal/pkg/server/config"
	"github.com/bmehdi777/moon/internal/pkg/server/database"
	pb "github.com/bmehdi777/moon/protos"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type TunnelServer struct {
	pb.UnimplementedTunnelServer

	ChannelsPerDomain *ChannelsDomains
	Database          *gorm.DB
}

func (ts *TunnelServer) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	srv := grpc.NewServer()
	pb.RegisterTunnelServer(srv, ts)

	listener, err := net.Listen("tcp", "localhost:4040")
	if err != nil {
		return err
	}

	go func() {
		_ = srv.Serve(listener)
		cancel()
	}()

	<-ctx.Done()

	srv.GracefulStop()
	return nil
}

func (ts *TunnelServer) Stream(stream pb.Tunnel_StreamServer) error {
	msg, err := stream.Recv()
	if err == io.EOF {
		// close
		return nil
	}
	if err != nil {
		return err
	}

	if _, ok := msg.Event.(*pb.StreamClient_Credentials); !ok {
		// close
		// send bad creds msg
		return nil
	}

	channelsName, err := ts.login(msg.Event.(*pb.StreamClient_Credentials).Credentials)
	if err != nil {
		fmt.Println("Error : ", err)
		return err
	}

	channels := ts.ChannelsPerDomain.Get(channelsName)
	if channels == nil {
		// close conn
		return nil
	}

	defer ts.ChannelsPerDomain.Delete(channelsName)
	defer ts.Database.Model(&database.DomainRecord{}).Where("dns_record = ?", channelsName).Update("connection_open", false)

	var reply *http.Request
	recvChan := make(chan *pb.StreamClient)

	go func() {
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				// close
				return
			}
			if err != nil {
				log.Fatalf("Error while receiving")
				continue
			}

			recvChan <- msg
		}
	}()

	for {
		select {
		case reply := <-channels.RequestChannel:
			var buf bytes.Buffer
			err := reply.Write(&buf)
			if err != nil {
				log.Fatalf("Error while writing request to wire : %v", err)
			}
			stream.Send(&pb.StreamServer{Event: &pb.StreamServer_HttpRequest{
				HttpRequest: buf.Bytes(),
			}})
		case response := <-recvChan:
			switch response.Event.(type) {
			case *pb.StreamClient_Ping_:
				fmt.Println("Ping")
				stream.Send(&pb.StreamServer{Event: &pb.StreamServer_Pong_{}})
				break
			case *pb.StreamClient_HttpResponse:
				reader := bytes.NewReader(msg.Event.(*pb.StreamClient_HttpResponse).HttpResponse)
				responseBufio := bufio.NewReader(reader)
				responseHttp, err := http.ReadResponse(responseBufio, reply)
				if err != nil {
					log.Fatalf("Error while converting bytes to HTTP response %v", err)
				}
				reply = nil
				channels.ResponseChannel <- responseHttp
				break
			case *pb.StreamClient_Logout_:
				return nil
			default:
				// ignore msg
				break
			}
		}
	}
}

func (ts *TunnelServer) login(creds *pb.StreamClient_Login) (string, error) {
	accessToken, err := authent.VerifyJwt(creds.AccessToken)
	if err != nil {
		return "", err
	}

	sub, err := accessToken.Claims.GetSubject()
	if err != nil {
		return "", err
	}

	// create domain record in db
	var user database.User
	ts.Database.First(&user, "kc_user_id = ?", sub)

	// no domain record registered
	var dnsRecord string
	if user.DomainRecordID == 0 {
		dnsRecord = uuid.NewString() + "." + config.GlobalConfig.App.GlobalDomainName
		record := database.DomainRecord{
			DNSRecord:      dnsRecord,
			ConnectionOpen: true,
		}
		ts.Database.Model(&user).Update("DomainRecord", record)
	} else {
		// refacto one day
		var domainRecord database.DomainRecord
		ts.Database.Table("domain_records").Select("domain_records.dns_record").Joins("INNER JOIN users ON domain_records.ID = users.domain_record_id").Scan(&domainRecord)
		dnsRecord = domainRecord.DNSRecord
		ts.Database.Model(&database.DomainRecord{}).Where("dns_record = ?", dnsRecord).Update("connection_open", true)
	}

	log.Printf("Connection open : http://%v", dnsRecord)

	// create channel
	ts.ChannelsPerDomain.Add(dnsRecord)

	return dnsRecord, nil
}

func (ts *TunnelServer) handleClient() {

}
