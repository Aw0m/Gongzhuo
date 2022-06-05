package main

import (
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"time"
	"wxProjectDev/user/router"
	"wxProjectDev/work/controllers"
	"wxprojectApiGateway/service/registry"
)

var (
	g errgroup.Group
)

func init() {
	registry.InitRegistry("public/config/service.yaml")
}

func main() {
	server01 := &http.Server{
		Addr:         ":8000",
		Handler:      router.UserRouter(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	registry.Register("user", "127.0.0.1", "8000", time.Second*2)

	server02 := &http.Server{
		Addr:         ":8001",
		Handler:      controllers.Controller(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	registry.Register("work", "127.0.0.1", "8001", time.Second*2)

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
