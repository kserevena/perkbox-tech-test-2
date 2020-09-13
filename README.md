# perkbox-tech-test-2

# Prerequisites
A recent version of docker (with docker-compose)

## Execution
Start the database with the following docker compose command:

`docker-compose up -d`

Stop the database with the following command:

`docker-compose down`

## Production settings
This app connects to a MongoDB database. The environment variable `MONGODB_CONNECTION_STRING`
should be set to the required connection string when running in a production environment.

## Interface mocking
For the benefit of unit testing, interfaces have been used in this application which can be mocked
during tests.

Mockery (see https://github.com/vektra/mockery) has been used to generate the mocks in this repository.