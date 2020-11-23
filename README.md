# zNet

![Build Status](https://github.com/xaque208/znet/workflows/Compiling/badge.svg)

[![PkgGoDev](https://pkg.go.dev/badge/github.com/xaque208/znet)](https://pkg.go.dev/github.com/xaque208/znet)

> It isn't the destination, but the journey.

    agent       Run a znet agent
    arpwatch    Export Junos ARP data Pometheus
    builder     Run a git repo builder
    completion  Generates zsh completion scripts
    gitwatch    Run a git watcher
    harvest     Run an mqtt harvester
    help        Help about any command
    inv         Report on inventory
    lights      Collect status of the HUE lights for reporting
    netconfig   Configure Junos Devices
    server      Listen for commands/events/messages sent to the RPC server
    timer       Run a timer
    version     Show znet version

## Daemons

This project consists of a few daemons.

### Server

`znet server`

The zNet server is the gRPC server that has a few tasks.

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

#### Events

There is also a somewhat generalized eventing framework to allow various
components of the system to send events into the gRPC server. The server can
then forward those events to subscribers, either local or remote. This allows
events from one component to propagate to other parts of the system to perform
actions, store data, update metrics or something else.

### Agent

The agent receives events as part of a subscription call to the events server.
These events are used to trigger lighting changes, or execute commands based on
configuration.

### Builder

The builder receives events as part of a subscription call to the events
server.  This allows the builder to keep tabs on git repos, and perform certain
actions based on a configuration stored with the repository.
