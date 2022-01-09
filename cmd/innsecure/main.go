package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/form3tech/innsecure"
	"github.com/form3tech/innsecure/jwtauth"
	"github.com/form3tech/innsecure/postgres"
	"github.com/go-kit/kit/log"
	_ "github.com/lib/pq"
)

func main() {
	var (
		httpAddr = flag.String("http.addr", ":8080", "HTTP listen address")
	)
	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewJSONLogger(os.Stdout)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var s innsecure.Service
	{
		host := os.Getenv("DB_HOST")
		user := os.Getenv("DB_USER")
		pass := os.Getenv("DB_PASSWORD")
		db, err := postgres.NewConnection(host, user, pass)
		if err != nil {
			panic(err)
		}
		defer db.Close()

		r := postgres.NewRepo(db)
		s = innsecure.NewBookingService(r)
	}

	var h http.Handler
	{
		jwtmw := jwtauth.NewMiddleware("SigningString")
		e := innsecure.MakeServerEndpoints(s, jwtmw)
		h = innsecure.MakeHTTPHandler(e, log.With(logger, "component", "HTTP"))
	}

	errs := make(chan error)

	// Shutdown handler
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	// HTTP Transport
	go func() {
		logger.Log("transport", "HTTP", "addr", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, h)
	}()

	logger.Log("exit", <-errs)
}
