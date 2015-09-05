go-pcd
======
[![Build Status](https://travis-ci.org/brimstone/go-pcd.svg)](https://travis-ci.org/brimstone/go-pcd)[![Coverage Status](https://coveralls.io/repos/brimstone/go-pcd/badge.svg?branch=master&service=github)](https://coveralls.io/github/brimstone/go-pcd?branch=master)

API daemon and controlling program for [pcd](https://github.com/brimstone/pcd).

Usage
-----

```bash
Pancake Crop Deli Control Program

Usage: 
  ./pcd [command]

Available Commands: 
  daemon      Run the API daemon
  docker/bip  Get or Set Docker Bridge IP
  hostname    Get or Set Hostname
  help        Help about any command

Flags:
  -a, --address="127.0.0.1:8080": Address for API server
  -h, --help[=false]: help for ./pcd


Use "./pcd [command] --help" for more information about a command.
```

Daemon
------

```bash
$ pcd daemon
```
