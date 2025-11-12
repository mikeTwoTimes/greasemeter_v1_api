package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sendgrid/sendgrid-go"
)

type App struct {
	port      int
	db        *pgxpool.Pool
	jwtSecret string
	mailer    *sendgrid.Client
}

func NewApp(addr, dbConn, secret, sgKey string) (*App, error) {
	port, err := strconv.Atoi(addr)

	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.New(context.Background(), dbConn)

	if err != nil {
		return nil, err
	}

	return &App{
		port:      port,
		db:        pool,
		jwtSecret: secret,
		mailer:    sendgrid.NewSendClient(sgKey),
	}, nil
}

func (a *App) Serve() error {
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", a.port),
		Handler:      a.handler(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(stop)

	go func() {
		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			log.Fatalf("could not listen: %v\n", err)
		}
	}()

	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	a.db.Close()

	return nil
}
