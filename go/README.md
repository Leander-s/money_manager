# Go backend

## Requirements
- Go
- docker-compose

## Setup
The repository has a .env.example file that contains the environment variables needed to run the backend server. Copy this file to a new file named .env and update the values as needed:
```bash
cp .env.example .env
```

## Run
To start the backend and database servers run:
```bash
docker-compose up --build
```

to stop the servers run:
```bash
docker-compose down
```
