package models

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}
