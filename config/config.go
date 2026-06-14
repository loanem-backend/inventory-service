package config

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	s3config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/loanem-backend/api-gateway/pkg/storage"
	amqp "github.com/rabbitmq/amqp091-go"
)

func InitDB() *pgxpool.Pool {
	ctx := context.Background()

	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		panic(fmt.Errorf("failed parsing databse config: %w", err))
	}
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeExec

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		panic(fmt.Errorf("failed creating connection pool: %w", err))
	}

	return pool
}

func InitStorageClient() *storage.S3Client {
	cfg, err := s3config.LoadDefaultConfig(context.Background(),
		s3config.WithRegion(os.Getenv("STORAGE_REGION")),
	)
	if err != nil {
		panic(fmt.Errorf("failed setting up storage client: %w", err))
	}

	client := s3.NewFromConfig(cfg)

	return &storage.S3Client{
		Client:        client,
		PresignClient: s3.NewPresignClient(client),
		Bucket:        os.Getenv("STORAGE_BUCKET"),
	}
}

func InitListener() net.Listener {
	listener, err := net.Listen("tcp", ":"+os.Getenv("APP_PORT"))
	if err != nil {
		panic(fmt.Errorf("failed listening: %w", err))
	}

	return listener
}

func InitAMQPChannel() (*amqp.Connection, *amqp.Channel) {
	url := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		os.Getenv("RABBITMQ_USER"),
		os.Getenv("RABBITMQ_PASS"),
		os.Getenv("RABBITMQ_HOST"),
		os.Getenv("RABBITMQ_PORT"),
	)

	var (
		conn *amqp.Connection
		err  error
	)

	for i := 1; i <= 5; i++ {
		conn, err = amqp.Dial(url)
		if err == nil {
			break
		}

		fmt.Printf("failed connecting to RabbitMQ -> attempt %d\n", i+1)
		time.Sleep(time.Second)
	}
	if err != nil {
		panic(fmt.Errorf("failed connecting to RabbitMQ: %w", err))
	}

	ch, err := conn.Channel()
	if err != nil {
		panic(fmt.Errorf("failed opening channel: %w", err))
	}

	return conn, ch
}
