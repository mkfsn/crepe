# crepe

crepe is a tool that helps you extract the data from ~~tons-of-layers~~ HTML in an easier way.

This uses a brilliant project called [goquery](https://github.com/PuerkitoBio/goquery) to extract the data.


# Example

By specifying some selectors in struct tags, crepe helps you unmarshal the data from HTML.

```go
	html := `
		<div id="header">
			<h1>Employees</h1>
		</div>
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
		Title     string `crepe:"div#header>h1,text"`
		Employees []*struct {
			Id     string `crepe:"attr=data-id"`
			Name   string `crepe:"td,eq:0,text"`
			Gender string `crepe:"td,eq:1,text"`
			Age    int    `crepe:"td,eq:2,text"`
			Role   string `crepe:"attr=role"`
		} `crepe:"table,eq:1,tbody>tr"`
	}
	if err := Unmarshal([]byte(html), &data); err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	fmt.Printf("result: %v\n", data)
```