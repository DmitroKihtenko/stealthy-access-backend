# Stealthy Access Backend
Access management backend service component of Stealthy application.
REST API web server. Encapsulates business logic of Stealthy application
for users entities management:
 - viewing list of application users
 - adding new users
 - deleting users

### Technologies
Build with:
 - Golang 1.21.1
 - Gin framework 1.9.1
 - Go Kratos Client 1.0.1

Works on HTTP web protocol.

### Requirements
Installed Docker and Docker-compose plugin

### API
You can view API details using Openapi standard mapping `/swagger/index.html`

### How to up and run
Configure application
1. Copy files: `kratos.env.example` to `kratos.env`, `config.yaml.example` to
`config.yaml`, `postgres.env.example` to `postgres.env`
2. Make changes you need in copied configuration files (details about
configs can be found in these files)

Build docker images
```bash
docker compose build
```

Build docker images and start service
```bash
docker compose up
```

Stop and remove containers after application use
```bash
docker compose down
```
