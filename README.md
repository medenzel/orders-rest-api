
# GO orders REST API

Simple pet-project rest api for working with orders.
Stores orders in database(Postgres) and is able to get/update/delete them via HTTP requests.


## Prerequisites
Go 1.18

Docker
## Installation & Run

To build and run use docker compose command from root of the repo(runs locally on port 8080)
```bash
  docker compose up -d
```
    
## API

#### Get all orders
Requiered query parameters: page_id (min = 1), page_size (min = 5, max = 10)

```bash
  GET /api/v1/orders?page_id=X&page_size=Y
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


