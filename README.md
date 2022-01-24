# cph-wallet
Repository by Glenn Harvey T. Liwanag

## Design Considerations

## Structure
The source is separated into four major layers accompanied by utility packages, and heavily operates using dependency injection with interfaces.
```
internal/
```

Migration files and other project-level scripts are found in `scripts/`.

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

Prepare the db container.
```
$ make startTestDb; make migrateDb
```

Run unit and integration tests.
```
$ make testWithIntegrations
```

Start the app container. `Ctrl+C` for graceful exit.
```
$ make startDev
```

On another shell session, perform API calls.
```
$ curl localhost:8080/transaction/v1/accounts

$ curl localhost:8080/metrics
```

Stop all containers
```
$ make stopContainers
```

See `Makefile` for details on commands.

## License
Explicit exclusive copyright as default for creative work. See [choosealicense section](https://choosealicense.com/no-permission/) on "No License."
