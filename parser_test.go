package narg

import (
	"bytes"
	"testing"

	"io/ioutil"

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

func BenchmarkParse(b *testing.B) {
	in, err := ioutil.ReadFile("test-fixtures/parser/ok/2-typical/caddy.txt")
	if err != nil {
		b.Error(err)
		return
	}
	for i := 0; i < b.N; i++ {
		_, err := Parse(bytes.NewReader(in))
		if err != nil {
			b.Fatal(err)
		}
	}
}
