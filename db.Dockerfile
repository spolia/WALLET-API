FROM mysql:8.0.23
COPY ./migrations/init.sql /docker-entrypoint-initdb.d/