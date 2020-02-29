package crepe

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Selector = *goquery.Selection

type queryer interface {
	query(Selector) *goquery.Selection
}

func newQueryer(selector string) (queryer, error) {
	switch {
	case strings.HasPrefix(selector, "eq:"):
		return newEqualQueryer(selector)
	}
	return newBasicQueryer(selector), nil
}

type basicQueryer struct {
	selector string
}

func newBasicQueryer(selector string) *basicQueryer {
	return &basicQueryer{selector: selector}
}

func (q *basicQueryer) String() string {
	return fmt.Sprintf("basicQueryer{selector:%q}", q.selector)
}

func (q *basicQueryer) query(selection Selector) *goquery.Selection {
	return selection.Find(q.selector)
}

type equalQueryer struct {
	number int
}

func newEqualQueryer(selector string) (*equalQueryer, error) {
	var number int
	if _, err := fmt.Sscanf(selector, "eq:%d", &number); err != nil {
		return nil, err
	}
	return &equalQueryer{number: number}, nil
}

func (q *equalQueryer) query(selection Selector) *goquery.Selection {
	return selection.Eq(q.number)
}

func (q *equalQueryer) String() string {
	return fmt.Sprintf("equalQueryer{eq:%d}", q.number)
}
