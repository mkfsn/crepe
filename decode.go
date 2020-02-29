package crepe

import (
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type decoder struct {
	// selection is like jQuery's selection.
	selection string

	// extensions are like jQuery's extension in selection.
	// For example, eq:1 equals to :eq(1) in jQuery
	extensions map[string]interface{}

	// target denotes how to extract the data from selection.
	target target
}

func newDecoderFromTag(tag string) (*decoder, error) {
	sections := strings.Split(tag, ",")
	if len(sections) < 1 {
		return nil, ErrInvalidTag
	}

	selection, options := sections[0], sections[1:]
	if len(options) == 0 {
		return &decoder{selection: selection, target: target{name: "html"}}, nil
	}

	extensionsOptions, targetOptions := make([]string, 0, len(options)), make([]string, 0, len(options))
	for _, argument := range options {
		switch {
		case strings.Contains(argument, ":"):
			extensionsOptions = append(extensionsOptions, argument)
		default:
			targetOptions = append(targetOptions, argument)
		}
	}

	extensions, target := newExtensions(extensionsOptions), newTarget(targetOptions)
	return &decoder{selection: selection, extensions: extensions, target: target}, nil
}

func (s *decoder) Delegate(selection *goquery.Selection) *goquery.Selection {
	if s.selection != "" {
		selection = selection.Find(s.selection)
	}

	for key, value := range s.extensions {
		switch key {
		case "eq":
			selection = selection.Eq(value.(int))
		}
	}
	return selection
}

type result struct {
	raw string
	ok  bool
}

func (s *decoder) Decode(selection *goquery.Selection) (*result, error) {
	selection = s.Delegate(selection)

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

type extensions map[string]interface{}

func newExtensions(arguments []string) extensions {
	extensions := make(map[string]interface{})
	for _, argument := range arguments {
		items := strings.Split(argument, ":")
		key, value := items[0], items[1]
		switch key {
		case "eq":
			if n, err := strconv.Atoi(value); err == nil {
				extensions[key] = n
			}
		}
	}
	return extensions
}

type target struct {
	// name of the target.
	// For example, html, text and attr.
	name string

	// options for the name.
	// For example, if name is attr, options could be: src, data-id, ... etc.
	options []string
}

func newTarget(arguments []string) target {
	target := target{name: "html"}
	for _, argument := range arguments {
		items := strings.Split(argument, "=")
		switch name := items[0]; name {
		case "html", "text":
			target.name = name

		case "attr":
			target.name = name
			target.options = items[1:]
		}
	}
	return target
}
