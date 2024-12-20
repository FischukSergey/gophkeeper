FROM postgres:16-alpine

# Аргументы для настройки базы данных (не чувствительные данные)
ARG POSTGRES_DB=gophkeeper
ARG POSTGRES_USER=postgres

# Установка переменных окружения (пароль будет передан при запуске)
ENV POSTGRES_DB=$POSTGRES_DB \
    POSTGRES_USER=$POSTGRES_USER

# Копирование всех SQL миграций из корректного пути
COPY ./migrations/*.sql /docker-entrypoint-initdb.d/

# Настройка прав доступа к директории с данными
RUN chmod 0700 /var/lib/postgresql/data

# Открываем порт PostgreSQL
EXPOSE 5432

# Точка монтирования для данных
VOLUME ["/var/lib/postgresql/data"]
