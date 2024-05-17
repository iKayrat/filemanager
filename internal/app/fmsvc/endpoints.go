package fmsvc

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	UploadFileEndpoint     endpoint.Endpoint
	GetListFileEndpoint    endpoint.Endpoint
	SearchFilenameEndpoint endpoint.Endpoint
	DownloadFileEndpoint   endpoint.Endpoint
}

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		UploadFileEndpoint:     makeUploadFileEndpoint(s),
		GetListFileEndpoint:    makeGetListFileEndpoint(s),
		SearchFilenameEndpoint: makeSearchFilenameEndpoint(s),
		DownloadFileEndpoint:   makeDownloadEndpoint(s),
	}
}

// MakeUploadFileEndpoint returns an endpoint via the passed service.
func makeUploadFileEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UploadFileRequest)
		e := s.Upload(ctx, req.Filename, req.Data)
		return UploadFileResponse{Err: e}, nil
	}
}

// MakeGetListFileEndpoint returns an endpoint via the passed service.
func makeGetListFileEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetListFileRequest)
		files, e := s.List(ctx, req.SortBy, req.Order)
		return GetListFileResponse{
			Files: files,
			Err:   e,
		}, nil
	}
}

func makeSearchFilenameEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(SearchFilenameRequest)
		files, e := s.SearchByName(ctx, req.Filename)
		return SearchFilenameResponse{
			Files: files,
			Err:   e,
		}, nil
	}
}

func makeDownloadEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(DownloadFileRequest)
		fileData, e := s.Download(ctx, req.Filename)
		return DownloadFileResponse{
			FileData: fileData,
			Err:      e,
		}, nil
	}
}


