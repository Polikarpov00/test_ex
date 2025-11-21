# go-qa-api

Вопрос-ответ сервис на Go.

## Запуск через docker compose:

docker compose up --build

## API

- GET /questions/
- POST /questions/
- GET /questions/{id}
- DELETE /questions/{id}
- POST /questions/{id}/answers/
- GET /answers/{id}
- DELETE /answers/{id}

## Тесты
go test ./...
