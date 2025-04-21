# 1. The project used:

* Standard CRUD operations of a database table
* JWT-based authentication
* Environment dependent application configuration management
* Structured logging with contextual information
* Error handling with proper error response generation
* Database migration
* Data validation
* Containerizing an application in Docker

# 2. The kit uses the following Go packages:

* RPC framework: [grpc](google.golang.org/grpc)
* Database access: [pgx](https://github.com/jackc/pgx)
* Database migration: [golang-migrate](https://github.com/golang-migrate/migrate)
* Data validation: [go-playground validator](https://github.com/go-playground/validator)
* Logging: [log/slog](https://pkg.go.dev/golang.org/x/exp/slog)
* JWT: [jwt-go](https://github.com/dgrijalva/jwt-go)
* Config reader: [cleanv](github.com/ilyakaznacheev/cleanenv)
* Env reader: [godotenv](github.com/joho/godotenv)

# 3. Getting Started

After installing Go, Docker and TaskFile, run the following commands to start experiencing:

```shell
## RUN SSO
# download the project
git clone https://github.com/cuongpiger/opswat.git
cd sso

# create config.env with that text:
$ nano config.env {
CONFIG_PATH=./config/local.yaml
POSTGRES_DB=url
POSTGRES_USER=myuser
POSTGRES_PASSWORD=mypass
}

# start a PostgreSQL database server in a Docker container
task db-start

# run the SSO server
go run ./cmd/sso
```

Also, you can start project in dev mode. For that you need rename in config.env
"CONFIG_PATH=./config/local.yaml" to "CONFIG_PATH=./config/dev.yaml" in both projects
and run following commads:

```shell
# run the SSO server
cd sso/
docker compose up --build
```

SSO-grpc Server running at http://localhost:44044. The server provides the following endpoints:

## 3.1 auth

* `Regiter`: register new user in db
* `Login`: log in to the application
* `GETUserID`: get user ID by name

## 3.2 permissions

* `SetAdmin`: set exists user to admin in your app. You need be creator of app
* `DelAdmin`: delete exists user from admin in your app. You need be creator of app
* `IsAdmin`: is the user an admin by userID
* `IsCreator`: is the user a creator by userID

## 3.3 apps

* `SetApp`: set new app in db. You will be creator of the app
* `DelApp`: delete exists apps. You need be creator of app
* `UpdApp`: update app name and secret
* `GetAppID`: get app id by app name

# 4. Project Layout

Project has the following project layout:

```
opswat/
├── cmd/                       start of applications of the project
├── config/                    configuration files for different environments
├── deployment/                configuration for create daemon in linux
├── internal/                  private application and library code
│   ├── app/                   application assembly
│   ├── config/                configuration library
│   ├── domain/                models of apps and users
│   ├── grpc/                  grpc handlers
│   │   ├── apps/              handlers of apps
│   │   ├── auth/              handlers of auth
│   │   └── permissions/       handlers of permissions
│   ├── lib/                   additional functions for logging, error handling, migration
│   ├── services/              logics of handlers
│   │   ├── apps/              handlers of apps
│   │   ├── auth/              handlers of auth
│   │   └── permissions/       handlers of permissions
│   └── storage/               storage library
├── migrations/                migrations
└── config.env                 config for sercret variables
```

The top level directories `cmd`, `internal`, `lib` are commonly found in other popular Go projects, as explained in
[Standard Go Project Layout](https://github.com/golang-standards/project-layout).

Within each feature package, code are organized in layers (grpc server, service, db), following the dependency
guidelines
as described in the [clean architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html).
