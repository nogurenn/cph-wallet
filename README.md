# cph-wallet
Repository by Glenn Harvey T. Liwanag

## Design Considerations

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
