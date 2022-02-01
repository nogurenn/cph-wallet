package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shopspring/decimal"

	"github.com/nogurenn/cph-wallet/dbutil"
	"github.com/nogurenn/cph-wallet/transaction"
)

func main() {
	httpAddress := flag.String("http.addr", ":8080", "HTTP listen address")
	flag.Parse()

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	cfg := dbutil.NewConfig()
	db, err := dbutil.NewDb(cfg)
	if err != nil {
		logger.Log("fatal", "db connection could not be established")
		panic(err)
	}

	fieldKeys := []string{"method"}

	var tdb transaction.Repository
	tdb = transaction.NewPostgresDb(db)

	var ts transaction.Service
	ts = transaction.NewService(tdb)
	ts = transaction.NewLoggingService(log.With(logger, "component", "transaction"), ts)
	ts = transaction.NewInstrumentingService(
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "transaction_service",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "transaction_service",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys),
		ts,
	)

	// make loading of test data mandatory for this exam code.
	err = setupTestData(ts)
	if err != nil {
		logger.Log("fatal", "test data could not be loaded to the repository")
		panic(err)
	}

	httpLogger := log.With(logger, "component", "http")

	mux := http.NewServeMux()
	mux.Handle("/transaction/v1/", transaction.MakeHandler(ts, httpLogger))

	http.Handle("/", mux)
	http.Handle("/metrics", promhttp.Handler())

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

// setupTestData loads test data to repositories.
func setupTestData(ts transaction.Service) error {
	usernames := []string{"bob123", "alice456", "karen789"}
	initialBalance := decimal.NewFromFloat(200.00)

	for _, username := range usernames {
		err := ts.CreateAccount(username)
		if err != nil {
			return err
		}

		err = ts.Deposit(username, initialBalance)
		if err != nil {
			return err
		}
	}

	return nil
}
