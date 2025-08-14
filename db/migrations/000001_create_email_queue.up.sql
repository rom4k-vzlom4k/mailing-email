CREATE TABLE email_queue (
    id SERIAL PRIMARY KEY,
    to_email VARCHAR(255) NOT NULL,
    subject TEXT NOT NULL,
    body TEXT NOT NULL,
    scheduled_at TIMESTAMP NOT NULL,
    sent_at TIMESTAMP,
    status VARCHAR(50) DEFAULT 'pending',
    error TEXT
);
