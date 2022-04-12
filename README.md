# Wallet-API

![technology Go](https://img.shields.io/badge/technology-go-blue.svg)

## Overview

This project was implemented in golang, and simulates an virtual wallet allowing users operate with multiple currencies (`USDT`, `BTC`, `ARS`).

## APP structure

This application uses `package oriented design` as project structure, and also uses `gorilla/mux` for request router and dispatcher, `gorilla/sessions` to manage the users session and `mysql` as database.

The service is composed of two main components:

- `user` : manage users creations and credentials, here the repository is defined.
- `movements`: manage the user account transactions like send money to other account, deposits, account history, account balance, etc; here the repository is defined.

## Endpoints

- `POST /login` : APP Login with alias and password.
- `POST /logout` : APP Logout.
- `POST /users` : User registration. Users with the same alias nor the same email are not allowed. Every time a new user is registered, all accounts for each currency are also initialized.
- `GET /internal/movements/balance` : Get the balance for each user currency.
- `GET /internal/movements/history` : Get the transactions history for each user currency.
- `POST /internal/movements/send` : Send money to other user.
- `POST /internal/movements/deposit` : Deposit money deposit in own account.


## How To Run This Project

- Make sure you have already installed both Docker Engine and Docker Compose in the last version (Engine: 20.10.13 and Compose: v2.3.3).
- Make sure you have these variables set in your environment `export DOCKER_BUILDKIT=0` and `export COMPOSE_DOCKER_CLI_BUILD=0`.
- Type `make build` to build the docker compose and then `make up` to up the compose.

# Test

- Type `make test` to run the unit tests.
- You can find test cases to test the endpoints in the directory `cmd/api/internal/testdata` .


## Improvements

There are several changes and improvements to be made:

- Improve the unit test cases to cover all border cases
- Use a container to test against a real database instead use mock
- Create and use environment variables
