## Commander - REST API for running commands

This application provides a REST API for managing and executing commands in the form of bash scripts.


### Configuration
An example configuration file has been "example.env". Copy this file to the same directory (config) and in the new ".env" enter the credentials needed to connect to your database.

### REST API
- POST request to /commands --> Create command
- GET request to /commands --> Get all commands
- GET request to /commands/1 --> Get command with ID 1
- DELETE request to /commands/1 --> Delete command with ID 1


### Run
Be sure to add the configuration file before starting.

#### Docker
``` bash
docker compose up
```