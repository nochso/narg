package narg

import (
	"testing"

	"github.com/nochso/golden"
)

func TestItem_String(t *testing.T) {
	i := Item{
		Name: "key",
		Args: []string{"arg1", "argument 2"},
		Children: []Item{{
			Name: "wait, there's more",
			Children: []Item{{
				Name: "foo#not-a-comment",
				Args: []string{"bar", "now with \"quotes\""},
			}},
		}, {
			Name: `!"§$%&/()=?{[]}\;:_,.-Ü*ü+ÖÄ'öä#|<>@µ`,
			Children: []Item{
				{Name: "1"},
				{Name: "2"},
			},
		}},
	}
	actual := i.String()

	tc := golden.NewCase(t, "test-fixtures/item.txt")
	if *update {
		tc.Out.Update([]byte(actual))
	}
	tc.Out.TestString(actual)
}
