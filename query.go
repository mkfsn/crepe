package crepe

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type queryer interface {
	query(*goquery.Selection) *goquery.Selection
}

func newQueryer(selector string) (queryer, error) {
	switch {
	case strings.HasPrefix(selector, ":eq"):
		return newEqualQueryer(selector)
	case strings.HasPrefix(selector, ":first"):
		return newFirstQueryer(), nil
	case strings.HasPrefix(selector, ":last"):
		return newLastQueryer(), nil
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

func (q *basicQueryer) query(selection *goquery.Selection) *goquery.Selection {
	return selection.Find(q.selector)
}

type equalQueryer struct {
	number int
}

func newEqualQueryer(selector string) (*equalQueryer, error) {
	var number int
	if _, err := fmt.Sscanf(selector, ":eq(%d)", &number); err != nil {
		return nil, err
	}
	return &equalQueryer{number: number}, nil
}

func (q *equalQueryer) query(selection *goquery.Selection) *goquery.Selection {
	return selection.Eq(q.number)
}

func (q *equalQueryer) String() string {
	return fmt.Sprintf("equalQueryer{eq:%d}", q.number)
}

type firstQueryer struct{}

func newFirstQueryer() *firstQueryer {
	return &firstQueryer{}
}

func (q *firstQueryer) query(selection *goquery.Selection) *goquery.Selection {
	return selection.First()
}

type lastQueryer struct{}

func newLastQueryer() *lastQueryer {
	return &lastQueryer{}
}

func (q *lastQueryer) query(selection *goquery.Selection) *goquery.Selection {
	return selection.Last()
}
