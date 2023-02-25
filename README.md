# mafia_bot

Проект дискорд бота для парсинга статистики с [полемики](https://polemicagame.com/)

## Запуск dev контейнера

Перед запуском необходимо настроить переменные окружения\
`docker-compose -f ./deployments/docker-compose.yml -f ./deployments/docker-compose-dev.yml up --build`

## Запуск тестов

`go test ./test/...`\
Для тестов БД сперва необходимо войти в контейнер  
`docker exec -it deployments-mafia-bot-1 /bin/bash`

## Переменные окружения бота

### Переменные окружения для бота

- MAFIA_BOT_DB_PASSWORD - Пароль для подключения к БД
- MAFIA_BOT_DB_USER - Пользователь для подключения к БД
- MAFIA_BOT_DB_NAME - Название БД бота
- MAFIA_BOT_DB_HOST - Название хоста БД, по умолчанию для контейнера `mafia-db`

#### Параметры подключения к дискорд API

- MAFIA_BOT_DISCORD_TOKEN - Токен дискорд бота
- MAFIA_BOT_STATUS_CHANNELS - Перечисленные через запятую ID каналов куда бот будет отправлять диагностические
  сообщения (например что бот активен)

#### Параметры подключения к API (взяты путем реверс инженеринга)

- MAFIA_BOT_CSRF - csrf токен берется из лога общения после логина
- MAFIA_BOT_CSRF_COOKIE - csrf cockie берется из лога обмена после обмена
- MAFIA_BOT_POLEMICA_HOST - хост polemica
- MAFIA_BOT_POLEMICA_LOGIN - логин на polemica
- MAFIA_BOT_POLEMICA_PASSWORD - пароль на polemica

#### Параметры для переодически задач

- MAFIA_BOT_POLEMICA_PARSE_HISTORY_TASK_DELAY - интервал получения новых игр, по умолчанию раз в час
