package server

import (
	"context"
	"log"

	protos "github.com/kjunn2000/grpc-server/protos/currency"
)

type Currency struct {
	protos.UnimplementedCurrencyServer
}

func NewCurrency() *Currency {
	return &Currency{}
}

func (c *Currency) GetRate(ctx context.Context, rr *protos.RateRequest) (*protos.RateResponse, error) {
	log.Println("Request is " + rr.GetBase() + " " + rr.GetDestination())
	return &protos.RateResponse{Rate: 0.5}, nil
}
