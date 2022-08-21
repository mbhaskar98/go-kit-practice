package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
)

type UppercaseRequest struct {
	Str string `json:"str"`
}

type UppercaseResponse struct {
	Str string `json:"str"`
	Err string `json:"err"`
}

type countRequest struct {
	Str string `json:"str"`
}

type countResponse struct {
	Count int `json:"count"`
}

func makeUppercaseEndpoint(svc StringService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UppercaseRequest)
		str, err := svc.Uppercase(req.Str)
		if err != nil {
			return UppercaseResponse{Str: str, Err: err.Error()}, nil
		}
		return UppercaseResponse{Str: str, Err: ""}, nil
	}
}

func makeCountEndpoint(svc StringService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(countRequest)
		count := svc.Count(req.Str)
		return countResponse{Count: count}, nil
	}
}

func decodeUppercaseRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request UppercaseRequest
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
