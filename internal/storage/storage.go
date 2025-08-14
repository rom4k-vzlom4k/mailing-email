package storage

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rom4k-vzlom4k/mailing-email/internal/models"
)

type EmailRepository interface {
	AddEmail(ctx context.Context, email models.AddEmail) (int64, error)
	GetPendingEmails(ctx context.Context, before time.Time) ([]models.AddEmail, error)
	UpdateStatus(ctx context.Context, id int64, status models.SentStatus, sentAt ...*time.Time) error
}

type emailRepo struct {
	pool *pgxpool.Pool
}

func NewEmailRepository(pool *pgxpool.Pool) EmailRepository {
	return &emailRepo{pool: pool}
}

func (r *emailRepo) AddEmail(ctx context.Context, email models.AddEmail) (int64, error) {
	sql := `INSERT INTO email_queue (to_email, subject, body, scheduled_at, sent_at, status)
	        VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	row := r.pool.QueryRow(ctx, sql,
		email.ToEmail, email.Subject, email.Body, email.ScheduledAt, email.SentAt, email.Status)
	var id int64
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *emailRepo) GetPendingEmails(ctx context.Context, before time.Time) ([]models.AddEmail, error) {
	sql := `SELECT id, to_email, subject, body, scheduled_at, sent_at, status
	        FROM email_queue
	        WHERE status = 'pending' AND scheduled_at < $1`
	rows, err := r.pool.Query(ctx, sql, before)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []models.AddEmail
	for rows.Next() {
		var ae models.AddEmail
		if err := rows.Scan(&ae.ID, &ae.ToEmail, &ae.Subject, &ae.Body, &ae.ScheduledAt, &ae.SentAt, &ae.Status); err != nil {
			return nil, err
		}
		data = append(data, ae)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return data, nil
}

func (r *emailRepo) UpdateStatus(ctx context.Context, id int64, status models.SentStatus, sentAt ...*time.Time) error {
	var sa *time.Time
	if len(sentAt) > 0 {
		sa = sentAt[0]
	}
	sql := `UPDATE email_queue SET status = $1, sent_at = $2 WHERE id = $3`
	_, err := r.pool.Exec(ctx, sql, status, sa, id)
	return err
}
