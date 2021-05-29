package main

import (
	"log"
	"net"
	"os"

	protos "github.com/kjunn2000/grpc-server/protos/currency"
	"github.com/kjunn2000/grpc-server/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	gs := grpc.NewServer()
	cs := server.NewCurrency()

	protos.RegisterCurrencyServer(gs, cs)

	reflection.Register(gs)

	l, err := net.Listen("tcp", ":9050")
	if err != nil {
		log.Fatal("Connection error" + err.Error())
		os.Exit(1)
	}

	gs.Serve(l)

}
