package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	auction "github.com/frederikgantriis/AuctionSystem-DISYS/gRPC"
	"google.golang.org/grpc"
)

func main() {
	username := os.Args[1]
	file, _ := openLogFile("./client/clientlog.log")

	mw := io.MultiWriter(os.Stdout, file)
	log.SetOutput(mw)
	log.SetFlags(2 | 3)

	log.Println("Hello World")

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
			bidAmount, _ := strconv.Atoi(command[1])
			bid := &auction.BidRequest{User: username, Bid: int32(bidAmount)}
			for id, client := range clients {
				res, err := client.Bid(ctx, bid)
				if err != nil {
					log.Printf("ERROR: %v", err)
					return
				}

				log.Printf("user %v: %v", username, res.GetOutcome())

				if res.GetOutcome() == auction.Outcomes_SUCCESS {
					log.Printf("user %v made a succesfull bid at server %v, for amount: %v", username, id, bidAmount)
				} else if res.GetOutcome() == auction.Outcomes_FAIL {
					log.Printf("user %v: bid was either too low or auction has ended", username)
				}
			}
		} else if command[0] == "result" {
			for _, client := range clients {
				client.Result(ctx, &auction.ResultRequest{})
			}
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
