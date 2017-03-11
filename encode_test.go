package narg

import (
	"testing"

	"github.com/kylelemons/godebug/diff"
)

func TestEncode(t *testing.T) {
	v := &testConf{
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
	exp := `name foo
port 80
debug true
float 3.14
hosts a bee "cee e"
ports 1024 1025
portsu 255 255
admin {
	id 1
	name "Phil \"Tandy\" Miller"
}
users {
	id 2
	name "Carol Pilbasian"
}
users {
	id 4
	name Todd
}`
	act, err := EncodeString(v)
	if err != nil {
		t.Fatal(err)
	}
	if exp != act {
		t.Error(diff.Diff(exp, act))
	}
}
