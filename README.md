# BL Wallet Service

BlueLabs Wallet Service, is a simple service to manage a digital wallet.

This service provides the following features

- Creates a wallet for a user
- Gets a wallet for a user
- Add funds to the wallet balance
- Remove funds from the wallet balance

⚠️ For the purpose of the exercise I've committed the `.env` file

Also, the service inits with some users already in database, you can use their ids to do wallet actions.

| ID  | Name  |
|-----|:-----:|
| 1   | Jane  |
| 2   | John  |
| 3   | Peter |

## Usage

#### Prerequisites

Before running the API, make sure you have Docker and Docker Compose installed on your machine.

## Usage

#### Prerequisites

Before running the API, make sure you have Docker and Docker Compose installed on your machine.

### Running the service

```bash
make build
```

Wait for the service to start. You will see log messages indicating the service is running.

Once the service is up and running, you can access it at http://localhost:9000

#### Stopping the service

To stop the service, press `Ctrl + C` in the terminal where the service is running.
You can also run `docker-compose down`

### Running Unit Tests

To execute tests for the API, follow these steps:

Run the following command:

```shell
make test
```

### API Documentation

Once the API is running, you can access the API documentation at http://localhost:9000/swagger/index.html.

## Service Calls and Responses

Local Endpoint: `http://localhost:9000`

The following endpoints are available in the app:

### Health Check

- Method: GET
- Path: `/health`
- Description: Returns the status of the service.

**Curl Request**

```shell
curl --location 'http://localhost:9000/health'
```

**Example Response:**

```json
{
  "status": "up"
}
```

### Get A Wallet by UserID

- Method: GET
- Path: `/users/:id/wallet`
- Description: Retrieves a wallet by userID.

**Curl Request**

```shell
curl --location 'http://localhost:9000/users/1/wallet'
```

**Example Response:**

```json
{
  "id": "25a08945-1caf-4734-bd12-bb6cf66dae52",
  "userID": "1",
  "balance": 15.0
}
```

### Create a Wallet

- Method: POST
- Path: `users/wallet`
- Description: Creates a new wallet for a certain user.

**Curl Request**

```shell
curl --location 'http://localhost:9000/users/wallet' \
--header 'Content-Type: application/json' \
--data '{
    "userID": "1"
}'
```

**Example Request Body:**

```json
{
  "userID": "1"
}
```

**Example Response Body:**

```json
{
  "message": "wallet created",
  "userID": "1"
}
```

### Add funds to wallet

- Method: POST
- Path: `users/wallet/add`
- Description: Adds fund to a user wallet.

**Curl Request**

```shell
curl --location 'http://localhost:9000/users/wallet/add' \
  --header 'x-idempotency-key: 1234' \
  --header 'Content-Type: application/json' \
  --data '{
    "userID": "1",
    "amount": 20.0
}'

```

**Example Request Body:**

```json
{
  "userID": "1",
  "amount": 10.0
}
```

**Example Response Body:**

```json
{
  "message": "funds were added",
  "userID": "1"
}
```

### Remove funds from wallet

- Method: POST
- Path: `users/wallet/remove`
- Description: Removes fund to a user wallet.

**Curl Request**

```shell
curl --location 'http://localhost:9000/users/wallet/remove' \
--header 'x-idempotency-key: 1234' \
--header 'Content-Type: application/json' \
--data '{
    "userID": "1",
    "amount": 10.0
}'
```

**Example Request Body:**

```json
{
  "userID": "1",
  "amount": 10.0
}
```

**Example Response Body:**

```json
{
  "message": "funds were removed",
  "userID": "1"
}
```
