# Microservices Store

![CI Pipeline](https://github.com/SurferAlex/microservices_store/workflows/CI%20Pipeline/badge.svg)

## Описание

Микросервисная архитектура для интернет-магазина.

## Сервисы

- **auth_service** - Сервис аутентификации и авторизации
- **profile_service** - Сервис управления профилями пользователей

## CI/CD

Проект использует GitHub Actions для автоматической проверки кода, сборки Docker образов и проверки миграций БД.

## Запуск

docker-compose up -d


### Что добавлено:

1. `workflow_dispatch` — возможность запускать workflow вручную из интерфейса GitHub
2. Summary для каждой задачи — краткий отчет о результатах
3. Badge в README — визуальный индикатор статуса CI/CD

### Проверка:

1. Сохраните обновленный файл `.github/workflows/ci.yml`
2. Создайте/обновите `README.md` с badge
3. Закоммитьте и запушьте:
   git add .github/workflows/ci.yml README.md
   git commit -m "Добавлены summary и badge статуса CI/CD"
   git push
   