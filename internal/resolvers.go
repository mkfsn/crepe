package internal

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Result struct {
	Raw string
	Ok  bool
}

type Resolver interface {
	Resolve(*goquery.Selection) (*Result, error) // Or change to (ResolvedResult, error)
}

func newResolver(spec string) (Resolver, error) {
	items := strings.Split(spec, "=")
	switch v := items[0]; v {
	case "html", "":
		return &htmlResolver{}, nil
	case "text":
		return &textResolver{}, nil
	case "attr":
		return &attrResolver{target: items[1]}, nil
	}
	return nil, ErrUnknownTarget
}

type htmlResolver struct{}

func (r *htmlResolver) Resolve(s *goquery.Selection) (*Result, error) {
	html, err := s.Html()
	return &Result{Raw: html, Ok: err == nil}, err
}

type textResolver struct{}

func (r *textResolver) Resolve(s *goquery.Selection) (*Result, error) {
	return &Result{Raw: s.Text(), Ok: true}, nil
}

type attrResolver struct {
	target string
}

func (r *attrResolver) Resolve(s *goquery.Selection) (*Result, error) {
	attr, ok := s.Attr(r.target)
	return &Result{Raw: attr, Ok: ok}, nil
}
