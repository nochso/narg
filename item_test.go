package narg

import (
	"testing"

	"github.com/nochso/golden"
)

func getTestItem() Item {
	return Item{
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
}

func TestItem_String(t *testing.T) {
	i := getTestItem()
	actual := i.String()

	tc := golden.NewCase(t, "test-fixtures/item/item.txt")
	if *update {
		tc.Out.Update([]byte(actual))
	}
	tc.Diff(actual)
}

func TestItemSlice_String(t *testing.T) {
	doc := ItemSlice{getTestItem(), getTestItem()}
	actual := doc.String()
	tc := golden.NewCase(t, "test-fixtures/item/doc.txt")
	if *update {
		tc.Out.Update([]byte(actual))
	}
	tc.Diff(actual)
}
