## pscf yaml - Instructions for use

```
$ ./releases/pscf yaml -h
NAME:
   pscf yaml - Provides tools for set, delete, and get operations for a node in YAML format.

USAGE:
   pscf yaml [-c config.yaml] [--set node.node1.node2=value]
   pscf yaml [-c config.yaml] [--del node.node1.node2]
   pscf yaml [-c config.yaml] [--get node.node1.node2]
   pscf yaml [-c config.yaml] [--getkeys node.node1.node2]

OPTIONS:
   --config FILE, -c FILE  Load configuration from FILE
   --del value             Delete a node in the configuration file.(--del node)
   --get value             Get the value of the node in the configuration file.(--get node)
   --set value             Set the value of the node and add a new node if the node does not exist.(--set node=value)
   --verbose, -v           Debug mode.(--verbose, -v)
   --tostr value           Force some tags into strings.(--tostr "yes,no")
   --getkeys value         Get the keys value of the node in the configuration file.(--getkeys node)
   --comment               Load comment enable.(--comment)
```

## Examples

```
$ cat test_1.yaml
fruits:
- Apple
- Orange
- Strawberry
- Mango
martin:
  name: Martin D'vloper
  job: Developer
  skill: Elite

```

```
$ cat test_2.yaml
- martin:
    name: Martin D'vloper
    job: Developer
    skills:
    - python
    - perl
    - pascal
- tabitha:
    name: Tabitha Bitumen
    job: Developer
    skills:
    - lisp
    - fortran
    - erlang
```

### --get

```
$ cat test_1.yaml | ./pscf yaml -c - --get martin.job
Developer

$ cat test_1.yaml | ./pscf yaml -c - --get martin
name: Martin D'vloper
job: Developer
skill: Elite
```

```
$ cat test_1.yaml | ./pscf yaml -c - --get fruits.[1]
Orange
```

```
$ cat test_2.yaml | ./pscf yaml -c - --get [1].tabitha.name
Tabitha Bitumen

$ cat test_2.yaml | ./pscf yaml -c - --get [1].tabitha.skills
- lisp
- fortran
- erlang

$ cat test_2.yaml | ./pscf yaml -c - --get [1].tabitha.skills.[2]
erlang
```

### --set
Update a value:
```
$ cat test_1.yaml | ./pscf yaml -c - --set martin.name=liu
...
martin:
  name: liu
  job: Developer
  skill: Elite
```

Set a new value:
```
$ cat test_1.yaml | ./pscf yaml -c - --set netpas.dns=114.114.114.114
...
netpas:
  dns: 114.114.114.114
```
Update a array:
```
$ cat test_1.yaml | ./pscf yaml -c - --set fruits.[1]=Peach
fruits:
- Apple
- Peach
- Strawberry
- Mango
...

$ cat test_1.yaml | ./pscf yaml -c - --set fruits.[+]=Peach
fruits:
- Apple
- Orange
- Strawberry
- Mango
- Peach
...
```

### --del

```
$ cat test_1.yaml | ./pscf yaml -c - --del martin.name
fruits:
- Apple
- Orange
- Strawberry
- Mango
martin:
  job: Developer
  skill: Elite
```

```
$ cat test_1.yaml | ./pscf yaml -c - --del martin
fruits:
- Apple
- Orange
- Strawberry
- Mango
```

```
$ cat test_1.yaml | ./pscf yaml -c - --del fruits.[2]
fruits:
- Apple
- Orange
- Mango
martin:
  name: Martin D'vloper
  job: Developer
  skill: Elite
```

### --getkeys

```
$ cat test_1.yaml | ./pscf yaml -c - --getkeys martin
name
job
skill
```

### --tostr
data:
```
$ cat test_3.yaml
ipv4: yes
ipv6: no
```

```
$ cat test_3.yaml | ./pscf yaml --get ipv4 -c -
true

$ cat test_3.yaml | ./pscf yaml --get ipv6 -c -
false

$ cat test_3.yaml | ./pscf yaml --get ipv4 -c - --tostr yes,no
yes

$ cat test_3.yaml | ./pscf yaml --get ipv6 -c - --tostr yes,no
no
```

### --comment
data:
```

$ cat test_4.yaml
# An employee record
martin:
  name: Martin D'vloper
  # job data
  job: Developer
  skill: Elite
```

```
$ cat ../test_4.yaml | ./pscf yaml -c - --get martin --comment
name: Martin D'vloper
# job data
job: Developer
skill: Elite


$ cat ../test_4.yaml | ./pscf yaml -c - --set martin.netpas=letvpn --comment
# An employee record
martin:
  name: Martin D'vloper
  # job data
  job: Developer
  skill: Elite
  netpas: letvpn
netpas: liu
```


