package school

import (
	"net/http"

	"context"
	"encoding/json"

	"io/ioutil"

	"strconv"

	errs "schoolmgt/errors"

	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

type errorer interface {
	error() error
}

func MakeHandler(svc Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	// Define all endpoints
	create := makeCreateSchoolEndpoint(svc)
	read := makeReadSchoolEndpoint(svc)
	//update := makeUpdateSchoolColorEndpoint(svc)
	list := makeReadAllSchoolsEndpoint(svc)

	createHandler := kithttp.NewServer(
		create,
		decodeCreateSchoolRequest,
		encodeCreateSchoolResponse,
		opts...,
	)

	readHandler := kithttp.NewServer(
		read,
		decodeReadSchoolRequest,
		encodeReadSchoolResponse,
		opts...,
	)

	listHandler := kithttp.NewServer(
		list,
		decodeListSchoolsRequest,
		encodeListSchoolsResponse,
		opts...,
	)

	r := mux.NewRouter()
	r.Handle("/api/v1/schools", listHandler).Methods("GET")
	r.Handle("/api/v1/schools", createHandler).Methods("POST")
	r.Handle("/api/v1/schools/{id}", readHandler).Methods("GET")
	//r.Handle("/api/v1/schools/{id}", updateHandler).Methods("PUT")

	return r
}

func decodeRequest(to interface{}, r *http.Request) (interface{}, error) {
	d, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(d, &to)
	return to, err
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	d, err := json.Marshal(response)
	if err != nil {
		encodeError(ctx, err, w)
	}
	_, err = w.Write(d)
	return err
}

func decodeCreateSchoolRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req schoolCreateRequest
	return decodeRequest(&req, r)
}

func encodeCreateSchoolResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}

	res := response.(schoolCreateResponse)
	return encodeResponse(ctx, w, res)
}

func decodeReadSchoolRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	req := schoolReadRequest{
		ID: id,
	}
	return req, err
}

func encodeReadSchoolResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}

	res := response.(schoolReadResponse)
	return encodeResponse(ctx, w, res)
}

func decodeListSchoolsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := schoolReadAllRequest{}
	return req, nil
}

func encodeListSchoolsResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}

	res := response.(schoolReadAllResponse)
	return encodeResponse(ctx, w, res)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	switch err {
	case errs.ErrInvalidArgument:
		w.WriteHeader(http.StatusBadRequest)
	case errs.ErrSchoolNotFound:
		w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
