FROM oryd/kratos:v1.1.0 AS kratos_application
COPY build/kratos/ /home/ory/
USER root
RUN apk add gettext
RUN ["chmod", "a+x", "/home/ory/init-script.sh"]
USER ory
ENTRYPOINT ["/bin/sh", "-c" , "sh /home/ory/init-script.sh && kratos serve --config /home/ory/kratos.yml"]
CMD []

FROM postgres:14.3-alpine3.16 AS kratos_storage
COPY build/postgres/docker-entrypoint-initdb.d/ /docker-entrypoint-initdb.d/
USER root
RUN ["chmod", "a+wr", "-R", "/docker-entrypoint-initdb.d"]
RUN ["chmod", "a+x", "/docker-entrypoint-initdb.d/1-forward-vars.sh"]
USER postgres

FROM golang:1.21.1-alpine AS builder
WORKDIR /app
COPY ./src/ ./
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.2 && swag init \
&& go build -o /app/access-backend

FROM golang:1.21.1-alpine AS application
ARG UID=1001
RUN adduser -u $UID -D app-user
COPY --from=builder /app/access-backend /app/access-backend
WORKDIR /app
RUN chown -R app-user:app-user /app && chmod u+x access-backend
COPY ./config.yaml .
USER app-user
CMD ["/app/access-backend"]
