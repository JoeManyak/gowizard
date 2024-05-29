package consts

const DefaultPostgresDockerfile = `version: '3.8'

services:
  postgres:
    image: postgres:latest
    container_name: postgres
    environment:
      POSTGRES_DB: default
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
`
