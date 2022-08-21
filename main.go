package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/go-kit/kit/endpoint"

	httptransport "github.com/go-kit/kit/transport/http"
)

// ErrEmpty is returned when input string is empty
var ErrEmpty = errors.New("empty string")

type StringService interface {
	UpperCase(*string) (string, error)
	Count(string) int
}

type stringService struct{}

func (s stringService) UpperCase(str *string) (string, error) {
	if *str == "" {
		return "", ErrEmpty
	}
	return strings.ToUpper(*str), nil
}

func (s stringService) Count(str string) int {
	return len(str)
}

type upperCaseRequest struct {
	Str string `json:"str"`
}

type upperCaseResponse struct {
	Str string `json:"str"`
	Err string `json:"err"`
}

type countRequest struct {
	Str string `json:"str"`
}

type countResponse struct {
	Count int `json:"count"`
}

func makeUpperCaseEndpoint(svc StringService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(upperCaseRequest)
		str, err := svc.UpperCase(&req.Str)
		if err != nil {
			return upperCaseResponse{Str: str, Err: err.Error()}, nil
		}
		return upperCaseResponse{Str: str, Err: ""}, nil
	}
}

func makeCountEndpoint(svc StringService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(countRequest)
		count := svc.Count(req.Str)
		return countResponse{Count: count}, nil
	}
}

func main() {
	svc := stringService{}

	upperCaseHandler := httptransport.NewServer(
		makeUpperCaseEndpoint(svc),
		decodeUppercaseRequest,
		encodeResponse,
	)

	countHandler := httptransport.NewServer(
		makeCountEndpoint(svc),
		decodeCountRequest,
		encodeResponse,
	)

	
	http.Handle("/uppercase", upperCaseHandler)
	http.Handle("/count", countHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func decodeUppercaseRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request upperCaseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeCountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request countRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
