package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"go-grpc-chat/chatClient"
	"go-grpc-chat/protoDir"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	host string = "global-grpc-chat-server-muoq7g5tta-uw.a.run.app:443"
	// host string = "localhost:3051"
)

func main() {
	var opts []grpc.DialOption
	if host != "" {
		opts = append(opts, grpc.WithAuthority(host))
	}

	// local, http connect
	// opts = append(opts, grpc.WithInsecure())

	// https connect
	systemRoots, err := x509.SystemCertPool()
	if err != nil {
		log.Panic(err)
	}
	cred := credentials.NewTLS(&tls.Config{
					RootCAs: systemRoots,
	})
	opts = append(opts, grpc.WithTransportCredentials(cred))

	conn, err := grpc.Dial(host, opts...)
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()

	client := protoDir.NewServicesClient(conn)
	stream, err := client.ChatService(context.Background())
	if err != nil {
		log.Panic(err)
	}

	ch := chatClient.ClientHandler{
		Stream: stream,
	}
	ch.ClientConfig()
	go ch.SendMessage()
	go ch.ReceiveMessage()

	block := make(chan bool)
	<- block
}
