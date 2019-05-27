# Project Spiderweb - Service

## License; use, modification, sharing, and distribution

* This project does **not** have an Open Source license and its copyright is only extended to the specified authors.
* You are not permitted to share the software.
* You are not permitted to distribute the software.
* You are not permitted to modify the software.
* You are not permitted to use the software, except at its hosted URL.

* You are, however, permitted to view and fork this repo.

You can read more about our permissions at https://choosealicense.com/no-permission/

## Development

### Contributing

If you want to get started on contributing, head over the [Root Project Page](https://github.com/Pergamene/project-spiderweb) and either check out the [Issues](https://github.com/Pergamene/project-spiderweb/issues) or [Projects](https://github.com/Pergamene/project-spiderweb/projects).  Not sure where to start?  You can [post your interest here](https://github.com/Pergamene/project-spiderweb/issues/2) and we'll get you started.

We keep a separate repo for Issues/Projects because the project spans more than one repo (front-end, back-end, etc).  If there is an issue specific to only this project, you can just [post an issue here](https://github.com/Pergamene/project-spiderweb-service/issues).

#### Adding dependencies

Dependencies are installed using [govendor](https://github.com/kardianos/govendor) until go2 is released.

##### Installation
```
go get -u github.com/kardianos/govendor
```

##### Use
```
govendor fetch <same as "go get" path>
```

### Setup

See the [Go Getting Started page](https://golang.org/doc/install) for details on how to set up your machine for Go development.

```
go get github.com/Pergamene/project-spiderweb-service
```

### Testing

First, install testify, which is used for unit test mocks and assertions: `go get github.com/stretchr/testify`.

To run the unit tests, you can run `go test ./...`.

If you're writing new unit tests, you'll likely need to install [mockery](https://github.com/vektra/mockery)
to mock interfaces.

If you want to run the intregated tests against the db, make sure to:

1. Set the env var `SETUP_SQL_FILEPATH` to
the filepath of where the db's [setup.sql](https://github.com/Pergamene/project-spiderweb-db/blob/master/setup.sql) file is on your local machine:

```
export SETUP_SQL_FILEPATH=/Users/rhyeen/Documents/repos/project-spiderweb/project-spiderweb-db/setup.sql
```

2. Have the database docker container running locally.  Follow the [README.md](https://github.com/Pergamene/project-spiderweb-db/blob/master/README.md) for instructions.

### Build/Run

To build the app, `cd cmd/server` and run `go build`. This will create a `server` executable that you can run

```
./server
```

You'll then be able to hit the service at `http://localhost:8782` try hitting `http://localhost:8782/healthcheck` to see the basic service is working or `http://localhost:8782/dbhealthcheck` to see if it can successfully connect to the database.

#### Serving API Docs locally

The API docs can be accessed when the server is running locally at: `http://localhost:8782/api/docs`.