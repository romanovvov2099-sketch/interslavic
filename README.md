# DRIVER SERVER 

HTTP Server to interaction and get metrics with KVANT Drivers

Server abilities:
 - Create, Modificate, Get clients
 - Create, Modificate, Get drivers 
 - Get license requests from drivers
 - Approve requests and Send licenses
 - Update version of drivers
 - Update configs of drivers
 - Regenerate driver requests and resend driver licenses
 - Get Logs (Error, Critical) from drivers
 - Get Logs (Concrete logs like 08092024.log) from drivers
 - Get Configs (drivercore.yaml, handler_config.yaml, models_config.yaml) from drivers
 - Automation update new versions from gitlab drivers projects 


To use admin features need auth_token

To interaction with driver need send driver uuid from driver(locate in drivercore.yaml)


More information in Dependencies:
- [DRIVER_LICENSE](http://gitlab.e-m-l.ru/drivers/driver_license) - To generate request and licenses
- [DRIVER_API]("http://gitlab.e-m-l.ru/drivers/driverapi") - Driver API to interaction with concrete drivers

## DOCS

Swagger documentation to API on endpoint host/swagger

## DB Migrations 

Download and add to path environment:
[Golang migrate CLI](https://github.com/golang-migrate/migrate/releases) 

After schema changes you need to create migration, dev database must be running `docker compose up db` and configured at `config.yaml`

#### How to handle migrations

You can create migration with
```sh
migrate create -ext sql -dir db/migrations -seq migration_name
```

You can apply your migrations on dev database

```sh
migrate -source file://db/migrations -database postgres://drivers:drivers@localhost:5438/drivers?sslmode=disable up
```

To migrate back to previous one:
```sh
migrate -source file://db/migrations -database postgres://drivers:drivers@localhost:5438/drivers?sslmode=disable down 1
```

## Usage

- Docs generate with Swagger: 
```sh
swag init -g .\cmd\app\main.go
```

- App start: 
```sh 
go run .\cmd\app\main.go
```

- App start docker: 
```sh 
docker compose up --build
```

## Testing

- If you want start db only with needed migrates:
```sh
docker-compose -f .\tests\database\docker-compose.dbtest.yml up --build
```

## Configs

App config
```yaml
http:
  port: 5005
  authorization: "test_auth"

database: 
  url: "postgres://drivers:drivers@localhost:5438/drivers?sslmode=disable"

logging:
  dir: "./logs" # MUST BE STRING, PATH TO SAVE LOG
  enable: true # MUST BE BOOL.   
  level: "DEBUG" # MUST BE STRING. Available levels: INFO, WARNING, ERROR. DEFAULT: INFO
  format: "LOG" # MUST BE STRING. Available formats: TXT, JSON. DEFAULT: TXT
  saving_days: 10 # MUST BE INTEGER AND MINIMAL VALUE IS 1.

drivers_logs:
  dir: "./drivers/logs" # MUST BE STRING, PATH TO SAVES LOG

gitlab_token: "token" # TOKEN TO GET REPO FROM KVANT GITLAB
drivers_path: "./applications" # APPLICATIONS OF DRIVERS
drivers: # DRIVERS WHERE GET NEW VERSIONS EXE.
  - name: "universal_driver" # NAME OF DRIVER (MUST BE CONSTANT)
    gitlab_url: "http://gitlab.e-m-l.ru/drivers/universal_driver"  # GITLAB OF DRIVER REPO
    deploy: "http://gitlab.e-m-l.ru/api/v4/projects/175/repository/files/deploy_$V%2Funiversal_driver_$V.exe/raw?ref=deploy" # URL TEMPLATE TO DONWLOAD NEW EXE ($V - version)

```

## Project structure

```yaml
.
├─ .docker          - package to deploy with docker in environments (prod, dev) 
├─ app              - package to inject dependencies with logic from project 
├─ cmd              - package to start application
├─ config           - package for config
├─ db
│   └─ migration    - package for migration of DB
├─ docs             - swagger docs
├─ internal         - realisation of hanlder that getting message from device, identify and send to lis
│   ├─ cron         - package for do some actions with shedule cron
│   ├─ database     - package for interaction with DB (repositories, postgres)
│   ├─ http         - package for interaction with http (controllers, middlewares, framework - fiber)
│   ├─ models       - package for interaction with models (driver, client, request...)
│   ├─ usecases     - internal logic of endpoints and sheduler 
│   └─ utils        - some other logic 
├─ logging          - package for logging 
├─ tests            - package for test (run with manualy)
└─ pkg/client       - package with client to provide in other gitlab project to use driver-server api
```
