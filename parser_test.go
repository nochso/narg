package narg

import (
	"testing"

	"github.com/nochso/golden"
)

func TestParse_Ok(t *testing.T) {
	golden.TestDir(t, "parser/ok", func(tc golden.Case) {
		tc.T.Parallel()
		r := tc.In.Reader()
		doc, err := Parse(r)
		r.Close()
		if err != nil {
			tc.T.Error(err)
		}
		if *update {
			tc.Out.Update([]byte(doc.String()))
		}
		tc.Diff(doc.String())
	})
}

func TestParse_Error(t *testing.T) {
	golden.TestDir(t, "parser/error", func(tc golden.Case) {
		tc.T.Parallel()
		r := tc.In.Reader()
		_, err := Parse(r)
		r.Close()
		if *update {
			tc.Out.Update([]byte(err.Error()))
		}
		tc.Diff(err.Error())
	})
}
