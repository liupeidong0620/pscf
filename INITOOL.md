### pscf ini - Instructions for use

```shell
$ ./pscf ini -h
NAME:
   pscf ini - Provides tools for set, delete, and get operations for section or key in INI format.

USAGE:
   pscf ini [-c config.yaml] [--set section.key=value]
   pscf ini [-c config.yaml] [--del section.key]
   pscf ini [-c config.yaml] [--get section.key]
   pscf ini [-c config.yaml] [--getkeys section]

OPTIONS:
   --config FILE, -c FILE  Load configuration from FILE
   --del value             Delete a node in the configuration file.(--del node)
   --get value             Get the value of the node in the configuration file.(--get node)
   --set value             Set the value of the node and add a new node if the node does not exist.(--set node=value)
   --verbose, -v           Debug mode.(--verbose, -v)
   --getkeys value         Get the keys value of the node in the configuration file.(--getkeys node)
```

## Examples

```
$ cat test.ini
; Package name
NAME        = ini
; Package version
VERSION     = v1
; Package import path
IMPORT_PATH = gopkg.in/%(NAME)s.%(VERSION)s

# Information about package author
# Bio can be written in multiple lines.
[author]
NAME   = Unknwon
E-MAIL = u@gogs.io
GITHUB = https://github.com/%(NAME)s

[features]
-: Support read/write comments of keys and sections
-: Support auto-increment of key names
-: Support load multiple files to overwrite key values

[advance]
value with quotes      = "some value"
value quote2 again     = 'some value'
```

### --get
```
$ cat test.ini | ./pscf ini -c - --get [].NAME
ini

$ cat test.ini | ./pscf ini -c - --get [author].NAME
Unknwon

$ cat test.ini | ./pscf ini -c - --get '[advance].value with quotes'
some value
```

```
$ cat test.ini | ./pscf ini -c - --get [features].[0]
Support read/write comments of keys and sections

```

### --set

Update a value:
```
$ cat test.ini | ./pscf ini -c - --set [].NAME=ini_2
; Package name
NAME        = ini_2
; Package version
VERSION     = v1
; Package import path
IMPORT_PATH = gopkg.in/%(NAME)s.%(VERSION)s
......
```

Set a new value:
```
$ cat test.ini | ./pscf ini -c - --set [google].domain=www.google.cn
.....
[google]
domain = www.google.cn

```

Set a new section:
```
$ cat test.ini | ./pscf ini -c - --set [google]=
......
[google]
```

Set array:
```
$ cat test.ini | ./pscf ini -c - --set [features].[0]=test
....
[features]
-  = test
-  = Support auto-increment of key names
-  = Support load multiple files to overwrite key values
....


$ cat test.ini | ./pscf ini -c - --set [features].[+]=test
....
[features]
-  = Support read/write comments of keys and sections
-  = Support auto-increment of key names
-  = Support load multiple files to overwrite key values
-  = test
....
```

### --del

delete a key:
```
$ cat test.ini | ./pscf ini -c - --del [].NAME
; Package version
VERSION     = v1
; Package import path
IMPORT_PATH = gopkg.in/%(NAME)s.%(VERSION)s
...
```

delete a section:
```
$ cat test.ini | ./pscf ini -c - --del [author]
; Package name
NAME        = ini
; Package version
VERSION     = v1
; Package import path
IMPORT_PATH = gopkg.in/%(NAME)s.%(VERSION)s

[features]
-  = Support read/write comments of keys and sections
-  = Support auto-increment of key names
-  = Support load multiple files to overwrite key values

[advance]
value with quotes  = some value
value quote2 again = some value
```

delete array key:
```
$ cat test.ini | ./pscf ini -c - --del [features].[1]
...
[features]
-  = Support read/write comments of keys and sections
-  = Support load multiple files to overwrite key values
...

$ cat test.ini | ./pscf ini -c - --del [features].[-]
...
[features]
-  = Support read/write comments of keys and sections
-  = Support auto-increment of key names
...
```

### --getkeys

```
$ cat test.ini | ./pscf ini -c - --getkeys [author]
NAME
E-MAIL
GITHUB
```
