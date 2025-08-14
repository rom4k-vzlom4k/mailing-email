package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/rom4k-vzlom4k/mailing-email/internal/models"
	"github.com/rom4k-vzlom4k/mailing-email/internal/service"
	"github.com/rom4k-vzlom4k/mailing-email/internal/storage"
)

func main() {
	// подгружаем env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// прерываем контекст
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// коннект к бд
	dsn := os.Getenv("DATABASE_URL")
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("cannot connect to DB: %v", err)
	}
	defer pool.Close()

	// читаем конфиг SMTP
	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		log.Fatalf("invalid SMTP_PORT: %v", err)
	}

	smtpCfg := models.SMTPConfig{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     port,
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
		From:     os.Getenv("SMTP_FROM"),
	}

	repo := storage.NewEmailRepository(pool)
	svc := service.NewEmailService(repo, smtpCfg)

	// тикер для воркера
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	fmt.Println("Worker started")
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Shutting down worker")
			return
		case <-ticker.C:
			fmt.Println("Processing pending emails")
			if err := svc.ProcessPendingEmails(context.Background()); err != nil {
				fmt.Println("Error processing emails:", err)
			}
		}
	}
}
