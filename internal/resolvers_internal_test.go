package internal

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ResolversSuite struct {
	suite.Suite
}

func (suite *ResolversSuite) TestHtmlResolvers() {
	r, err := newResolver("html")
	suite.Require().NoError(err)
	suite.Equal(r, &htmlResolver{})
}

func (suite *ResolversSuite) TestTextResolvers() {
	r, err := newResolver("text")
	suite.Require().NoError(err)
	suite.Equal(r, &textResolver{})
}

func (suite *ResolversSuite) TestAttrResolvers() {
	r, err := newResolver("attr=title")
	suite.Require().NoError(err)
	suite.Equal(r, &attrResolver{target: "title"})
}

func (suite *ResolversSuite) TestUnknownResolvers() {
	r, err := newResolver("data=title")
	suite.Equal(ErrUnknownTarget, err)
	suite.Equal(nil, r)
}

func TestResolvers(t *testing.T) {
	suite.Run(t, new(ResolversSuite))
}
