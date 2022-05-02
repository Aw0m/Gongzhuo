package main

import (
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"time"
	"wxProjectDev/user"
	"wxProjectDev/work"
)

var (
	g errgroup.Group
)

func main() {
	server01 := &http.Server{
		Addr:         ":8000",
		Handler:      user.Controller(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	server02 := &http.Server{
		Addr:         ":8001",
		Handler:      work.Controller(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	g.Go(func() error {
		return server01.ListenAndServe()
	})

	g.Go(func() error {
		return server02.ListenAndServe()
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
