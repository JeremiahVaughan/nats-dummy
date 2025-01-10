package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

var ns *server.Server

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := initNats()
	if err != nil {
		log.Fatalf("error, when initNats() for main(). Error: %v", err)
	}
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		<-sigs
		cancel()
	}()

	// Replace the infinite select with a listen to the context
	select {
	case <-ctx.Done():
		log.Println("Received interrupt, shutting down...")
	}

	log.Println("Exiting...")
}

func initNats() error {
	opts := &server.Options{
		Port: 3000,
		// You can set other options here (e.g., authentication, clustering)
	}

	var err error
	ns, err = server.NewServer(opts)
	if err != nil {
		return fmt.Errorf("error, when creating NATS server. Error: %v", err)
	}

	go ns.Start()

	// Ensure the server has started
	if !ns.ReadyForConnections(10 * time.Second) {
		return errors.New("error, NATS Server didn't start in time")
	}

	// Retrieve the server's listen address
	addr := ns.Addr()
	var port int
	if tcpAddr, ok := addr.(*net.TCPAddr); ok {
		port = tcpAddr.Port
	} else {
		return fmt.Errorf("error, filed to get nats server port")
	}
	fmt.Printf("NATS server is running on port %d\n", port)
	return nil
}

func connectToNats() (*nats.Conn, error) {
	nc, err := nats.Connect(ns.ClientURL())
	if err != nil {
		return nil, fmt.Errorf("error, when connecting to NATS server. Error: %v", err)
	}
	return nc, nil
}
