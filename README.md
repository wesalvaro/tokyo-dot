# Tokyo Trains Dot

Tokyo trains organized into a [GraphViz](http://www.graphviz.org) Dot file.

[Preview](https://rawgit.com/wesalvaro/tokyo-dot/master/renders/index.html)

Contributions welcome!

# Data

## Line (graph with nodes only)

Each line is a subgraph named `%LINE_NAMEStations` (e.g. `HibiyaStations`). The
node color should be set as a graph default to the line's color.

## Station (node)

Stations are named `%BADGE%02NUM` (e.g. `H03`).

### Attributes

- `label` is a record format string: `{%NAME_JP|%NAME_EN}|{%BADGE|%NUM}`
- `pos` contains the position as `%LATITUDE,%LONGITUDE!` (where `!` indicates
  that the node position should not change)

## Train Type (graph with edges only)

Each train type is a subgraph named `%LINE_NAME%TRAIN_TYPE` (e.g. `ToyokoLocal).
The edge color should be set as a graph default for the train type. This graph
contains no nodes and only edges. The edges connect the stations that the train
will stop at as it travels from the first station to the last (i.e. the reverse
direction should not be added).

### Attributes

- `len` is the length (in minutes) that it takes the train to travel the edge.
- `return` may be `no` (default `yes`) to indicate this is a one-way route.
- `reserved` may be `yes` (default `no`) to indicate seats may be reserved.

## Transfers (graph with edges only)

Similar to train types, the transfers out from each line are kept in a graph
named `%LINE_NAMETransfers` (e.g. `ChiyodaTransfers`).

### Attributes

- `len` is the length (in minutes) that it takes transfer (by walking).
- `label` denotes which cars are best positioned for the transfer and, if
  applicable, the best exit to use when transferring to a different station. The
  two values should be separated by a `|`.
