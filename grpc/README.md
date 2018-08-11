# Mango-Store- Service

This is the Mango-Store- service

Generated with

```
micro new mango-store- --namespace=go.micro --type=srv
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: go.micro.srv.mango-store-
- Type: srv
- Alias: mango-store-

## Dependencies

Micro services depend on service discovery. The default is consul.

```
# install consul
brew install consul

# run consul
consul agent -dev
```

## Usage

A Makefile is included for convenience

Build the binary

```
make build
```

Run the service
```
./mango-store--srv
```

Build a docker image
```
make docker
```