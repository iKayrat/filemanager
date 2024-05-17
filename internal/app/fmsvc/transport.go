package fmsvc

// import (
// 	"context"
// 	"encoding/json"

// 	"github.com/go-kit/kit/log"

// 	"net/http"

// 	"github.com/go-kit/kit/transport"
// 	httptransport "github.com/go-kit/kit/transport/http"
// 	"github.com/gorilla/mux"
// )

// func MakeHTTPHandler(s Service, logger log.Logger) http.Handler {
// 	r := mux.NewRouter()
// 	e := MakeServerEndpoints(s)
// 	options := []httptransport.ServerOption{
// 		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
// 		httptransport.ServerErrorEncoder(encodeError),
// 	}

// 	// POST    /profiles/                          adds another profile
// 	// GET     /profiles/:id                       retrieves the given profile by id
// 	// PUT     /profiles/:id                       post updated profile information about the profile
// 	// PATCH   /profiles/:id                       partial updated profile information
// 	// DELETE  /profiles/:id                       remove the given profile
// 	// GET     /profiles/:id/addresses/            retrieve addresses associated with the profile
// 	// GET     /profiles/:id/addresses/:addressID  retrieve a particular profile address
// 	// POST    /profiles/:id/addresses/            add a new address
// 	// DELETE  /profiles/:id/addresses/:addressID  remove an address

// 	r.Methods("POST").Path("/files/").Handler(httptransport.NewServer(
// 		e.PostFileEndpoint,
// 		decodePostFileRequest,
// 		encodeResponse,
// 		options...,
// 	))

// 	return r

// }

// func decodePostFileResponse(_ context.Context, resp *http.Response) (interface{}, error) {
// 	var response UploadFileResponse
// 	err := json.NewDecoder(resp.Body).Decode(&response)
// 	return response, err
// }

// func decodePostFileRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
// 	var req UploadFileRequest
// 	if e := json.NewDecoder(r.Body).Decode(&req.Path); e != nil {
// 		return nil, e
// 	}
// 	return req, nil
// }

// func encodePostFileRequest(ctx context.Context, req *http.Request, request interface{}) error {
// 	// r.Methods("POST").Path("/profiles/")
// 	req.URL.Path = "/files/"
// 	return encodeRequest(ctx, req, request)
// }

// type errorer interface {
// 	error() error
// }

// func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
// 	if e, ok := response.(errorer); ok && e.error() != nil {
// 		// Not a Go kit transport error, but a business-logic error.
// 		// Provide those as HTTP errors.
// 		encodeError(ctx, e.error(), w)
// 		return nil
// 	}
// 	w.Header().Set("Content-Type", "application/json; charset=utf-8")
// 	return json.NewEncoder(w).Encode(response)
// }

// func encodeError(_ context.Context, err error, w http.ResponseWriter) {
// 	if err == nil {
// 		panic("encodeError with nil error")
// 	}
// 	w.Header().Set("Content-Type", "application/json; charset=utf-8")
// 	w.WriteHeader(codeFrom(err))
// 	json.NewEncoder(w).Encode(map[string]interface{}{
// 		"error": err.Error(),
// 	})
// }

// func codeFrom(err error) int {
// 	switch err {
// 	case ErrNotFound:
// 		return http.StatusNotFound
// 	case ErrAlreadyExists, ErrInconsistentIDs:
// 		return http.StatusBadRequest
// 	default:
// 		return http.StatusInternalServerError
// 	}
// }
