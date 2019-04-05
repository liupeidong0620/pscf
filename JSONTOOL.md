## pscf json Instructions for use

```shell
$ ./releases/pscf json -h
NAME:
   pscf json - Provides tools for set, delete, and get operations for a node in JSON format.

USAGE:
   pscf json [-c config.yaml] [--set node.node1.node2=value]
   pscf json [-c config.yaml] [--del node.node1.node2]
   pscf json [-c config.yaml] [--get node.node1.node2]

OPTIONS:
   --config FILE, -c FILE  Load configuration from FILE
   --del value             Delete a node in the configuration file.(--del node)
   --get value             Get the value of the node in the configuration file.(--get node)
   --set value             Set the value of the node and add a new node if the node does not exist.(--set node=value)
   --verbose, -v           Debug mode.(--verbose, -v)
```

## Examples

### --get 

```
$ echo '{"name":{"first":"Tom","last":"Smith"}}' | ./pscf json -c - --get name.last
Smith
```

```
$ echo '{"name":{"first":"Tom","last":"Smith"}}' | ./pscf json -c - --get name
{
  "first": "Tom",
  "last": "Smith"
}
```

```
$ echo '{"name":{"first":"Tom","last":"Smith"}}' | ./pscf json -c - --get name2
null
```

```
$ echo '{"friends":["Tom","Jane","Carol"]}' | ./pscf json -c - --get friends.[1]
Jane
```

### --set

Update a value:
```
$ echo '{"name":{"first":"Tom","last":"Smith"}}' | ./pscf json -c - --set name.first=Andy
{
  "name": {
    "first": "Andy",
    "last": "Smith"
  }
}
```

Set a new value:
```
$ echo '{"name":{"first":"Tom","last":"Smith"}}' | ./pscf json -c - --set age=46
{
  "age": 46,
  "name": {
    "first": "Tom",
    "last": "Smith"
  }
}

$ echo '{"name":{"first":"Tom","last":"Smith"}}' | ./pscf json -c - --set 'age="46"'
{
  "age": "46",
  "name": {
    "first": "Tom",
    "last": "Smith"
  }
}
```

Set a new nested value:
```
$ echo '{"name":{"first":"Tom","last":"Smith"}}' | ./pscf json -c - --set task.today=relax
{
  "task": {
    "today": "relax"
  },
  "name": {
    "first": "Tom",
    "last": "Smith"
  }
}
```

Replace an array value by index:
```
$ echo '{"friends":["Tom","Jane","Carol"]}' | ./pscf json -c - --set friends.[1]=Andy
{
  "friends": ["Tom", "Andy", "Carol"]
}
```

Append an array:
```
$ echo '{"friends":["Tom","Jane","Carol"]}' | ./pscf json -c - --set friends.[+]=Andy
{
  "friends": ["Tom", "Jane", "Carol", "Andy"]
}
```

Set an array value that's past the bounds:
```
$ echo '{"friends":["Tom","Jane","Carol"]}' | ./pscf json -c - --set friends.5=liu
{
  "friends": ["Tom", "Jane", "Carol", null, null, "liu"]
}
```

### --del

```
$ echo '{"age":46,"name":{"first":"Tom","last":"Smith"}}' | ./pscf json -c - --del age
{
  "name": {
    "first": "Tom",
    "last": "Smith"
  }
}
```

```
$ echo '{"friends":["Andy","Carol"]}' | ./pscf json -c - --del friends.[0]
{
  "friends": ["Carol"]
}
```
