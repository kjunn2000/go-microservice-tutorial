package server

import (
	"context"

	protos "github.com/kjunn2000/go-microservice-tutorial/currency-server/protos/currency"
)

type Currency struct {
	protos.UnimplementedCurrencyServer
}

func NewCurrency() *Currency {
	return &Currency{}
}

func (c *Currency) GetRate(ctx context.Context, rr *protos.RateRequest) (*protos.RateResponse, error) {
	return &protos.RateResponse{Rate: 0.5}, nil
}
