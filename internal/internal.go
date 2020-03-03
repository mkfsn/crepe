package internal

import (
	"strings"
)

func NewSelectResolver(spec string) (SelectResolver, error) {
	specs := strings.Split(spec, selectorSpecSplitter)
	if len(specs) == 0 {
		return nil, ErrEmptyTag
	}
	selectorSpecs, targetOrSelectorSpec := specs[:len(specs)-1], specs[len(specs)-1]

	selectors := make([]Selector, 0, len(specs))
	for _, spec := range selectorSpecs {
		selector, err := newSelector(spec)
		if err != nil {
			return nil, err // invalid selector, early return
		}
		selectors = append(selectors, selector)
	}

	if r, err := newResolver(targetOrSelectorSpec); err == nil {
		return &chainSelector{selectors: selectors, resolver: r}, nil
	}

	// Last one is not a target, maybe a selector?
	selector, err := newSelector(targetOrSelectorSpec)
	if err != nil {
		return nil, err
	}
	return &chainSelector{selectors: append(selectors, selector)}, nil
}
