package main

import (
	"context"
	"github.com/vastzp/alerts/server"
	"github.com/vastzp/alerts/service"
	"github.com/vastzp/alerts/storage"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {

	// I implemented this task with 3 layers: storage, service, server.

	// this is initialization of storage layer
	storageFilename := "alerts.db"
	sqliteStorage, err := storage.NewSQLiteStorage(storageFilename)
	if err != nil {
		log.Fatal("failed to init storage")
		return
	}

	// this is initialization of service layer
	serviceInstance := service.NewService(sqliteStorage)

	// this is initialization of server layer
	serverInstance := server.NewServer(serviceInstance)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Run server instance
	go serverInstance.Run()

	// Create link with os signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// wait for os signal
	<-c

	// do shutting down routine
	log.Println("shutting down")

	// just gives 3 seconds to server to execute current requests
	ctx, cancel = context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := serverInstance.Shutdown(ctx); err != nil {
		log.Printf("error during server shutdown: %s\n", err)
	} else {
		log.Println("server shutdown complete.")
	}
}
