# Mailing Email Service 

Сервис для планирования и отправки email сообщений по расписанию на Go

## Статус проекта
**В разработке** — в данный момент умеет отправлять письма через SMTP (Gmail, Yandex и др.) и хранить очередь в PostgreSQL
Планы:
- Панель управления для просмотра и редактирования очереди писем
- Логи отправки и повторные попытки при ошибках

## Структура проекта
```
├── 📁 cmd/
│   └── 📁 mailing/
│       └── 🐹 main.go
├── 📁 db/
│   └── 📁 migrations/
│       ├── 🗄️ 000001_create_email_queue.down.sql
│       └── 🗄️ 000001_create_email_queue.up.sql
├── 📁 internal/
│   ├── 📁 models/
│   │   ├── 🐹 email.go
│   │   └── 🐹 smtp.go
│   ├── 📁 service/
│   │   └── 🐹 email_service.go
│   └── 📁 storage/
│       └── 🐹 storage.go
├── 📄 .env.example
├── 🚫 .gitignore
├── 🐳 Dockerfile
├── 📖 README.md
├── ⚙️ docker-compose.yml
├── 🐹 go.mod
└── 🐹 go.sum
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

## Запуск через Docker

1. Убедитесь, что Docker и Docker Compose установлены.

2. Скопируйте `.env.example` в `.env` и заполните переменные.

3. Поднимите контейнеры:
```bash
docker-compose up --build -d
```
4. Проверить логи:
```bash
docker-compose logs -f app
```
### Сервис будет работать, база создат таблицу email_queue автоматически через миграцию.

В данный момент реализовывается таким образом:  
Подключитесь к базе и вставьте несколько писем.
```bash
docker exec -it mailing_db psql -U mailing_user -d mailing_db
```
Письмо:
```sql
INSERT INTO email_queue (to_email, subject, body, scheduled_at, status)
VALUES 
('test1@example.com', 'Тестовое письмо', 'Это тест', now(), 'pending');
```
Сервис автоматически обработает письмо со статусом `pending`

## Логи
Логи Go приложения
```bash
docker-compose logs -f app
```
Логи базы
```bash
docker-compose logs -f db
```
