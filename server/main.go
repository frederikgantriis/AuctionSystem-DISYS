package main

import (
	"context"
	"io"
	"log"
	"net"
	"os"
	"time"

	auction "github.com/frederikgantriis/AuctionSystem-DISYS/gRPC"
	"google.golang.org/grpc"
)

func main() {
	file, _ := openLogFile("./server/serverlog.log")

	mw := io.MultiWriter(os.Stdout, file)
	log.SetOutput(mw)
	log.SetFlags(2 | 3)

	log.Println("Hello World!")

	if len(os.Args) != 2 {
		log.Printf("Please input a number to run the server on")
	}

	listen, _ := net.Listen("tcp", "localhost:300"+os.Args[1])

	grpcServer := grpc.NewServer()
	auction.RegisterAuctionServer(grpcServer, &Server{
		highestBid:        0,
		timeLeft:          -1,
		currentWinnerUser: "",
	})

	log.Printf("server listening at %v", listen.Addr())

	grpcServer.Serve(listen)
}

func (s *Server) bid(ctx context.Context, req *auction.BidRequest) (*auction.BidReply, error) {
	if s.timeLeft == -1 {
		s.timeLeft = 60

		go func() {
			for s.timeLeft > 0 {
				s.timeLeft--
				time.Sleep(time.Second)
			}
		}()

		log.Printf("Auction started")
	}

	if (req.Bid > s.highestBid) && (s.timeLeft > 0) {
		s.highestBid = req.Bid
		s.currentWinnerUser = req.User
		return &auction.BidReply{Outcome: auction.Outcomes(SUCCESS)}, nil
	} else {
		return &auction.BidReply{Outcome: auction.Outcomes(FAIL)}, nil
	}
}

func (s *Server) result(ctx context.Context, resReq *auction.ResultRequest) (*auction.ResultReply, error) {
	return &auction.ResultReply{User: s.currentWinnerUser, HighestBid: s.highestBid}, nil
}

func openLogFile(path string) (*os.File, error) {
	logFile, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return logFile, nil
}

type Server struct {
	auction.UnimplementedAuctionServer
	highestBid        int32
	currentWinnerUser string
	timeLeft          int32
}

type Outcomes int32

const (
	FAIL      Outcomes = 0
	SUCCESS   Outcomes = 1
	EXCEPTION Outcomes = 2
)
