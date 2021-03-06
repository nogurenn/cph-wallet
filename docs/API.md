# Show Accounts

**URL** : `/transaction/v1/accounts`

**Method** : `GET`

## Success Response

**Code** : `200 OK`

**Content** : Sorted by `id (username)` ascending.

```json
{
  "accounts": [
    {
      "id": "alice456",
      "balance": "270.62",
      "currency": "USD"
    },
    {
      "id": "bob123",
      "balance": "44.47",
      "currency": "USD"
    },
    {
      "id": "karen789",
      "balance": "284.91",
      "currency": "USD"
    }
  ],
  "error": null
}
```

# Show Payment Transactions

**URL** : `/transaction/v1/payments`

**Method** : `GET`

## Success Response

**Code** : `200 OK`

**Content** : Sorted by creation date of `Transaction` descending (latest first).

```json
{
  "payments": [
    {
      "id": "bbd569d4-9154-4e5e-ab34-e1e75e27c1c8",
      "name": "payment",
      "entries": [
        {
          "account": "karen789",
          "amount": "44.79",
          "to_account": "alice456",
          "direction": "outgoing"
        },
        {
          "account": "alice456",
          "amount": "44.79",
          "from_account": "karen789",
          "direction": "incoming"
        }
      ],
      "created_at": "2022-02-01T20:33:14.520032Z",
      "updated_at": "2022-02-01T20:33:14.520032Z"
    }
  ],
  "error": null
}
```

# Send Payment to Target Account

**URL** : `/transaction/v1/payments`

**Method** : `POST`

**Content**:
```json
{
  "username": "karen789",
  "target_username": "alice456",
  "amount": "44.79"
}
```

## Success Response

**Code** : `201 CREATED`

**Content** :

```json
{
  "error": null
}
```