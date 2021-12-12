# zNet

[![PkgGoDev](https://pkg.go.dev/badge/github.com/xaque208/znet)](https://pkg.go.dev/github.com/xaque208/znet)
![Build Status](https://github.com/xaque208/znet/workflows/Compiling/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/xaque208/znet/badge.svg)](https://coveralls.io/github/xaque208/znet)

> It isn't the destination, but the journey.

    Usage of ./bin/linux/znet:
      -config.file value
          Configuration file to load
      -otel_endpoint string
          otel endpoint, eg: tempo:4317
      -target string
          target module (default "all")

## Modules

This project consists of a few modules that are all loaded by default, some of
which are dependencies for a few key modules.

#### Inventory

The server is responsible for providing an interface into an inventory service.
Currently, this means LDAP, but this could potentially be a number of back-end
databases. I only have eyes for LDAP at the moment.

To complement this inventory storage is a modicum of code generation that is
based off the `rpc.proto` file. The `make proto` command will process the
template files and write out a few files.

- `internal/ldap.ldif` is the LDAP schema that must manually be written to the
  LDAP server.
- `internal/inventory/rpc_types.go` is the Go data types for all of the schema
  interface. We avoid using the RPC types directly, and create ones that map
  directly to what we will store in LDAP.
- `internal/inventory/network.go` contains all the CRUD methods for operating
  each of the various network types.

This has been quite handy, as it has allowed the schema to expand over time,
while also keeping the client implementation of that schema up-to-date.

#### Telemetry

The server has a few methods to report an observation of an inventory device.
This might come from a variety of sources, each with potentially partial
information. When an observation is made, the data is sent to the server for
handling. This might mean updating the data in the inventory database, or it
might mean some other alert.

#### Harvester

The harvester is used to read Zigbee messages from an MQTT bus that are
produced using [Zigbee2mqtt][z2m]. These messages are used to perform various
actions against the server for actions such as changing the state of other
Zigbee devices, such as lights, or other such event driven motivations.

[z2m]: https://www.zigbee2mqtt.io/
