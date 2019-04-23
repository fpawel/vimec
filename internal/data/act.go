package data

import "time"

//go:generate reform

// Act represents a row in act table.
//reform:act
type Act struct {
	ActID         int64      `reform:"act_id,pk"`
	Year          int        `reform:"year"`
	Month         time.Month `reform:"month"`
	Day           int        `reform:"day"`
	ActNumber     int        `reform:"act_number"`
	DocCode       int        `reform:"doc_code"`
	RouteSheet    string     `reform:"route_sheet"`
	ProductsCount int        `reform:"products_count"`
}

