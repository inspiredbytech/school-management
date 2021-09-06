package school

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

// schoolCreateRequest represents an HTTP request from the client for school creation
type schoolCreateRequest struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Country  string   `json:"country"`
	City     string   `json:"city"`
	Address  string   `json:"Address"`
	Contacts []string `json:"Contacts"`
}

// schoolCreateResponse represents an HTTP response from our server for school creation
type schoolCreateResponse struct {
	ID    int   `json:"id,omitempty"`
	Error error `json:"error,omitempty"`
}

// error is the schoolCreateResponse errorer implementation
func (r schoolCreateResponse) error() error { return r.Error }

// makeCreateSchoolEndpoint generates a service endpoint for schools
func makeCreateSchoolEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*schoolCreateRequest)

		id, err := s.CreateSchool(*New(req.ID, req.Name, req.Country, req.City, req.Address, req.Contacts))

		return schoolCreateResponse{ID: id, Error: err}, nil
	}
}

// schoolReadRequest represents an HTTP request to read a single school from the client
type schoolReadRequest struct {
	ID int `json:"id"`
}

// schoolReadResponse represents an HTTP response containing a school or the error when fetching
type schoolReadResponse struct {
	School School `json:"school,omitempty"`
	Error  error  `json:"error,omitempty"`
}

// error is the schoolReadResponse errorer implementation
func (r schoolReadResponse) error() error { return r.Error }

func makeReadSchoolEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(schoolReadRequest)
		u, err := s.GetSchool(req.ID)
		return schoolReadResponse{School: u, Error: err}, nil
	}
}

// schoolReadAllRequest represents an HTTP request from the client to get all schools
type schoolReadAllRequest struct{}

// schoolReadAllResponse represents an HTTP response from the server listing all schools
type schoolReadAllResponse struct {
	Schools []*School `json:"schools,omitempty"`
	Error   error     `json:"error,omitempty"`
}

// error is an errorer implementation for schoolReadAllResponse
func (r schoolReadAllResponse) error() error { return r.Error }

// makeReadAllSchoolsEndpoint creates an HTTP endpoint for retrieving all schools
func makeReadAllSchoolsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		schools := s.GetSchools()
		return schoolReadAllResponse{Schools: schools, Error: nil}, nil
	}
}
