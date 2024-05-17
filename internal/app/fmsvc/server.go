package fmsvc

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
)

func NewHTTPServer(ctx context.Context, endpoints Endpoints, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	r.Use(commonMiddleware)

	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Methods("POST").Path("/file").Handler(httptransport.NewServer(
		endpoints.UploadFileEndpoint,
		decodeUploadFileRequest,
		encodeResponse,
		options...,
	))

	r.Methods("GET").Path("/files/{sortby}/{orderby}").Handler(httptransport.NewServer(
		endpoints.GetListFileEndpoint,
		decodeGetListFileRequest,
		encodeResponse,
		options...,
	))

	r.Methods("GET").Path("/files/{filename}").Handler(httptransport.NewServer(
		endpoints.SearchFilenameEndpoint,
		decodeSearchFilenameRequest,
		encodeResponse,
		options...,
	))

	r.Methods("GET").Path("/download").Handler(httptransport.NewServer(
		endpoints.DownloadFileEndpoint,
		decodeDownloadRequest,
		encodeDownloadFileResponse,
		options...,
	))

	return r
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}
