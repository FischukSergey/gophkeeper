FROM postgres:16-alpine

COPY build/db/init-db.sh /docker-entrypoint-initdb.d/
COPY build/db/ensure-db.sh /usr/local/bin/
RUN chmod +x /docker-entrypoint-initdb.d/init-db.sh \
    && chmod +x /usr/local/bin/ensure-db.sh

# Обернем оригинальную точку входа
COPY build/db/docker-entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/docker-entrypoint.sh

ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["postgres"]
