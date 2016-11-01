package narg

import (
	"testing"

	"github.com/nochso/golden"
)

func TestParse_Ok(t *testing.T) {
	for tc := range golden.Dir(t, "parser/ok") {
		r := tc.In.Reader()
		doc, err := Parse(r)
		r.Close()
		if err != nil {
			t.Errorf("%s: %s", tc.In.Path, err)
		}
		if *update {
			tc.Out.Update([]byte(doc.String()))
		}
		tc.Diff(doc.String())
	}
}

func TestParse_Error(t *testing.T) {
	for tc := range golden.Dir(t, "parser/error") {
		r := tc.In.Reader()
		_, err := Parse(r)
		r.Close()
		if *update {
			tc.Out.Update([]byte(err.Error()))
		}
		tc.Diff(err.Error())
	}
}
