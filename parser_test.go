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
