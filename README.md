# Тестовое задание L0

### Использование

Запустить контейнеры с базой данных и nats-streaming:

`docker-compose up --detach`

Прописать переменные окружения можно напрямую в yml-конфиге или создать файл .env вида:

```
DB_USER=user
DB_PASSWORD=password
DB_NAME=db
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=postgres
```

### Флаги запуска сервиса
- `addr` (по умолчанию ":8080")
- `dbURL` для подключения к БД (по умолчанию составляется из переменных окружения): "postgres://<user>:<password>"@localhost:5432/<db_name>"

### Nats-streaming
Приложение pub публикует три сообщения: валидный json, его дубликат и невалидный json. Сервис с подписчиком принимает эти сообщения, добавляет валидный заказ в БД и кэш, а о дубликате и невалидном json выдает уведомление
