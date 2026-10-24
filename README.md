# reward-point

A small Go HTTP service that manages a **rewards points ledger**. Customers
(members) earn points through activities like purchases, referrals and cashback,
and can later redeem them. Every activity is stored as a reward entry, and a
member's balance is the sum of all their entries.

It uses only the Go standard library: `net/http` routing, `log/slog` for
logging, and an in-memory store with no external dependencies.

---

## Contents

- [Quick start](#quick-start)
- [Point types](#point-types)
- [API](#api)
- [Validation rules](#validation-rules)
- [Project structure](#project-structure)
- [Testing](#testing)

---

## Quick start

### Prerequisites

- **Go 1.26+** to run locally, **or**
- **Docker** with Docker Compose

### Run with Docker Compose

```bash
docker compose up --build
```

### Run locally

```bash
go run ./cmd/reward-point
```

Either way, the service listens on **`http://localhost:8080`**.

### Try it out

```bash
# 1. Create a member
curl -s -X POST localhost:8080/members \
  -d '{"name":"Ron","email":"ron@xyz.com"}'

# 2. Earn points (purchase, referral, cashback)
curl -s -X POST localhost:8080/rewards \
  -d '{"member_id":1,"point_type_id":1,"points":500,"description":"Purchase at Store A"}'
curl -s -X POST localhost:8080/rewards \
  -d '{"member_id":1,"point_type_id":2,"points":200,"description":"Referred user Jack"}'
curl -s -X POST localhost:8080/rewards \
  -d '{"member_id":1,"point_type_id":3,"points":50,"description":"Cashback on payment of $132"}'

# 3. Redeem points (only accepts positive value)
curl -s -X POST localhost:8080/rewards \
  -d '{"member_id":1,"point_type_id":4,"points":300,"description":"Redeemed for e gift card"}'

# 4. Check the balance: 500 + 200 + 50 - 300 = 450
curl -s localhost:8080/members/1

# 5. List the member's reward entries
curl -s localhost:8080/members/1/rewards
```

Stop the service with `Ctrl+C` (or `docker compose down`).

---

## Point types

### Point acquire and redemption details

| `point_type_id` | Description      | Direction        |
|-----------------|------------------|------------------|
| 1               | Purchase Earning | Credit (+)       |
| 2               | Referral Bonus   | Credit (+)       |
| 3               | Cashback         | Credit (+)       |
| 4               | Redemption       | Debit (−)        |

### Business rules

- Purchase Earning, Referral Bonus and Cashback entries are stored with
  **positive** points.
- Redemption entries are stored with **negative** points.
- Clients always send `points` as a positive number. The service applies the
  sign based on `point_type_id`.
- A member's available balance is the **sum of all their reward entries**.

---

## API

All requests and responses use JSON. Errors have this shape:

```json
{ "error": "human readable message" }
```

| Method | Path                      | Description                               |
|--------|---------------------------|-------------------------------------------|
| POST   | `/members`                | Create a member                           |
| GET    | `/members/{id}`           | Get a member and their points balance     |
| POST   | `/rewards`                | Create a reward entry                     |
| GET    | `/members/{id}/rewards`   | List all reward entries for a member      |

### `POST /members`

Creates a member account. The service sets `member_id` and `created_at`.

**Request**

```json
{
  "name": "Alice Johnson",
  "email": "alice@example.com"
}
```

**Response `200 OK`**

```json
{
  "name": "Alice Johnson",
  "email": "alice@example.com",
  "member_id": 1,
  "created_at": "2024-01-15T09:00:00Z"
}
```

### `GET /members/{id}`

Returns the member's details and their current points balance.

**Response `200 OK`**

```json
{
  "name": "Alice Johnson",
  "email": "alice@example.com",
  "member_id": 1,
  "created_at": "2024-01-15T09:00:00Z",
  "point_balance": 450
}
```

### `POST /rewards`

Creates a reward entry for a member. The service sets `reward_id` and
`event_date`.

**Request**

```json
{
  "member_id": 1,
  "point_type_id": 1,
  "points": 500,
  "description": "Purchase at Store A"
}
```

**Response `200 OK`**

```json
{
  "reward_id": 1,
  "member_id": 1,
  "point_type_id": 1,
  "points": 500,
  "description": "Purchase at Store A",
  "event_date": "2024-02-01T14:22:10Z"
}
```

For a redemption (`point_type_id: 4`), the response shows the stored negative
value, for example `"points": -300`.

### `GET /members/{id}/rewards`

Returns every reward entry for the member in the order they were created. A
member with no entries gets an empty array `[]`.

**Response `200 OK`**

```json
[
  {
    "reward_id": 1,
    "member_id": 1,
    "point_type_id": 1,
    "points": 500,
    "description": "Purchase at Store A",
    "event_date": "2024-02-01T14:22:10Z"
  },
  {
    "reward_id": 2,
    "member_id": 1,
    "point_type_id": 2,
    "points": 200,
    "description": "Referred user Bob",
    "event_date": "2024-02-05T10:15:33Z"
  }
]
```

### Status codes

| Status | When                                                              |
|--------|-------------------------------------------------------------------|
| 200    | Success                                                           |
| 400    | Request body is not valid JSON, or `{id}` in the path is not a number |
| 500    | The store rejected the request (see [Validation rules](#validation-rules)) |

---

## Validation rules

The store layer ([internal/database/store.go](internal/database/store.go))
enforces these rules. Each has its own exported error.

| Rule                                                        | Error                   |
|-------------------------------------------------------------|-------------------------|
| A member cannot be created with a duplicate email (case and surrounding whitespace are ignored) | `ErrDuplicateEmail`     |
| `member_id` must reference an existing member               | `ErrMemberNotFound`     |
| `point_type_id` must be a valid type (1–4)                  | `ErrInvalidPointType`   |
| `points` must be a positive number                          | `ErrNegativePoints`     |
| A redemption must not exceed the member's available balance | `ErrInsufficientPoints` |

---

## Project structure

```
.
├── cmd/
│   └── reward-point/
│       └── main.go          # Entry point
├── internal/
│   ├── models/
│   │   └── models.go        # Has models for Member, MemberInfo, Reward, PointTypeID
│   ├── database/
│   │   ├── store.go         # Temporary DB like interface
│   │   └── store_test.go    # Unit tests (coverage 95%+)
│   └── server/
│       ├── server.go        # HTTP server setup, start/stop
│       ├── handler.go       # Route registration and request handlers
│       └── dto.go           # Request payloads
├── Dockerfile
└── docker-compose.yml
```

---

## Testing

Run the unit test using:

```bash
go test -v ./...
```

With coverage:

```bash
go test -cover ./...
```
