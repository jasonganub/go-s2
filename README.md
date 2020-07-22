# go-s2id

go-s2id is a library written in Go. This contains only a client due to being entirely local.

This library consists of many utilities for S2ID such as diffing between [S2ID](https://s2geometry.io) files,
generating S2ID based on levels on geojsons from [geojson.io](https://geojson.io)

One useful functionality of this library is to be able to have two files of cell IDs and to be able to diff and remove
the unions from just one file. This is useful when you need S2ID cells for a polygon that need to be adjacent to
each other without overlapping or gaps.

## Requirement
- Golang 1.13
- Go Modules

## Example

```go
import "github.com/jasonganub/go-s2id"
```

```

## Build

```
go build
```


Reach out to me on twitter @jasonganub for any inquiries.