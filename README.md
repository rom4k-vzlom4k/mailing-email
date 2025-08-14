# Mailing Email Service 

Сервис для планирования и отправки email сообщений по расписанию на Go

## Статус проекта
**В разработке** — в данный момент умеет отправлять письма через SMTP (Gmail, Yandex и др.) и хранить очередь в PostgreSQL
Планы:
- Панель управления для просмотра и редактирования очереди писем
- Логи отправки и повторные попытки при ошибках
- Docker-compose для быстрого запуска

## Структура проекта
```
mailing_email/
├── cmd/
│ └── main.go
├── db/
│ └── migrations/
│ ├── 000001_create_email_queue.down.sql # SQL для отката миграции
│ └── 000001_create_email_queue.up.sql # SQL для создания таблицы email_queue
├── internal/
│ ├── models/
│ │ ├── email.go # Модель Email
│ │ └── smtp.go # Модель SMTPConfig
│ ├── service/
│ │ └── email_service.go # Логика отправки писем
│ └── storage/
│ └── storage.go # Работа с БД
├── .env.example # Пример конфигурации окружения
├── go.mod
├── go.sum
└── README.md
```

## Установка и запуск

1. Клонируйте репозиторий:
```bash
git clone https://github.com/username/mailing_email.git
cd mailing_email
```

2. Создайте .env на основе .env.example и заполните:
```env
DATABASE_URL=postgres://user:password@localhost:5432/dbname
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your_email@gmail.com
SMTP_PASSWORD=app_password
SMTP_FROM=your_email@gmail.com
```

3. Запустите миграции:

```bash
migrate -path db/migrations -database "$DATABASE_URL" up
```

4. Запустите сервис:
```bash
go run cmd/main.go
```