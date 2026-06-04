package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/loanem-backend/inventory-service/config"
	"github.com/loanem-backend/inventory-service/internal/consumer"
	"github.com/loanem-backend/inventory-service/internal/server"
	"github.com/loanem-backend/inventory-service/internal/service"
	"google.golang.org/grpc"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	godotenv.Load()

	db := config.InitDB()
	defer db.Close()

	amqpConn, amqpCh := config.InitAMQPChannel()
	defer amqpCh.Close()
	defer amqpConn.Close()

	s := grpc.NewServer()

	is, ts, cs := service.Initialize(db)

	server.Start(s, is, ts)

	if err := consumer.Start(ctx, amqpCh, cs); err != nil {
		log.Fatalf("failed starting service: %v\n", err)
	}

	lis := config.InitListener()

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed serving grpc: %v\n", err)
	}

	<-ctx.Done()
	log.Println("Stopping server gracefully...")
}
