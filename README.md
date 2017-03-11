narg - Nested arguments as configuration
========================================

[![MIT license](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![GitHub release](https://img.shields.io/github/release/nochso/narg.svg)](https://github.com/nochso/narg/releases)
[![GoDoc](https://godoc.org/github.com/nochso/narg?status.svg)](http://godoc.org/github.com/nochso/narg)
[![Go Report Card](https://goreportcard.com/badge/github.com/nochso/narg)](https://goreportcard.com/report/github.com/nochso/narg)
[![Build Status](https://travis-ci.org/nochso/narg.svg?branch=master)](https://travis-ci.org/nochso/narg)
[![Coverage Status](https://coveralls.io/repos/github/nochso/narg/badge.svg?branch=master)](https://coveralls.io/github/nochso/narg?branch=master)

Installation
------------

    go get github.com/nochso/narg


Documentation
-------------

See [godoc](https://godoc.org/github.com/nochso/narg) for API docs and examples.

### Decoding configuration

Given a configuration struct consisting of scalar types, structs and slices:

```go
type testConf struct {
	Name   string
	Port   int
	Debug  bool
	Float  float32
	Hosts  []string
	Ports  []int
	PortsU []uint8
	Admin  testUser
	Users  []testUser
}

type testUser struct {
	ID   int    `narg:"0"`
	Name string `narg:"1"`
}
```

Note that users can be defined using positional arguments instead of named ones.
The order of the arguments is defined using the struct tag `narg` with the
zero-based index (see testUser struct above).

The following narg input can be decoded into a pointer to `testConf`:

```go
in := `name foo
port 80
debug 1
float 3.14
hosts a bee "cee e"
ports 1024 1025
portsu 0xff 255
admin {
	id 1
	name "Phil \"Tandy\" Miller"
}
users {
	id 2
	name "Carol Pilbasian"
}
users 4 Todd`
cfg := &testConf{}
err := Decode(strings.NewReader(in), cfg)
```

Field names are case-insensitive.

`cfg` has now been populated with the given narg input:

```go
cfg = &testConf{
    Name:   "foo",
    Port:   80,
    Debug:  true,
    Float:  3.14,
    Hosts:  []string{"a", "bee", "cee e"},
    Ports:  []int{1024, 1025},
    PortsU: []uint8{0xff, 0xff},
    Admin: testUser{
        ID:   1,
        Name: `Phil "Tandy" Miller`,
    },
    Users: []testUser{
        {2, "Carol Pilbasian"},
        {4, "Todd"},
    },
}
```

Currently only supplied values are overwritten. You can set up your struct with
sane defaults and when decoding a sparse config the defaults will be kept.

License
-------

This package is released under the [MIT license](LICENSE).