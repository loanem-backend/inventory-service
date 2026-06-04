package config

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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
