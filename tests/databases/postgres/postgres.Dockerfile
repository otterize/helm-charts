FROM postgres:latest

ENV POSTGRES_DB=""
ENV POSTGRES_USER=""
ENV POSTGRES_PASSWORD=""

COPY postgres.sql /docker-entrypoint-initdb.d/database.sql

# Expose the PostgreSQL port
EXPOSE 5432