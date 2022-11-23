package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	auction "github.com/frederikgantriis/AuctionSystem-DISYS/gRPC"
	"google.golang.org/grpc"
)

func main() {
	file, _ := openLogFile("./serverlog.log")

	log.SetOutput(file)
	log.SetFlags(2 | 3)

	clients := make([]auction.AuctionClient, 3)

	for i := 0; i < 3; i++ {
		port := int32(3000) + int32(i)

		fmt.Printf("Trying to dial: %v\n", port)
		conn, err := grpc.Dial(fmt.Sprintf(":%v", port), grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("Could not connect: %s", err)
		}
		clients[i] = auction.NewAuctionClient(conn)
		defer conn.Close()
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		command := strings.Split(scanner.Text(), " ")
		command[0] = strings.ToLower(command[0])

		if command[0] == "bid" {
			for _, client := range clients {
				log.Println(client)
			}
		} else if command[0] == "result" {

		}

		if err != nil {
			log.Panicln(err)
		}

	}
}

func openLogFile(path string) (*os.File, error) {
	logFile, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return logFile, nil
}
