# cph-wallet
Repository by Glenn Harvey T. Liwanag

## Design Considerations
* The project design is based on **double-entry bookkeeping**, with `Transaction` as parent, and `Entry` as children. This was chosen to prioritize **auditability**, a critical aspect of finance-related systems.
  * `Transaction` (`payment`)
    * `Entry` (`outgoing`)
    * `Entry` (`incoming`)
  * `Transaction` (`deposit`)
    * `Entry` (`incoming`)
* Represent financial numbers using `shopspring/decimal` to safely operate on them and avoid silent precision loss.
* Serialize decimal fields in responses to string instead of JSON number to prevent precision loss.
* Use `UUID` as primary keys in db tables.
* Calculations, such as `SUM(credit, debit)`, are deferred to the DB as much as possible, trading away simpler SQL queries and commands in exchange for easier performance gains right from the start.

## Structure
```
cph-wallet
├── dbutil/         - db-related utilities e.g. driver/conn
├── docs/           - various docfiles e.g. API.md
├── mocks/          - mock interfaces used in tests
├── scripts/        - various project scripts and sql migration files
├── transaction/    - domain src
└── util/           - utility functions
```

## Installation
Install [Docker Engine](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/). Make sure `docker-compose` is at least `v2`.

## Usage
Go to project root.
```
$ cd cph-wallet/
```

Prepare the local envfile.
```
$ make env
```

Pull required `go` modules.
```
$ make deps
```

Run unit and integration tests.
```
$ make testWithIntegrations
```

Start the app . `Ctrl+C` for graceful exit.
```
$ make startTestDB; make migrateDB; make startDev
```

On another shell session, perform API calls.
```
$ curl localhost:8080/transaction/v1/accounts

$ curl -X POST -H "Content-Type: application/json" \
--data '{"username":"karen789","target_username":"alice456","amount": "44.79"}' \
localhost:8080/transaction/v1/payments

$ curl localhost:8080/transaction/v1/payments

$ curl localhost:8080/metrics
```

Stop all containers.
```
$ make stopContainers
```

See `Makefile` for details on commands.

## License
Explicit exclusive copyright as default for creative work. See [choosealicense section](https://choosealicense.com/no-permission/) on "No License."
