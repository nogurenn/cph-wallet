package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/log"
)

func main() {
	httpAddress := flag.String("http.addr", ":8080", "HTTP listen address")
	flag.Parse()

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	//httpLogger := log.With(logger, "component", "http")

	mux := http.NewServeMux()

	mux.HandleFunc("/hello", func(writer http.ResponseWriter, request *http.Request) {
		io.WriteString(writer, "Hello World!\n")
	})

	http.Handle("/", mux)

	errs := make(chan error, 2)
	go func() {
		logger.Log("transport", "http", "address", *httpAddress, "msg", "listening")

		errs <- http.ListenAndServe(*httpAddress, nil)
	}()
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("terminated", <-errs)
}
