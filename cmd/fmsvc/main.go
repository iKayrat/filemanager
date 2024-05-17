package main

import (
	"context"
	"database/sql"
	"filemanagerService/internal/app/fmsvc"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/log"
	_ "github.com/lib/pq"
)

func main() {
	var (
		httpAddr = flag.String("http.addr", ":8080", "HTTP listen address")
	)
	flag.Parse()
	ctx := context.Background()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	conn, err := sql.Open("postgres", "postgresql://root:root@localhost:5431/fmanagerdb?sslmode=disable")
	if err != nil {
		logger.Log("Cannot connect to db", err)
	}

	var s fmsvc.Service
	{
		s = fmsvc.New(conn)
		s = fmsvc.LoggingMiddleware(logger)(s)
	}

	endpoints := fmsvc.MakeEndpoints(s)
	var h http.Handler
	{
		h = fmsvc.NewHTTPServer(ctx, endpoints, log.With(logger, "component", "HTTP"))
	}

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("transport", "HTTP", "addr", *httpAddr)

		errs <- http.ListenAndServe(*httpAddr, h)
	}()

	logger.Log("exit", <-errs)
}
