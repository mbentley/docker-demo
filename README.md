# mbentley/docker-demo

Docker Demo Application

Original Author: Evan Hazlett ([https://github.com/ehazlett/docker-demo](https://github.com/ehazlett/docker-demo))

This is a Go demo application with a PostgreSQL database used for demonstrating Docker and Docker Compose

## Demo

### Requirements

- Install Docker ([https://docs.docker.com/installation/](https://docs.docker.com/installation/)) — Docker Compose is included as a plugin with Docker Engine and Docker Desktop

### Build

- `docker compose build`

### Run

- `docker compose up`

### Cleanup

- `docker compose down -v --rmi local`
