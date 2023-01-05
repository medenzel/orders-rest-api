
# GO orders REST API

Simple pet-project rest api for working with orders.
Stores orders in database(Postgres) and is able to get/update/delete them via HTTP requests.


## Prerequisites
Go 1.18

Docker
## Installation & Run

To download this app use command

```bash
  go get github.com/medenzel/orders-rest-api
```

To build and run use docker compose command from root of the repo(runs locally on port 8080)
```bash
  docker compose up -d
```
    
## API

#### Get all orders

```bash
  GET /api/v1/orders
```
#### Get order

```bash
  GET /api/v1/orders/{id}
```
#### Post new order

```bash
  POST /api/v1/orders
```
#### Update order

```bash
  PUT /api/v1/orders/{id}
```
#### Delete order

```bash
  DELETE /api/v1/orders/{id}
```


