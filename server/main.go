package main

import (
	"context"
	"io"
	"log"
	"net"
	"os"

	auction "github.com/frederikgantriis/AuctionSystem-DISYS/gRPC"
	"google.golang.org/grpc"
)

type Server struct {
	auction.UnimplementedAuctionServer
}

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
	auction.RegisterAuctionServer(grpcServer, &Server{})

	log.Printf("server listening at %v", listen.Addr())

	grpcServer.Serve(listen)
}

func (s *Server) bid(ctx context.Context, req *auction.BidRequest) (*auction.BidReply, error) {
	return nil, nil
}

func (s *Server) result(ctx context.Context, resReq *auction.ResultRequest) (*auction.ResultReply, error) {
	return nil, nil
}

func openLogFile(path string) (*os.File, error) {
	logFile, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return logFile, nil
}

type Outcomes int32

const (
	FAIL      Outcomes = 0
	SUCCESS   Outcomes = 1
	EXCEPTION Outcomes = 2
)
