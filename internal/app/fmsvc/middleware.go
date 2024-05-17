package fmsvc

import (
	"context"
	"time"

	"github.com/go-kit/log"
)

type Middleware func(Service) Service

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

// Download implements Service.
func (mw *loggingMiddleware) Download(ctx context.Context, filename string) (fileData []byte, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "Download", "filename", filename, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.Download(ctx, filename)
}

// List implements Service.
func (mw *loggingMiddleware) List(ctx context.Context, sortBy string, order string) (files []File, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "ListFile", "sortby", sortBy, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.List(ctx, sortBy, order)
}

// SearchByName implements Service.
func (mw *loggingMiddleware) SearchByName(ctx context.Context, filename string) (files []File, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "SearchByName", "filename", filename, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.SearchByName(ctx, filename)
}

func (mw loggingMiddleware) Upload(ctx context.Context, filename string, data []byte) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "UploadFile", "filename", filename, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.Upload(ctx, filename, data)
}
