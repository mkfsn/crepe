package crepe_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/mkfsn/crepe"

	"github.com/stretchr/testify/suite"
)

type UnmarshalSuite struct {
	suite.Suite
}

func (suite *UnmarshalSuite) TestFirstQuery() {
	html := suite.loadHTML("page1.html")
	var data struct {
		Employees []*struct {
			Id     string `crepe:"attr=data-id"`
			Name   string `crepe:"td,:eq(0),text"`
			Gender string `crepe:"td,:eq(1),text"`
			Age    int    `crepe:"td,:eq(2),text"`
			Role   string `crepe:"attr=role"`
		} `crepe:"table,:eq(1),tbody>tr,:first"`
	}
	err := crepe.Unmarshal([]byte(html), &data)
	suite.Require().NoError(err)

	b, err := json.Marshal(data)
	suite.Require().NoError(err)
	suite.Require().Equal(`{"Employees":[{"Id":"aaaa","Name":"Tony","Gender":"Male","Age":20,"Role":"engineer"}]}`, string(b))
}

func (suite *UnmarshalSuite) TestLastQuery() {
	html := suite.loadHTML("page1.html")
	var data struct {
		Employees []*struct {
			Id     string `crepe:"attr=data-id"`
			Name   string `crepe:"td,:eq(0),text"`
			Gender string `crepe:"td,:eq(1),text"`
			Age    int    `crepe:"td,:eq(2),text"`
			Role   string `crepe:"attr=role"`
		} `crepe:"table,:eq(1),tbody>tr,:last"`
	}
	err := crepe.Unmarshal([]byte(html), &data)
	suite.Require().NoError(err)

	b, err := json.Marshal(data)
	suite.Require().NoError(err)
	suite.Require().Equal(`{"Employees":[{"Id":"bbbb","Name":"Mary","Gender":"Female","Age":23,"Role":"manager"}]}`, string(b))
}

func (suite *UnmarshalSuite) TestSliceOfPtrToStruct() {
	html := suite.loadHTML("page1.html")
	var data struct {
		Employees []*struct {
			Id     string `crepe:"attr=data-id"`
			Name   string `crepe:"td,:eq(0),text"`
			Gender string `crepe:"td,:eq(1),text"`
			Age    int    `crepe:"td,:eq(2),text"`
			Role   string `crepe:"attr=role"`
		} `crepe:"table,:eq(1),tbody>tr"`
	}
	err := crepe.Unmarshal([]byte(html), &data)
	suite.Require().NoError(err)

	b, err := json.Marshal(data)
	suite.Require().NoError(err)
	suite.Require().Equal(`{"Employees":[{"Id":"aaaa","Name":"Tony","Gender":"Male","Age":20,"Role":"engineer"},{"Id":"bbbb","Name":"Mary","Gender":"Female","Age":23,"Role":"manager"}]}`, string(b))
}

func (suite *UnmarshalSuite) loadHTML(filename string) []byte {
	f, err := os.Open(fmt.Sprintf("./testdata/%s", filename))
	suite.Require().NoError(err)

	b, err := ioutil.ReadAll(f)
	suite.Require().NoError(err)

	return b
}

func TestUnmarshal(t *testing.T) {
	suite.Run(t, new(UnmarshalSuite))
}
