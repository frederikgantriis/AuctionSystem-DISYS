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
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	username := os.Args[1]

	file, _ := openLogFile("./client/clientlog.log")

	mw := io.MultiWriter(os.Stdout, file)
	log.SetOutput(mw)
	log.SetFlags(2 | 3)

	servers := make([]auction.AuctionClient, 3)

	for i := 0; i < 3; i++ {
		port := int32(3000) + int32(i)

		fmt.Printf("Trying to dial: %v\n", port)
		conn, err := grpc.Dial(fmt.Sprintf(":%v", port), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
		if err != nil {
			log.Fatalf("Could not connect: %s", err)
		}
		servers[i] = auction.NewAuctionClient(conn)
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

			w := 0
			var reply *auction.BidReply
			for _, server := range servers {
				res, err := server.Bid(ctx, bid)
				if err != nil {
					log.Printf("user %v: ERROR - %v", username, err)
					continue
				}
				w++

				log.Printf("user %v: %v", username, res.GetOutcome())

				if reply == nil || (reply.GetOutcome() != res.GetOutcome() && reply.GetOutcome() == auction.Outcomes_SUCCESS) {
					reply = res
				}
			}

			if w >= 2 {
				if reply.GetOutcome() == auction.Outcomes_SUCCESS {
					log.Printf("user %v: made a succesfull bid, for amount: %v", username, bidAmount)
				} else if reply.GetOutcome() == auction.Outcomes_FAIL {
					log.Printf("user %v: bid was either too low or auction has ended", username)
				}
			} else {
				log.Printf("user %v: Call successful server writes and delete write (not implemented)", username)
			}
		} else if command[0] == "result" {
			r := 0
			var reply *auction.ResultReply
			for _, server := range servers {
				res, err := server.Result(ctx, &auction.ResultRequest{})
				if err != nil {
					log.Printf("user %v: ERROR - %v", username, err)
					continue
				}
				r++

				// Take the first result
				if reply == nil {
					reply = res
					continue
				}

				// Update highestBid
				if res.GetHighestBid() > reply.GetHighestBid() {
					reply.HighestBid = res.GetHighestBid()
					reply.User = res.GetUser()
				}

				// Update timeleft
				if res.GetTimeLeft() < reply.GetTimeLeft() {
					reply.TimeLeft = reply.GetTimeLeft()
				}
			}

			if r >= 2 {
				if reply.GetTimeLeft() > 0 {
					log.Printf("user %v: Current highest bid: %v by %v\n The auction ends in: %v seconds", username, reply.GetHighestBid(), reply.GetUser(), reply.GetTimeLeft())
				} else {
					log.Printf("user %v: The auction has ended. The winner was %v with the bid: %v", username, reply.GetUser(), reply.GetHighestBid())
				}
			} else {
				log.Printf("user %v: Didn't receive enough answers to get a usefull result", username)
			}
		} else if command[0] == "reset" {
			for _, server := range servers {
				_, err := server.Reset(ctx, &auction.ResetRequest{})
				if err != nil {
					log.Printf("user %v: ERROR - %v", username, err)
					continue
				}
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
