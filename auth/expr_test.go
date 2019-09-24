package auth

import (
	"testing"

	"github.com/antonmedv/expr"
)

func TestExpr(t *testing.T) {
	rule := `Token.Attr.role == "reader" && Resource.Get("kind") == "article" && Action == "GET"`

	program, err := expr.Compile(rule)
	if err != nil {
		t.Fatal(err)
	}

	output, err := expr.Run(program, &RuleContext{
		Token: &Token{
			Attr: map[string]string{
				"role": "reader",
			},
		},
		Resource: Kind("article"),
		Action:   ActionGet,
	})
	if err != nil {
		t.Fatal(err)
	}

	if !output.(bool) {
		t.Error("failed")
	}
}
