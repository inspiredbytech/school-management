package stringsvc

import (
	"context"
	"encoding/json"
	"net/http"

	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

// MakeHandler builds a go-kit http transport and returns it
func MakeHandler(svc StringService, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
	}

	uppercase := makeUppercaseEndpoint(svc)

	uppercaseHandler := kithttp.NewServer(
		uppercase,
		decodeUppercaseRequest,
		encodeResponse,
		opts...,
	)

	strcount := makeCountEndpoint(svc)

	strcountHandler := kithttp.NewServer(
		strcount,
		decodeCountRequest,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()
	//r.Handle("/api/v1/users", upperCaseHandler).Methods("GET")
	//r.Handle("/api/v1/users", createHandler).Methods("POST")
	r.Handle("/api/v1/string", uppercaseHandler).Methods("GET")
	r.Handle("/api/v1/string", uppercaseHandler).Methods("POST")
	r.Handle("/api/v1/string/upper", uppercaseHandler).Methods("POST")
	r.Handle("/api/v1/string/count", strcountHandler).Methods("POST")
	//r.Handle("/api/v1/users/{s}", updateHandler).Methods("PUT")
	return r
}

type uppercaseRequest struct {
	S string `json:"s"`
}

type uppercaseResponse struct {
	V   string `json:"v"`
	Err string `json:"err,omitempty"` // errors don't JSON-marshal, so we use a string
}

type countRequest struct {
	S string `json:"s"`
}

type countResponse struct {
	V int `json:"v"`
}

func decodeUppercaseRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request uppercaseRequest
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
