# Pscan - a simple port scanning program
Pscan is a simple port scanning program written in Go.
## Usage
```
pscan [options] [target1 target2 ...]
```
## Options
* `-p | --port` - the ports to be scanned. For example, `-p 80`, `-p 22,80,443`, `-p 0-99,443,900-1000`. 
Default value - `0-1023`.
* `-t | --timeout` - the timeout for a connection in seconds. If the timeout is expired, the port is considered to be filtered.
Default value - 20.
* `max-requests` - the maximum amount of simultaneous requests to a single host.
Default value - 10.
* `max-routines` - the maximum amount of goroutines.
Default value - 10000.
* `protocol` - the protocol used for connection. Available protocols are tcp, tcp4 (i.e. tcp only with ipv4), tcpv6 (i.e. tcp only with ipv6).
Default value - tcp.