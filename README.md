# trading-post
Store stock trading transactions with funny money

A hosted version of this can be found at [https://trading-post.club](https://trading-post.club).

## API

### GET `/profile/`

Return the current user profile, including all currently owned stocks.

```shell
curl --header 'Authorization: Bearer <id_token>' http://localhost:3000/profile/
# {
#   "name": "Roy van de Water",
#   "riches": -1067.8401,
#   "stocks": [
#     {
#       "quantity": 2,
#       "ticker": "ctxs"
#     },
#     {
#       "quantity": 1,
#       "ticker": "goog"
#     }
#   ]
# }
```

### POST `/profile/buy-orders`

Create a new buy order and subtract the market rate from the profile's riches.

```shell
curl \
  --request 'POST' \
  --header 'Authorization: Bearer <id_token>' \
  --data '{"ticker": "goog", "quantity": 3}' \
  http://localhost:3000/profile/buy-orders
# {
#   "id": "2fa48bba-dddb-4f23-8c95-bd2dd9e07ed3",
#   "price": 905.96,
#   "quantity": 3,
#   "ticker": "goog"
# }
```

### POST `/profile/sell-orders`

Create a new sell order and add the market rate to the profile's riches

```shell
curl \
  --request 'POST' \
  --header 'Authorization: Bearer <id_token>' \
  --data '{"ticker": "goog", "quantity": 3}' \
  http://localhost:3000/profile/sell-orders
# {
#   "id": "ffe2ea60-9695-4e2b-8fd3-c978020f213b",
#   "price": 905.96,
#   "quantity": 3,
#   "ticker": "goog"
# }
```
