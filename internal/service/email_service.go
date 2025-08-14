package service

import (
	"context"
	"fmt"
	"time"

	"github.com/rom4k-vzlom4k/mailing-email/internal/models"
	"github.com/rom4k-vzlom4k/mailing-email/internal/storage"
	"gopkg.in/gomail.v2"
)

type EmailService interface {
	AddEmail(ctx context.Context, email models.AddEmail) (int64, error)
	ProcessPendingEmails(ctx context.Context) error
}

type emailService struct {
	repo       storage.EmailRepository
	stmpConfig models.SMTPConfig
}

func NewEmailService(repo storage.EmailRepository, cfg models.SMTPConfig) EmailService {
	return &emailService{
		repo:       repo,
		stmpConfig: cfg,
	}
}

func (s *emailService) AddEmail(ctx context.Context, email models.AddEmail) (int64, error) {
	if email.Status == "" {
		email.Status = models.StatusPending
	}
	return s.repo.AddEmail(ctx, email)
}

func (s *emailService) ProcessPendingEmails(ctx context.Context) error {
	now := time.Now()
	emails, err := s.repo.GetPendingEmails(ctx, now)
	if err != nil {
		return fmt.Errorf("get pending emails %w", err)
	}
	for _, e := range emails {
		if err != s.repo.UpdateStatus(ctx, e.ID, models.StatusFailed) {
			fmt.Println("failed to set in_progress:", err)
			continue
		}
		if err := s.sendEmail(e); err != nil {
			fmt.Println("send email failed:", err)
			s.repo.UpdateStatus(ctx, e.ID, models.StatusFailed, nil)
			continue
		}
		sentAt := time.Now()
		if err := s.repo.UpdateStatus(ctx, e.ID, models.StatusDone, &sentAt); err != nil {
			fmt.Println("failed to set done:", err)
		}
	}
	return nil
}

func (s *emailService) sendEmail(email models.AddEmail) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.stmpConfig.From)
	m.SetHeader("To", email.ToEmail)
	m.SetHeader("Subject", email.Subject)
	m.SetBody("text/plain; charset=UTF-8", email.Body)
	d := gomail.NewDialer(s.stmpConfig.Host, s.stmpConfig.Port, s.stmpConfig.Username, s.stmpConfig.Password)
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil

}
