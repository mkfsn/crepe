package crepe

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type QuerySuite struct {
	suite.Suite
}

func (suite *QuerySuite) TestEqualQuery() {

}

func (suite *QuerySuite) TestQueryers() {
	testCases := []struct {
		name             string
		selectors        []string
		expectedQueryers []queryer
	}{
		{
			name:      "basic queryer",
			selectors: []string{"table>tr"},
			expectedQueryers: []queryer{
				&basicQueryer{selector: "table>tr"},
			},
		},
		{
			name:      "equal queryer",
			selectors: []string{":eq(3)"},
			expectedQueryers: []queryer{
				&equalQueryer{number: 3},
			},
		},
		{
			name:      "first queryer",
			selectors: []string{":first"},
			expectedQueryers: []queryer{
				&firstQueryer{},
			},
		},
		{
			name:      "last queryer",
			selectors: []string{":last"},
			expectedQueryers: []queryer{
				&lastQueryer{},
			},
		},
	}
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			actualQueryers := make([]queryer, 0)
			for _, selector := range tc.selectors {
				queryer, err := newQueryer(selector)
				suite.Require().NoError(err)
				actualQueryers = append(actualQueryers, queryer)
			}
			suite.Require().Equal(tc.expectedQueryers, actualQueryers)
		})
	}
}

func TestQuery(t *testing.T) {
	suite.Run(t, new(QuerySuite))
}
