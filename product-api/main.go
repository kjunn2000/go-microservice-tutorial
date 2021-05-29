package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/kjunn2000/product-api/files"
	"github.com/kjunn2000/product-api/handlers"
)

func main() {

	stor, err := files.NewLocal("./imagestore", 1000*1024*5)
	if err != nil {
		log.Fatalln("Error: ", err)
	}

	fh := handlers.NewFileHandler(stor)

	sm := mux.NewRouter()

	ph := sm.Methods(http.MethodPost).Subrouter()
	ph.HandleFunc("/images/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-z]{3}}", fh.ServeHTTP)

	gh := sm.Methods(http.MethodGet).Subrouter()
	gh.Handle(
		"/iamges/{id:[0-9]+}/{filename:[a-aA-Z]+\\.[a-z]{3}}",
		http.StripPrefix("/images/", http.FileServer(http.Dir("/imagestore"))),
	)

	s := &http.Server{
		Addr:         ":8080",
		Handler:      sm,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}
	go func() {
		s.ListenAndServe()
	}()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	res := <-ch
	fmt.Println("Server shut down", res)

	c, _ := context.WithTimeout(context.Background(), time.Millisecond*50000)

	s.Shutdown(c)
}
