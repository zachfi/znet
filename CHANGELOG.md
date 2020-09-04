<a name="v0.25.2"></a>
## [v0.25.2] - 2020-09-02
<a name="v0.25.1"></a>
## [v0.25.1] - 2020-09-02
<a name="v0.25.0"></a>
## [v0.25.0] - 2020-09-01
### Bug Fixes
- **lint:** adjust
- **lint:** disperse complexity
- **rpc:** avoid nil pointer reference

### Features
- **iot:** track bridge state
- **iot:** start handling bridge log messages

<a name="v0.24.7"></a>
## [v0.24.7] - 2020-08-21
### Bug Fixes
- **agent:** include execution duration

<a name="v0.24.6"></a>
## [v0.24.6] - 2020-08-21
<a name="v0.24.5"></a>
## [v0.24.5] - 2020-08-21
<a name="v0.24.4"></a>
## [v0.24.4] - 2020-08-21
<a name="v0.24.3"></a>
## [v0.24.3] - 2020-08-21
<a name="v0.23.4"></a>
## [v0.23.4] - 2020-08-17
### Bug Fixes
- **build:** use git-chglog for release process

<a name="v0.24.2"></a>
## [v0.24.2] - 2020-08-17
### Bug Fixes
- **lint:** move ldap config validation up the stack

<a name="v0.24.1"></a>
## [v0.24.1] - 2020-08-17
### Bug Fixes
- **iot:** use pointers for floats to ensure stats handling

<a name="v0.24.0"></a>
## [v0.24.0] - 2020-08-16
### Bug Fixes
- **inv:** ensure empty basedn is detected
- **server:** ensure RPC is closed properly and avoid remote deadlock

<a name="v0.23.3"></a>
## [v0.23.3] - 2020-08-16
### Bug Fixes
- **ci:** force reset of checkout

<a name="v0.23.2"></a>
## [v0.23.2] - 2020-08-16
### Bug Fixes
- **iot:** include iot events

<a name="v0.23.1"></a>
## [v0.23.1] - 2020-08-16
### Bug Fixes
- **lint:** reduce complexity

<a name="v0.23.0"></a>
## [v0.23.0] - 2020-08-16
### Features
- **gitwatch:** restructure for collection interval

<a name="v0.22.0"></a>
## [v0.22.0] - 2020-08-14
### Bug Fixes
- **metrics:** improve metric names

### Features
- **lights:** initial zigbee integration
- **server:** include the eventmachine  on telemetryServer

<a name="v0.21.4"></a>
## [v0.21.4] - 2020-08-14
<a name="v0.21.3"></a>
## [v0.21.3] - 2020-08-14
<a name="v0.21.2"></a>
## [v0.21.2] - 2020-08-14
<a name="v0.21.1"></a>
## [v0.21.1] - 2020-08-14
<a name="v0.21.0"></a>
## [v0.21.0] - 2020-08-10
### Bug Fixes
- **ldap:** avoid renumbering oids for schema generation

### Features
- **zigbee:** include ZigbeeDevice struct

<a name="v0.20.6"></a>
## [v0.20.6] - 2020-07-26
<a name="v0.20.5"></a>
## [v0.20.5] - 2020-07-26
<a name="v0.20.4"></a>
## [v0.20.4] - 2020-07-26
<a name="v0.20.3"></a>
## [v0.20.3] - 2020-07-26
### Bug Fixes
- **timer:** ensure repeat events are handled promptly

<a name="v0.20.2"></a>
## [v0.20.2] - 2020-07-06
<a name="v0.20.1"></a>
## [v0.20.1] - 2020-07-06
### Bug Fixes
- **vault:** better fall-through for login handling

<a name="v0.20.0"></a>
## [v0.20.0] - 2020-07-06
### Documentation Updates
- include changelog

### Features
- **vault:** implement cert login

<a name="v0.19.0"></a>
## [v0.19.0] - 2020-06-28
### Features
- **arpwatch:** include ARP reporter command
- **inventory:** overhaul inventory usage
- **netconfig:** include initial netconfig work
- **rpc:** overhaul the rpc servers
- **rpc:** begin code generation from proto file

<a name="v0.18.3"></a>
## [v0.18.3] - 2020-05-31
<a name="v0.18.2"></a>
## [v0.18.2] - 2020-05-26
<a name="v0.18.1"></a>
## [v0.18.1] - 2020-05-17
<a name="v0.18.0"></a>
## [v0.18.0] - 2020-05-17
### Features
- **znet:** support TLS for gRPC connectivity

<a name="v0.17.3"></a>
## [v0.17.3] - 2020-05-15
### Bug Fixes
- **ci:** reset the working tree, without the hash
- **ci:** reset the working tree
- **lint:** adjust

### Features
- **agent:** include filtering on git collection
- **ci:** include semver parsing for latest tag retrieval
- **inventory:** include utility methods for ldap calls
- **inventory:** provide more structure for recording inventory
- **iot:** parse led config and color messages

<a name="v0.17.2"></a>
## [v0.17.2] - 2020-05-09
<a name="v0.17.1"></a>
## [v0.17.1] - 2020-05-09
### Bug Fixes
- **lint:** adjust
- **lint:** adjust

<a name="v0.17.0"></a>
## [v0.17.0] - 2020-05-08
### Bug Fixes
- **gitlab:** drop support for gitlab
- **lint:** adjust

### Features
- **iot:** include iot package for homeassistant spec
- **rpc:** implement new things service
- **thingharvest:** include MQTT harvester to publish over RPC

<a name="v0.16.7"></a>
## [v0.16.7] - 2020-05-03
### Bug Fixes
- **znet:** flag handling for server listening

<a name="v0.16.6"></a>
## [v0.16.6] - 2020-05-03
### Documentation Updates
- improve

<a name="v0.16.5"></a>
## [v0.16.5] - 2020-05-03
<a name="v0.16.4"></a>
## [v0.16.4] - 2020-05-03
<a name="v0.16.3"></a>
## [v0.16.3] - 2020-05-02
<a name="v0.16.2"></a>
## [v0.16.2] - 2020-05-02
<a name="v0.16.1"></a>
## [v0.16.1] - 2020-04-30
<a name="v0.16.0"></a>
## [v0.16.0] - 2020-04-30
<a name="v0.15.0"></a>
## [v0.15.0] - 2020-04-26
### Bug Fixes
- **ci:** avoid calls on nil repo

<a name="v0.14.9"></a>
## [v0.14.9] - 2020-04-26
<a name="v0.14.8"></a>
## [v0.14.8] - 2020-04-26
### Features
- **agent:** include duration on ExecutionResult

<a name="v0.14.7"></a>
## [v0.14.7] - 2020-04-26
### Features
- **znet:** start to capture execution result metrics

<a name="v0.14.6"></a>
## [v0.14.6] - 2020-04-26
### Bug Fixes
- **lint:** adjust

<a name="v0.14.5"></a>
## [v0.14.5] - 2020-04-26
<a name="v0.14.4"></a>
## [v0.14.4] - 2020-04-26
<a name="v0.14.3"></a>
## [v0.14.3] - 2020-04-26
<a name="v0.14.2"></a>
## [v0.14.2] - 2020-04-26
### Bug Fixes
- **cmd:** nil check on builder command
- **cmd:** nil check on agent command

<a name="v0.14.1"></a>
## [v0.14.1] - 2020-04-26
### Bug Fixes
- **cmd:** deprecations

<a name="v0.14.0"></a>
## [v0.14.0] - 2020-04-26
<a name="v0.13.14"></a>
## [v0.13.14] - 2020-04-26
### Bug Fixes
- **znet:** avoid closing connections were never created

<a name="v0.13.13"></a>
## [v0.13.13] - 2020-04-25
<a name="v0.13.12"></a>
## [v0.13.12] - 2020-04-25
### Bug Fixes
- **release:** update release process

<a name="v0.13.11"></a>
## [v0.13.11] - 2020-04-25
### Bug Fixes
- **lint:** adjust

<a name="v0.13.10"></a>
## [v0.13.10] - 2020-04-25
<a name="v0.13.9"></a>
## [v0.13.9] - 2020-04-25
<a name="v0.13.8"></a>
## [v0.13.8] - 2020-04-25
<a name="v0.13.7"></a>
## [v0.13.7] - 2020-04-25
<a name="v0.13.6"></a>
## [v0.13.6] - 2020-04-25
<a name="v0.13.5"></a>
## [v0.13.5] - 2020-04-25
### Features
- **ci:** new package for git repo tracking

<a name="v0.13.4"></a>
## [v0.13.4] - 2020-04-25
### Bug Fixes
- **lint:** adjust

<a name="v0.13.3"></a>
## [v0.13.3] - 2020-04-25
### Features
- **builder:** handle git cache and checkout

<a name="v0.13.0"></a>
## [v0.13.0] - 2020-04-25
<a name="v0.13.2"></a>
## [v0.13.2] - 2020-04-25
<a name="v0.13.1"></a>
## [v0.13.1] - 2020-04-25
### Documentation Updates
- improve

### Features
- **builder:** start a gitwatch builder

<a name="v0.12.18"></a>
## [v0.12.18] - 2020-04-24
<a name="v0.12.14"></a>
## [v0.12.14] - 2020-04-24
<a name="v0.12.21"></a>
## [v0.12.21] - 2020-04-24
<a name="v0.12.17"></a>
## [v0.12.17] - 2020-04-24
<a name="v0.12.16"></a>
## [v0.12.16] - 2020-04-24
<a name="v0.12.15"></a>
## [v0.12.15] - 2020-04-24
<a name="v0.12.19"></a>
## [v0.12.19] - 2020-04-24
<a name="v0.12.20"></a>
## [v0.12.20] - 2020-04-24
<a name="v0.12.22"></a>
## [v0.12.22] - 2020-04-24
<a name="v0.12.9"></a>
## [v0.12.9] - 2020-04-24
<a name="v0.12.8"></a>
## [v0.12.8] - 2020-04-24
<a name="v0.12.12"></a>
## [v0.12.12] - 2020-04-24
<a name="v0.12.11"></a>
## [v0.12.11] - 2020-04-24
<a name="v0.12.13"></a>
## [v0.12.13] - 2020-04-24
<a name="v0.12.10"></a>
## [v0.12.10] - 2020-04-24
<a name="v0.12.7"></a>
## [v0.12.7] - 2020-04-24
<a name="v0.12.6"></a>
## [v0.12.6] - 2020-04-24
<a name="v0.12.5"></a>
## [v0.12.5] - 2020-04-24
### Bug Fixes
- **agent:** better error checking

<a name="v0.12.4"></a>
## [v0.12.4] - 2020-04-24
<a name="v0.12.3"></a>
## [v0.12.3] - 2020-04-24
### Features
- **agent:** begin command execution result response

<a name="v0.12.2"></a>
## [v0.12.2] - 2020-04-24
<a name="v0.12.1"></a>
## [v0.12.1] - 2020-04-24
<a name="v0.12.0"></a>
## [v0.12.0] - 2020-04-22
### Bug Fixes
- **znet:** fix signal handling in a few places

<a name="v0.11.13"></a>
## [v0.11.13] - 2020-04-21
<a name="v0.11.14"></a>
## [v0.11.14] - 2020-04-21
### Documentation Updates
- improve command help

<a name="v0.11.12"></a>
## [v0.11.12] - 2020-04-20
<a name="v0.11.11"></a>
## [v0.11.11] - 2020-04-20
<a name="v0.11.9"></a>
## [v0.11.9] - 2020-04-19
### Bug Fixes
- **gitwatch:** better public key handling

<a name="v0.11.8"></a>
## [v0.11.8] - 2020-04-19
### Bug Fixes
- **gitwatch:** better public key handling

<a name="v0.11.7"></a>
## [v0.11.7] - 2020-04-19
<a name="v0.11.10"></a>
## [v0.11.10] - 2020-04-19
### Bug Fixes
- **gitwatch:** better public key handling
- **gitwatch:** better public key handling
- **gitwatch:** better public key handling for existing clones

<a name="v0.11.6"></a>
## [v0.11.6] - 2020-04-18
### Bug Fixes
- **timer:** set a default future_limit to avoid fast loop

<a name="v0.11.5"></a>
## [v0.11.5] - 2020-04-18
### Bug Fixes
- **scheduler:** handle rescheduling

<a name="v0.11.4"></a>
## [v0.11.4] - 2020-04-17
<a name="v0.11.3"></a>
## [v0.11.3] - 2020-04-17
### Bug Fixes
- **scheduler:** improve retry and recent event detection

<a name="v0.11.2"></a>
## [v0.11.2] - 2020-04-17
<a name="v0.11.1"></a>
## [v0.11.1] - 2020-04-16
<a name="v0.11.0"></a>
## [v0.11.0] - 2020-04-14
<a name="v0.10.0"></a>
## [v0.10.0] - 2020-04-13
### Bug Fixes
- **agent:** update
- **agent:** try to clean up connections
- **znet:** improve repo auth handling for http/git
- **znet:** update connection handling and timeouts

<a name="v0.9.9"></a>
## [v0.9.9] - 2020-04-12
### Bug Fixes
- **lights:** change subscription names to include SolarEvent
- **timer:** ensure that zero value names are not present

<a name="v0.9.8"></a>
## [v0.9.8] - 2020-04-12
### Bug Fixes
- **timer:** make sure events within the last five seconds are scheduled

### Features
- **agent:** begin an agent to receive events

<a name="v0.9.7"></a>
## [v0.9.7] - 2020-04-11
### Bug Fixes
- **lint:** adjust
- **listen:** ensure remoteChan is created, drop eventNames init as duplication

### Features
- **agent:** begin an agent to receive events
- **eventserver:** implement rpc SubscribeEvents method
- **rpc:** remote event subscription

<a name="v0.9.6"></a>
## [v0.9.6] - 2020-04-10
### Features
- **ssh:** add support for key file auth

<a name="v0.9.5"></a>
## [v0.9.5] - 2020-04-10
<a name="v0.9.4"></a>
## [v0.9.4] - 2020-04-10
<a name="v0.9.3"></a>
## [v0.9.3] - 2020-04-09
<a name="v0.9.1"></a>
## [v0.9.1] - 2020-04-09
<a name="v0.9.2"></a>
## [v0.9.2] - 2020-04-09
<a name="v0.9.0"></a>
## [v0.9.0] - 2020-04-09
### Features
- **gitwatch:** include gitwatch command

<a name="v0.8.3"></a>
## [v0.8.3] - 2020-04-05
### Bug Fixes
- **astro:** handle empty name set

<a name="v0.8.2"></a>
## [v0.8.2] - 2020-04-03
### Bug Fixes
- **scheduler:** nil pointer check

<a name="v0.8.1"></a>
## [v0.8.1] - 2020-03-31
### Bug Fixes
- **producer:** ensure that multiple producers run with each other

<a name="v0.8.0"></a>
## [v0.8.0] - 2020-03-30
### Bug Fixes
- **build:** closer to passing new build requirements

### Features
- **completion:** include zsh shell completion
- **events:** update the producer contract to include start and stop methods
- **lights:** implement rpc light alerts
- **znet:** implement light alerts

<a name="v0.7.11"></a>
## [v0.7.11] - 2020-03-22
<a name="v0.7.10"></a>
## [v0.7.10] - 2020-03-19
### Bug Fixes
- **log:** raise log level

<a name="v0.7.9"></a>
## [v0.7.9] - 2020-03-19
### Bug Fixes
- **timer:** raise log level for timer start messages

<a name="v0.7.8"></a>
## [v0.7.8] - 2020-03-19
### Bug Fixes
- **log:** pass lint
- **znet:** allow missing LDAP configuration

<a name="v0.7.6"></a>
## [v0.7.6] - 2020-03-15
<a name="v0.7.7"></a>
## [v0.7.7] - 2020-03-15
<a name="v0.7.5"></a>
## [v0.7.5] - 2020-03-15
<a name="v0.7.4"></a>
## [v0.7.4] - 2020-03-15
<a name="v0.7.3"></a>
## [v0.7.3] - 2020-03-15
### Features
- **astro:** include preSunset event

<a name="v0.7.2"></a>
## [v0.7.2] - 2020-03-15
### Bug Fixes
- **lights:** modify yaml config for puppet

<a name="v0.7.1"></a>
## [v0.7.1] - 2020-03-14
<a name="v0.7.0"></a>
## [v0.7.0] - 2020-03-13
### Features
- **release:** start using goreleaser

<a name="v0.6.5"></a>
## [v0.6.5] - 2020-01-31
<a name="v0.6.4"></a>
## [v0.6.4] - 2020-01-31
<a name="v0.6.3"></a>
## [v0.6.3] - 2020-01-31
<a name="v0.6.2"></a>
## [v0.6.2] - 2020-01-26
<a name="v0.6.1"></a>
## [v0.6.1] - 2020-01-26
<a name="v0.6.0"></a>
## [v0.6.0] - 2020-01-25
<a name="v0.5.2"></a>
## [v0.5.2] - 2020-01-16
<a name="v0.5.1"></a>
## [v0.5.1] - 2020-01-16
<a name="v0.5.0"></a>
## [v0.5.0] - 2020-01-16
<a name="v0.4.0"></a>
## [v0.4.0] - 2019-12-02
<a name="v0.3.0"></a>
## [v0.3.0] - 2019-08-31
<a name="v0.2.0"></a>
## [v0.2.0] - 2019-08-20
<a name="v0.1.0"></a>
## [v0.1.0] - 2019-02-24
<a name="v0.0.9"></a>
## [v0.0.9] - 2018-09-18
<a name="v0.0.8"></a>
## [v0.0.8] - 2018-09-17
<a name="v0.0.7"></a>
## [v0.0.7] - 2018-09-13
<a name="v0.0.6"></a>
## [v0.0.6] - 2018-09-13
<a name="v0.0.5"></a>
## [v0.0.5] - 2018-09-08
<a name="v0.0.4"></a>
## [v0.0.4] - 2018-09-08
<a name="v0.0.3"></a>
## [v0.0.3] - 2018-09-03
<a name="v0.0.2"></a>
## [v0.0.2] - 2018-09-03
<a name="v0.0.1"></a>
## v0.0.1 - 2018-09-03
[Unreleased]: https://github.com/xaque208/znet/compare/v0.25.2...HEAD
[v0.25.2]: https://github.com/xaque208/znet/compare/v0.25.1...v0.25.2
[v0.25.1]: https://github.com/xaque208/znet/compare/v0.25.0...v0.25.1
[v0.25.0]: https://github.com/xaque208/znet/compare/v0.24.7...v0.25.0
[v0.24.7]: https://github.com/xaque208/znet/compare/v0.24.6...v0.24.7
[v0.24.6]: https://github.com/xaque208/znet/compare/v0.24.5...v0.24.6
[v0.24.5]: https://github.com/xaque208/znet/compare/v0.24.4...v0.24.5
[v0.24.4]: https://github.com/xaque208/znet/compare/v0.24.3...v0.24.4
[v0.24.3]: https://github.com/xaque208/znet/compare/v0.23.4...v0.24.3
[v0.23.4]: https://github.com/xaque208/znet/compare/v0.24.2...v0.23.4
[v0.24.2]: https://github.com/xaque208/znet/compare/v0.24.1...v0.24.2
[v0.24.1]: https://github.com/xaque208/znet/compare/v0.24.0...v0.24.1
[v0.24.0]: https://github.com/xaque208/znet/compare/v0.23.3...v0.24.0
[v0.23.3]: https://github.com/xaque208/znet/compare/v0.23.2...v0.23.3
[v0.23.2]: https://github.com/xaque208/znet/compare/v0.23.1...v0.23.2
[v0.23.1]: https://github.com/xaque208/znet/compare/v0.23.0...v0.23.1
[v0.23.0]: https://github.com/xaque208/znet/compare/v0.22.0...v0.23.0
[v0.22.0]: https://github.com/xaque208/znet/compare/v0.21.4...v0.22.0
[v0.21.4]: https://github.com/xaque208/znet/compare/v0.21.3...v0.21.4
[v0.21.3]: https://github.com/xaque208/znet/compare/v0.21.2...v0.21.3
[v0.21.2]: https://github.com/xaque208/znet/compare/v0.21.1...v0.21.2
[v0.21.1]: https://github.com/xaque208/znet/compare/v0.21.0...v0.21.1
[v0.21.0]: https://github.com/xaque208/znet/compare/v0.20.6...v0.21.0
[v0.20.6]: https://github.com/xaque208/znet/compare/v0.20.5...v0.20.6
[v0.20.5]: https://github.com/xaque208/znet/compare/v0.20.4...v0.20.5
[v0.20.4]: https://github.com/xaque208/znet/compare/v0.20.3...v0.20.4
[v0.20.3]: https://github.com/xaque208/znet/compare/v0.20.2...v0.20.3
[v0.20.2]: https://github.com/xaque208/znet/compare/v0.20.1...v0.20.2
[v0.20.1]: https://github.com/xaque208/znet/compare/v0.20.0...v0.20.1
[v0.20.0]: https://github.com/xaque208/znet/compare/v0.19.0...v0.20.0
[v0.19.0]: https://github.com/xaque208/znet/compare/v0.18.3...v0.19.0
[v0.18.3]: https://github.com/xaque208/znet/compare/v0.18.2...v0.18.3
[v0.18.2]: https://github.com/xaque208/znet/compare/v0.18.1...v0.18.2
[v0.18.1]: https://github.com/xaque208/znet/compare/v0.18.0...v0.18.1
[v0.18.0]: https://github.com/xaque208/znet/compare/v0.17.3...v0.18.0
[v0.17.3]: https://github.com/xaque208/znet/compare/v0.17.2...v0.17.3
[v0.17.2]: https://github.com/xaque208/znet/compare/v0.17.1...v0.17.2
[v0.17.1]: https://github.com/xaque208/znet/compare/v0.17.0...v0.17.1
[v0.17.0]: https://github.com/xaque208/znet/compare/v0.16.7...v0.17.0
[v0.16.7]: https://github.com/xaque208/znet/compare/v0.16.6...v0.16.7
[v0.16.6]: https://github.com/xaque208/znet/compare/v0.16.5...v0.16.6
[v0.16.5]: https://github.com/xaque208/znet/compare/v0.16.4...v0.16.5
[v0.16.4]: https://github.com/xaque208/znet/compare/v0.16.3...v0.16.4
[v0.16.3]: https://github.com/xaque208/znet/compare/v0.16.2...v0.16.3
[v0.16.2]: https://github.com/xaque208/znet/compare/v0.16.1...v0.16.2
[v0.16.1]: https://github.com/xaque208/znet/compare/v0.16.0...v0.16.1
[v0.16.0]: https://github.com/xaque208/znet/compare/v0.15.0...v0.16.0
[v0.15.0]: https://github.com/xaque208/znet/compare/v0.14.9...v0.15.0
[v0.14.9]: https://github.com/xaque208/znet/compare/v0.14.8...v0.14.9
[v0.14.8]: https://github.com/xaque208/znet/compare/v0.14.7...v0.14.8
[v0.14.7]: https://github.com/xaque208/znet/compare/v0.14.6...v0.14.7
[v0.14.6]: https://github.com/xaque208/znet/compare/v0.14.5...v0.14.6
[v0.14.5]: https://github.com/xaque208/znet/compare/v0.14.4...v0.14.5
[v0.14.4]: https://github.com/xaque208/znet/compare/v0.14.3...v0.14.4
[v0.14.3]: https://github.com/xaque208/znet/compare/v0.14.2...v0.14.3
[v0.14.2]: https://github.com/xaque208/znet/compare/v0.14.1...v0.14.2
[v0.14.1]: https://github.com/xaque208/znet/compare/v0.14.0...v0.14.1
[v0.14.0]: https://github.com/xaque208/znet/compare/v0.13.14...v0.14.0
[v0.13.14]: https://github.com/xaque208/znet/compare/v0.13.13...v0.13.14
[v0.13.13]: https://github.com/xaque208/znet/compare/v0.13.12...v0.13.13
[v0.13.12]: https://github.com/xaque208/znet/compare/v0.13.11...v0.13.12
[v0.13.11]: https://github.com/xaque208/znet/compare/v0.13.10...v0.13.11
[v0.13.10]: https://github.com/xaque208/znet/compare/v0.13.9...v0.13.10
[v0.13.9]: https://github.com/xaque208/znet/compare/v0.13.8...v0.13.9
[v0.13.8]: https://github.com/xaque208/znet/compare/v0.13.7...v0.13.8
[v0.13.7]: https://github.com/xaque208/znet/compare/v0.13.6...v0.13.7
[v0.13.6]: https://github.com/xaque208/znet/compare/v0.13.5...v0.13.6
[v0.13.5]: https://github.com/xaque208/znet/compare/v0.13.4...v0.13.5
[v0.13.4]: https://github.com/xaque208/znet/compare/v0.13.3...v0.13.4
[v0.13.3]: https://github.com/xaque208/znet/compare/v0.13.0...v0.13.3
[v0.13.0]: https://github.com/xaque208/znet/compare/v0.13.2...v0.13.0
[v0.13.2]: https://github.com/xaque208/znet/compare/v0.13.1...v0.13.2
[v0.13.1]: https://github.com/xaque208/znet/compare/v0.12.18...v0.13.1
[v0.12.18]: https://github.com/xaque208/znet/compare/v0.12.14...v0.12.18
[v0.12.14]: https://github.com/xaque208/znet/compare/v0.12.21...v0.12.14
[v0.12.21]: https://github.com/xaque208/znet/compare/v0.12.17...v0.12.21
[v0.12.17]: https://github.com/xaque208/znet/compare/v0.12.16...v0.12.17
[v0.12.16]: https://github.com/xaque208/znet/compare/v0.12.15...v0.12.16
[v0.12.15]: https://github.com/xaque208/znet/compare/v0.12.19...v0.12.15
[v0.12.19]: https://github.com/xaque208/znet/compare/v0.12.20...v0.12.19
[v0.12.20]: https://github.com/xaque208/znet/compare/v0.12.22...v0.12.20
[v0.12.22]: https://github.com/xaque208/znet/compare/v0.12.9...v0.12.22
[v0.12.9]: https://github.com/xaque208/znet/compare/v0.12.8...v0.12.9
[v0.12.8]: https://github.com/xaque208/znet/compare/v0.12.12...v0.12.8
[v0.12.12]: https://github.com/xaque208/znet/compare/v0.12.11...v0.12.12
[v0.12.11]: https://github.com/xaque208/znet/compare/v0.12.13...v0.12.11
[v0.12.13]: https://github.com/xaque208/znet/compare/v0.12.10...v0.12.13
[v0.12.10]: https://github.com/xaque208/znet/compare/v0.12.7...v0.12.10
[v0.12.7]: https://github.com/xaque208/znet/compare/v0.12.6...v0.12.7
[v0.12.6]: https://github.com/xaque208/znet/compare/v0.12.5...v0.12.6
[v0.12.5]: https://github.com/xaque208/znet/compare/v0.12.4...v0.12.5
[v0.12.4]: https://github.com/xaque208/znet/compare/v0.12.3...v0.12.4
[v0.12.3]: https://github.com/xaque208/znet/compare/v0.12.2...v0.12.3
[v0.12.2]: https://github.com/xaque208/znet/compare/v0.12.1...v0.12.2
[v0.12.1]: https://github.com/xaque208/znet/compare/v0.12.0...v0.12.1
[v0.12.0]: https://github.com/xaque208/znet/compare/v0.11.13...v0.12.0
[v0.11.13]: https://github.com/xaque208/znet/compare/v0.11.14...v0.11.13
[v0.11.14]: https://github.com/xaque208/znet/compare/v0.11.12...v0.11.14
[v0.11.12]: https://github.com/xaque208/znet/compare/v0.11.11...v0.11.12
[v0.11.11]: https://github.com/xaque208/znet/compare/v0.11.9...v0.11.11
[v0.11.9]: https://github.com/xaque208/znet/compare/v0.11.8...v0.11.9
[v0.11.8]: https://github.com/xaque208/znet/compare/v0.11.7...v0.11.8
[v0.11.7]: https://github.com/xaque208/znet/compare/v0.11.10...v0.11.7
[v0.11.10]: https://github.com/xaque208/znet/compare/v0.11.6...v0.11.10
[v0.11.6]: https://github.com/xaque208/znet/compare/v0.11.5...v0.11.6
[v0.11.5]: https://github.com/xaque208/znet/compare/v0.11.4...v0.11.5
[v0.11.4]: https://github.com/xaque208/znet/compare/v0.11.3...v0.11.4
[v0.11.3]: https://github.com/xaque208/znet/compare/v0.11.2...v0.11.3
[v0.11.2]: https://github.com/xaque208/znet/compare/v0.11.1...v0.11.2
[v0.11.1]: https://github.com/xaque208/znet/compare/v0.11.0...v0.11.1
[v0.11.0]: https://github.com/xaque208/znet/compare/v0.10.0...v0.11.0
[v0.10.0]: https://github.com/xaque208/znet/compare/v0.9.9...v0.10.0
[v0.9.9]: https://github.com/xaque208/znet/compare/v0.9.8...v0.9.9
[v0.9.8]: https://github.com/xaque208/znet/compare/v0.9.7...v0.9.8
[v0.9.7]: https://github.com/xaque208/znet/compare/v0.9.6...v0.9.7
[v0.9.6]: https://github.com/xaque208/znet/compare/v0.9.5...v0.9.6
[v0.9.5]: https://github.com/xaque208/znet/compare/v0.9.4...v0.9.5
[v0.9.4]: https://github.com/xaque208/znet/compare/v0.9.3...v0.9.4
[v0.9.3]: https://github.com/xaque208/znet/compare/v0.9.1...v0.9.3
[v0.9.1]: https://github.com/xaque208/znet/compare/v0.9.2...v0.9.1
[v0.9.2]: https://github.com/xaque208/znet/compare/v0.9.0...v0.9.2
[v0.9.0]: https://github.com/xaque208/znet/compare/v0.8.3...v0.9.0
[v0.8.3]: https://github.com/xaque208/znet/compare/v0.8.2...v0.8.3
[v0.8.2]: https://github.com/xaque208/znet/compare/v0.8.1...v0.8.2
[v0.8.1]: https://github.com/xaque208/znet/compare/v0.8.0...v0.8.1
[v0.8.0]: https://github.com/xaque208/znet/compare/v0.7.11...v0.8.0
[v0.7.11]: https://github.com/xaque208/znet/compare/v0.7.10...v0.7.11
[v0.7.10]: https://github.com/xaque208/znet/compare/v0.7.9...v0.7.10
[v0.7.9]: https://github.com/xaque208/znet/compare/v0.7.8...v0.7.9
[v0.7.8]: https://github.com/xaque208/znet/compare/v0.7.6...v0.7.8
[v0.7.6]: https://github.com/xaque208/znet/compare/v0.7.7...v0.7.6
[v0.7.7]: https://github.com/xaque208/znet/compare/v0.7.5...v0.7.7
[v0.7.5]: https://github.com/xaque208/znet/compare/v0.7.4...v0.7.5
[v0.7.4]: https://github.com/xaque208/znet/compare/v0.7.3...v0.7.4
[v0.7.3]: https://github.com/xaque208/znet/compare/v0.7.2...v0.7.3
[v0.7.2]: https://github.com/xaque208/znet/compare/v0.7.1...v0.7.2
[v0.7.1]: https://github.com/xaque208/znet/compare/v0.7.0...v0.7.1
[v0.7.0]: https://github.com/xaque208/znet/compare/v0.6.5...v0.7.0
[v0.6.5]: https://github.com/xaque208/znet/compare/v0.6.4...v0.6.5
[v0.6.4]: https://github.com/xaque208/znet/compare/v0.6.3...v0.6.4
[v0.6.3]: https://github.com/xaque208/znet/compare/v0.6.2...v0.6.3
[v0.6.2]: https://github.com/xaque208/znet/compare/v0.6.1...v0.6.2
[v0.6.1]: https://github.com/xaque208/znet/compare/v0.6.0...v0.6.1
[v0.6.0]: https://github.com/xaque208/znet/compare/v0.5.2...v0.6.0
[v0.5.2]: https://github.com/xaque208/znet/compare/v0.5.1...v0.5.2
[v0.5.1]: https://github.com/xaque208/znet/compare/v0.5.0...v0.5.1
[v0.5.0]: https://github.com/xaque208/znet/compare/v0.4.0...v0.5.0
[v0.4.0]: https://github.com/xaque208/znet/compare/v0.3.0...v0.4.0
[v0.3.0]: https://github.com/xaque208/znet/compare/v0.2.0...v0.3.0
[v0.2.0]: https://github.com/xaque208/znet/compare/v0.1.0...v0.2.0
[v0.1.0]: https://github.com/xaque208/znet/compare/v0.0.9...v0.1.0
[v0.0.9]: https://github.com/xaque208/znet/compare/v0.0.8...v0.0.9
[v0.0.8]: https://github.com/xaque208/znet/compare/v0.0.7...v0.0.8
[v0.0.7]: https://github.com/xaque208/znet/compare/v0.0.6...v0.0.7
[v0.0.6]: https://github.com/xaque208/znet/compare/v0.0.5...v0.0.6
[v0.0.5]: https://github.com/xaque208/znet/compare/v0.0.4...v0.0.5
[v0.0.4]: https://github.com/xaque208/znet/compare/v0.0.3...v0.0.4
[v0.0.3]: https://github.com/xaque208/znet/compare/v0.0.2...v0.0.3
[v0.0.2]: https://github.com/xaque208/znet/compare/v0.0.1...v0.0.2
