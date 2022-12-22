# mafia_bot

Проект дискорд бота для парсинга статистики с [полемики](https://polemicagame.com/)

## Запуск dev контейнера

Перед запуском необходимо настроить переменные окружения\
`docker-compose -f ./deployments/docker-compose.yml -f ./deployments/docker-compose-dev.yml up --build`

## Запуск тестов

`go test ./test/...`\
Для тестов тестов БД сперва необходимо войти в контейнер
`docker exec -it deployments-mafia-bot-1 /bin/bash`

## Переменные окружения бота

### Переменные окружения для бота

- MAFIA_BOT_DB_PASSWORD - Пароль для подключения к БД
- MAFIA_BOT_DB_USER - Пользователь для подключения к БД
- MAFIA_BOT_DB_NAME - Название БД бота
- MAFIA_BOT_DB_HOST - Название хоста БД, по умолчанию для контейнера `mafia-db`
- MAFIA_BOT_DISCORD_TOKEN - Токен дискорд бота
- MAFIA_BOT_STATUS_CHANNELS - Перечисленные через запятую ID каналов куда бот будет отправлять диагностические
  сообщения (например что бот активен)
