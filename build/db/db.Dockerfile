FROM postgres:16-alpine

# Аргументы для настройки базы данных (можно переопределить при сборке)
ARG POSTGRES_DB=gophkeeper
ARG POSTGRES_USER=postgres
ARG POSTGRES_PASSWORD=postgres

# Установка переменных окружения
ENV POSTGRES_DB=$POSTGRES_DB
ENV POSTGRES_USER=$POSTGRES_USER
ENV POSTGRES_PASSWORD=$POSTGRES_PASSWORD

# Копирование всех SQL миграций
COPY migration/*.sql /docker-entrypoint-initdb.d/

# Настройка прав доступа к директории с данными
RUN chmod 0700 /var/lib/postgresql/data

# Открываем порт PostgreSQL
EXPOSE 5432

# Точка монтирования для данных
VOLUME ["/var/lib/postgresql/data"]
