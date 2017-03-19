package narg

import (
	"reflect"
	"testing"
	"time"

	"github.com/kylelemons/godebug/pretty"
)

type testConf struct {
	Name     string
	Port     int
	Debug    bool
	Float    float32
	Hosts    []string
	Ports    []int
	PortsU   []uint8
	Creation time.Time
	Admin    testUser
	Users    []testUser
}

type testUser struct {
	ID   int    `narg:"0"`
	Name string `narg:"1"`
}

func TestDecode(t *testing.T) {
	act := &testConf{}
	exp := &testConf{
		Name:     "foo",
		Port:     80,
		Debug:    true,
		Float:    3.14,
		Hosts:    []string{"a", "bee", "cee e"},
		Ports:    []int{1024, 1025},
		PortsU:   []uint8{0xff, 0xff},
		Creation: time.Date(1999, time.December, 31, 23, 59, 59, 0, time.UTC),
		Admin: testUser{
			ID:   1,
			Name: `Phil "Tandy" Miller`,
		},
		Users: []testUser{
			{2, "Carol Pilbasian"},
			{4, "Todd"},
		},
	}
	in := `
name foo
port 80
debug 1
float 3.14
hosts a bee "cee e"
ports 1024 1025
portsu 0xff 255
creation 1999-12-31T23:59:59Z
admin {
	id 1
	name "Phil \"Tandy\" Miller"
}
users {
	id 2
	name "Carol Pilbasian"
}
users 4 Todd
`
	err := Decode(in, act)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(exp, act) {
		t.Fatal(pretty.Compare(act, exp))
	}
}
