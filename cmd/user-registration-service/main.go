package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/pdkonovalov/user-registration-service/pkg/config"
	"github.com/pdkonovalov/user-registration-service/pkg/email"
	"github.com/pdkonovalov/user-registration-service/pkg/http"
	"github.com/pdkonovalov/user-registration-service/pkg/jwt"
	"github.com/pdkonovalov/user-registration-service/pkg/storage/postgres"
)

func main() {
	err := run(context.Background(), os.Getenv)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, getenv func(string) string) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	config, err := config.ReadConfig(getenv)
	if err != nil {
		return fmt.Errorf("error load config: %s", err)
	}

	storage, err := postgres.Init(config)
	if err != nil {
		return fmt.Errorf("error init storage: %s", err)
	}

	email, err := email.Init(config)
	if err != nil {
		return fmt.Errorf("error init email: %s", err)
	}

	jwt, err := jwt.Init(config)
	if err != nil {
		return fmt.Errorf("error init jwt: %s", err)
	}

	server := http.MakeServer(config, storage, email, jwt)
	err = server.Start()
	if err != nil {
		return fmt.Errorf("error start server: %s", err)
	}

	log.Printf("start user-registration-service at %s:%s", config.Host, config.Port)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		log.Print("shutdown user-registration-service")
		err := server.Shutdown()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error shutdown server: %s\n", err)
		}
		err = storage.Shutdown()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error shutdown storage: %s\n", err)
		}
	}()
	wg.Wait()
	return nil
}
