# pscf - Processing standardized configuration file.

## Introduction
PSCF tools provide configuration file add, delete, and get operations.

### Supported configuration file formats

* YAML
* JSON
* INI

## Installation
`go get  github.com/liupeidong0620/pscf`

## Compile
```
linux
make build-linux

MacOS
make build-darwin

winOS
make build-windows
```

### Command line arguments
```
$ pscf -h
NAME:
   pscf - Processing standardized config file.

USAGE:
   pscf [yaml|json|ini] -h

VERSION:
   1.0.0 darwin (amd64) Build Date 2019-04-05 18:54:11.112875647 +0800 CST

AUTHOR:
   liupeidong <liupeidong0620@163.com>

COMMANDS:
     yaml     Provides tools for set, delete, and get operations for a node in YAML format.
     json     Provides tools for set, delete, and get operations for a node in JSON format.
     ini      Provides tools for set, delete, and get operations for section or key in INI format.
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version

COPYRIGHT:
   (c) 2018 LPD
```
### Example

* pscf yaml - YAMLTOOL.md
* pscf json - JSONTOOL.md
* pscf ini  - INITOOL.md

