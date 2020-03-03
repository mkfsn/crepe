package internal

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type SelectorsSuite struct {
	suite.Suite
}

func (suite *SelectorsSuite) TestSelectors() {
	testCases := []struct {
		name              string
		specs             []string
		expectedSelectors []Selector
	}{
		{
			name:  "basic selector",
			specs: []string{"table>tr"},
			expectedSelectors: []Selector{
				&basicSelector{spec: "table>tr"},
			},
		},
		{
			name:  "equal selector",
			specs: []string{":eq(3)"},
			expectedSelectors: []Selector{
				&equalSelector{number: 3},
			},
		},
		{
			name:  "first selector",
			specs: []string{":first"},
			expectedSelectors: []Selector{
				&firstSelector{},
			},
		},
		{
			name:  "last selector",
			specs: []string{":last"},
			expectedSelectors: []Selector{
				&lastSelector{},
			},
		},
	}
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			actualSelectors := make([]Selector, 0)
			for _, selector := range tc.specs {
				selector, err := newSelector(selector)
				suite.Require().NoError(err)
				actualSelectors = append(actualSelectors, selector)
			}
			suite.Require().Equal(tc.expectedSelectors, actualSelectors)
		})
	}
}

func TestSelectors(t *testing.T) {
	suite.Run(t, new(SelectorsSuite))
}
