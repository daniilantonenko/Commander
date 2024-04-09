## Commander - REST API for running commands

This application provides a REST API for managing and executing commands in the form of bash scripts.


### Configuration
An example configuration file has been "example.env". Edit if other application settings are needed and rename to ".env".

### REST API
- POST request to /commands --> Create command
- GET request to /commands --> Get all commands
- GET request to /commands/1 --> Get command with ID 1
- DELETE request to /commands/1 --> Delete command with ID 1


### Run
Be sure to add the configuration file before starting.

#### Release
``` bash
tar -czf [archive_name]
cd [app_path]
mv example.env .env
nano .env
./commander
```

#### Golang run
``` bash
go run ./cmd/commander/main.cmd
```

#### Docker
``` bash
docker compose up
```