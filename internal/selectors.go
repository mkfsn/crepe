package internal

import (
	"errors"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	selectorSpecSplitter = ","
)

type Selector interface {
	Select(*goquery.Selection) *goquery.Selection
}

type SelectResolver interface {
	Selector
	Resolver
}

type chainSelector struct {
	selectors []Selector
	resolver  Resolver
}

func (c *chainSelector) Select(s *goquery.Selection) *goquery.Selection {
	for _, selector := range c.selectors {
		s = selector.Select(s)
	}
	return s
}

func (c *chainSelector) Resolve(s *goquery.Selection) (*Result, error) {
	if c.resolver == nil {
		return nil, errors.New("nothing to resolve")
	}
	s = c.Select(s)
	return c.resolver.Resolve(s)
}

func newSelector(spec string) (Selector, error) {
	switch {
	case strings.HasPrefix(spec, ":eq"):
		return newEqualSelector(spec)
	case strings.HasPrefix(spec, ":first"):
		return newFirstSelector(), nil
	case strings.HasPrefix(spec, ":last"):
		return newLastSelector(), nil
	}
	return newBasicSelector(spec), nil
}

type basicSelector struct {
	spec string
}

func newBasicSelector(spec string) *basicSelector {
	return &basicSelector{spec: spec}
}

func (q *basicSelector) String() string {
	return fmt.Sprintf("basicSelector{spec:%q}", q.spec)
}

func (q *basicSelector) Select(selection *goquery.Selection) *goquery.Selection {
	return selection.Find(q.spec)
}

type equalSelector struct {
	number int
}

func newEqualSelector(selector string) (*equalSelector, error) {
	var number int
	if _, err := fmt.Sscanf(selector, ":eq(%d)", &number); err != nil {
		return nil, err
	}
	return &equalSelector{number: number}, nil
}

func (q *equalSelector) Select(selection *goquery.Selection) *goquery.Selection {
	return selection.Eq(q.number)
}

func (q *equalSelector) String() string {
	return fmt.Sprintf("equalSelector{eq:%d}", q.number)
}

type firstSelector struct{}

func newFirstSelector() *firstSelector {
	return &firstSelector{}
}

func (q *firstSelector) Select(selection *goquery.Selection) *goquery.Selection {
	return selection.First()
}

type lastSelector struct{}

func newLastSelector() *lastSelector {
	return &lastSelector{}
}

func (q *lastSelector) Select(selection *goquery.Selection) *goquery.Selection {
	return selection.Last()
}
