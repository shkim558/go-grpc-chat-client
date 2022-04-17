package chatClient

import (
	"bufio"
	"fmt"
	"go-grpc-chat/protoDir"
	"log"
	"os"
	"strings"
)

type ClientHandler struct {
	Stream protoDir.Services_ChatServiceClient
	ClientName string
}

func (ch *ClientHandler) ClientConfig() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Your Name: ")
	name, err := reader.ReadString('\n')
	if err != nil {
		log.Panic(err)
	}
	ch.ClientName = strings.Trim(name, "\r\n")
}

func (ch *ClientHandler) SendMessage() {
	for {
		reader := bufio.NewReader(os.Stdin)
		clientMessage, err := reader.ReadString('\n')
		if err != nil {
			log.Panic(err)
		}
		clientMessage = strings.Trim(clientMessage, "\r\n")

		clientMessageIntf := &protoDir.FromClient{
			Name: ch.ClientName,
			Body: clientMessage,
		}
		err = ch.Stream.Send(clientMessageIntf)
		if err != nil {
			log.Panic(err)
		}
	}
}

func (ch *ClientHandler) ReceiveMessage() {
	for {
		resp, err := ch.Stream.Recv()
		if err != nil {
			log.Panic(err)
		}
		log.Printf("%s: %s", resp.Name, resp.Body)
	}
}