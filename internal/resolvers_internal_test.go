package internal

import (
	"bytes"
	"testing"

	"github.com/PuerkitoBio/goquery"
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

func (suite *ResolversSuite) TestResolveHtml() {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(`<h1><span>Title</span></h1>`)))
	suite.Require().NoError(err)

	r, err := newResolver("html")
	suite.Require().NoError(err)
	suite.Equal(r, &htmlResolver{})

	actual, err := r.Resolve(doc.Selection.Find("h1"))
	suite.Require().NoError(err)
	excepted := &Result{Raw: "<span>Title</span>", Ok: true}
	suite.Equal(excepted, actual)
}

func (suite *ResolversSuite) TestResolveAttr() {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(`<p data-content="foo"></p>`)))
	suite.Require().NoError(err)

	r, err := newResolver("attr=data-content")
	suite.Require().NoError(err)
	suite.Equal(r, &attrResolver{target: "data-content"})

	actual, err := r.Resolve(doc.Selection.Find("p"))
	suite.Require().NoError(err)
	excepted := &Result{Raw: "foo", Ok: true}
	suite.Equal(excepted, actual)
}

func (suite *ResolversSuite) TestResolveText() {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(`<p>Hello</p>`)))
	suite.Require().NoError(err)

	r, err := newResolver("text")
	suite.Require().NoError(err)
	suite.Equal(r, &textResolver{})

	actual, err := r.Resolve(doc.Selection.Find("p"))
	suite.Require().NoError(err)
	excepted := &Result{Raw: "Hello", Ok: true}
	suite.Equal(excepted, actual)
}

func TestResolvers(t *testing.T) {
	suite.Run(t, new(ResolversSuite))
}
