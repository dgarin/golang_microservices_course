package main

import (
	"context"
	"log"
	"main/handlers"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

//var bindAddress = env.String("BIND_ADDRESS", false, ":8080", "Bind addres for the server")

func main() {
	l := log.New(os.Stdout, "product-api ", log.LstdFlags)
	// create the handlers
	ph := handlers.NewProducts(l)

	// create a new serve mux and register the handlers
	sm := mux.NewRouter()
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", ph.GetProducts)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", ph.UpdateProducts)
	putRouter.Use(ph.MiddlewareProductValidation)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", ph.AddProducts)
	postRouter.Use(ph.MiddlewareProductValidation)

	//create a new server
	s := http.Server{
		Addr:         ":8080",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		l.Println("Starting server on port 8080")
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)
	sig := <-sigChan
	l.Println("Received terminate, graceful shutdown", sig)
	tc, cancel := context.WithTimeout(context.Background(), time.Duration(30)*time.Second)
	defer cancel()
	s.Shutdown(tc)
}
