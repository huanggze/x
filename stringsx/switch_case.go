package stringsx

import (
	"fmt"
	"slices"
	"strings"
)

type (
	RegisteredCases struct {
		cases  []string
		actual string
	}
	errUnknownCase struct {
		*RegisteredCases
	}
)

func SwitchExact(actual string) *RegisteredCases {
	return &RegisteredCases{
		actual: actual,
	}
}

func (r *RegisteredCases) AddCase(cases ...string) bool {
	r.cases = append(r.cases, cases...)
	return slices.Contains(cases, r.actual)
}

func (r *RegisteredCases) String() string {
	return "[" + strings.Join(r.cases, ", ") + "]"
}

func (r *RegisteredCases) ToUnknownCaseErr() error {
	return errUnknownCase{r}
}

func (e errUnknownCase) Error() string {
	return fmt.Sprintf("expected one of %s but got %s", e.String(), e.actual)
}

func (e errUnknownCase) Is(err error) bool {
	_, ok := err.(errUnknownCase)
	return ok
}
