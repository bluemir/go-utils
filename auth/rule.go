package auth

import (
	"fmt"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
)

func (m *manager) Rules() RulesClause {
	return &rulesClause{m}
}

type rulesClause struct {
	*manager
}

func (c *rulesClause) Add(rule string) error {
	return c.manager.rules.Add(rule)
}
func (*rulesClause) List() ([]Rule, error) {
	return nil, nil // TODO
}
func (c *rulesClause) Delete(index int) error {
	return nil // TODO
}

type RuleContext struct {
	Token    *Token
	User     *User
	Resource Resource
	Action   Action
}

type Rules struct {
	items []Rule
}

func (rules *Rules) Add(rule string) error {
	program, err := expr.Compile(rule)
	if err != nil {
		return err
	}
	rules.items = append(rules.items, program)
	return nil
}

func (rules *Rules) check(rc *RuleContext) bool {
	rc.User = rc.Token.User

	for _, p := range rules.items {
		output, err := expr.Run(p, rc)
		if err != nil {
			// TODO log
			fmt.Println(err)
			continue
		}
		if output.(bool) {
			return true
		}
	}
	return false
}

type Rule *vm.Program
