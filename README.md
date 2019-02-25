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
go get https://github.com/Pergamene/project-spiderweb-service
```

### Testing

To run the unit tests, you can run `go test ./...`.

### Build

To build the app, `cd cmd/server` and run `go build`. This will create a `server` executable that you can run

```
./server
```

You'll then be able to hit the service at `http://localhost:8782` try hitting `http://localhost:8782/healthcheck` to see the basic service is working or `http://localhost:8782/dbhealthcheck` to see if it can successfully connect to the database.