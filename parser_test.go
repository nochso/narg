package narg

import (
	"io/ioutil"
	"testing"

	"github.com/nochso/golden"
)

func TestParse_Ok(t *testing.T) {
	golden.TestDir(t, "parser/ok", func(tc golden.Case) {
		tc.T.Parallel()
		doc, err := Parse(tc.In.String())
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
		_, err := Parse(tc.In.String())
		if *update {
			tc.Out.Update([]byte(err.Error()))
		}
		tc.Diff(err.Error())
	})
}

func BenchmarkParse(b *testing.B) {
	f, err := ioutil.ReadFile("test-fixtures/parser/ok/2-typical/caddy.txt")
	if err != nil {
		b.Error(err)
		return
	}
	in := string(f)
	for i := 0; i < b.N; i++ {
		_, err := Parse(in)
		if err != nil {
			b.Fatal(err)
		}
	}
}
