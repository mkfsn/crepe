package crepe_test

import (
	"encoding/json"
	"testing"

	"github.com/mkfsn/crepe"

	"github.com/stretchr/testify/suite"
)

type UnmarshalSuite struct {
	suite.Suite
}

func (suite *UnmarshalSuite) TestSliceOfPtrToStruct() {
	html := `
		<table></table>
		<table>
			<tbody>
			<tr data-id="aaaa" role="engineer">
				<td>Tony</td>
				<td>Male</td>
				<td>20</td>
			</tr>
			<tr data-id="bbbb" role="manager">
				<td>Mary</td>
				<td>Female</td>
				<td>23</td>
			</tr>
			</tbody>
		</table>
	`
	var data struct {
		Employees []*struct {
			Id     string `crepe:"attr=data-id"`
			Name   string `crepe:"td,eq:0,text"`
			Gender string `crepe:"td,eq:1,text"`
			Age    int    `crepe:"td,eq:2,text"`
			Role   string `crepe:"attr=role"`
		} `crepe:"table,eq:1,tbody>tr"`
	}

	err := crepe.Unmarshal([]byte(html), &data)
	suite.Require().NoError(err)

	b, err := json.Marshal(data)
	suite.Require().NoError(err)
	suite.Require().Equal(`{"Employees":[{"Id":"aaaa","Name":"Tony","Gender":"Male","Age":20,"Role":"engineer"},{"Id":"aaaa","Name":"Tony","Gender":"Male","Age":20,"Role":"engineer"}]}`, string(b))
}

func TestUnmarshal(t *testing.T) {
	suite.Run(t, new(UnmarshalSuite))
}
