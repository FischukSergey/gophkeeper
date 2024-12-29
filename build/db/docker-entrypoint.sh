#!/bin/bash
set -e

# Запускаем оригинальный entrypoint PostgreSQL в фоновом режиме
docker-entrypoint.sh postgres &

# Ждем некоторое время, чтобы PostgreSQL успел запуститься
sleep 5

# Выполняем наш скрипт
ensure-db.sh

# Ждем завершения основного процесса
wait $! 