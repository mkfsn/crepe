package crepe

import (
	"encoding/json"
	"fmt"
)

func ExampleUnmarshal() {
	html := `
		<div id="header">
			<h1>Employees</h1>
		</div>
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
		Title     string `crepe:"div#header>h1,text"`
		Employees []*struct {
			Id     string `crepe:",attr=data-id"`
			Name   string `crepe:"td,eq:0,text"`
			Gender string `crepe:"td,eq:1,text"`
			Age    int    `crepe:"td,eq:2,text"`
			Role   string `crepe:",attr=role"`
		} `crepe:"table>tbody>tr"`
	}
	if err := Unmarshal([]byte(html), &data); err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	b, _ := json.MarshalIndent(data, "", "\t")
	fmt.Printf("result: %s\n", b)
}
