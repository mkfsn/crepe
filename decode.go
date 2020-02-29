package crepe

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type decoder struct {
	// sequenced queryers from the structure tag
	queryers []queryer

	// target denotes how to extract the data from selection.
	target *target
}

func newDecoderFromTag(tag string) (*decoder, error) {
	sections := strings.Split(tag, tagSplitter)
	if len(sections) == 0 {
		return nil, ErrEmptyTag
	}
	selectors, targetOrSelector := sections[:len(sections)-1], sections[len(sections)-1]

	queryers := make([]queryer, 0, len(sections))
	for _, selector := range selectors {
		queryer, err := newQueryer(selector)
		if err != nil {
			return nil, err // invalid selector, early return
		}
		queryers = append(queryers, queryer)
	}

	if target, err := newTarget(targetOrSelector); err == nil {
		return &decoder{queryers: queryers, target: target}, nil
	}

	// Last one is not a target, maybe a selector?
	queryer, err := newQueryer(targetOrSelector)
	if err != nil {
		return nil, err
	}
	return &decoder{queryers: append(queryers, queryer)}, nil
}

func (s *decoder) Query(selection *goquery.Selection) *goquery.Selection {
	for _, q := range s.queryers {
		selection = q.query(selection)
	}
	return selection
}

type result struct {
	raw string
	ok  bool
}

func (s *decoder) Decode(parentSelection *goquery.Selection) (*result, error) {
	selection := s.Query(parentSelection)
	switch s.target.name {
	case "html":
		html, err := selection.Html()
		return &result{raw: html}, err
	case "text":
		return &result{raw: selection.Text(), ok: true}, nil
	case "attr":
		attr, ok := selection.Attr(s.target.options[0])
		return &result{raw: attr, ok: ok}, nil
	}
	return nil, ErrUnknownAction
}

type target struct {
	// name of the target.
	// For example, html, text and attr.
	name string

	// options for the name.
	// For example, if name is attr, options could be: src, data-id, ... etc.
	options []string
}

func newTarget(selector string) (*target, error) {
	items := strings.Split(selector, "=")
	switch v := items[0]; v {
	case "html", "":
		return &target{name: "html"}, nil
	case "text":
		return &target{name: v}, nil
	case "attr":
		return &target{name: v, options: items[1:]}, nil
	}
	return nil, ErrUnknownTarget
}
