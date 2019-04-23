package main

import (

	"github.com/fpawel/vimec/internal/data"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"gopkg.in/reform.v1"
	"time"
)

type ModelTableActs struct {
	walk.TableModelBase
	items []*data.Act
}

func (m *ModelTableActs) SelectByDate(year, month int) {
	const
	(
		qYear = "WHERE year = ?"
		qMonth = " AND month = ?"
		qOrder = " ORDER BY act_number"
	)

	var (
		xs []reform.Struct
		err error )

	if year == 0 {
		xs, err = data.DB.SelectAllFrom( data.ActTable,qOrder)
	} else if month == 0{
		xs, err = data.DB.SelectAllFrom( data.ActTable,qYear + qOrder, year)
	} else {
		xs, err = data.DB.SelectAllFrom( data.ActTable,qYear + qMonth + qOrder, year, month)
	}
	if err != nil {
		panic(err)
	}
	m.applyItems(xs)
}



func (m *ModelTableActs) applyItems(xs []reform.Struct)  {
	m.items = nil
	for _, x := range xs {
		p := x.(*data.Act)
		m.items = append(m.items, p)
	}
	m.PublishRowsReset()

}

func (m *ModelTableActs) RowCount() int {
	return len(m.items)
}

func (m *ModelTableActs) Value(row, col int) interface{} {
	x := m.items[row]
	switch ActColumn(col) {
	case ActColID:
		return x.ActID
	case ActColNumber:
		return x.ActNumber
	case ActColCreatedAt:

		return time.Date(x.Year,x.Month,x.Day,0,0,0,0, time.Local)
	case ActColDoc:
		return x.DocCode
	case ActColRouteSheet:
		return x.RouteSheet
	case ActColProductsCount:
		return x.ProductsCount
	}
	return ""
}

type ActColumn int

const (
	ActColID ActColumn = iota
	ActColNumber
	ActColCreatedAt
	ActColDoc
	ActColRouteSheet
	ActColProductsCount
)

func actColumns() []TableViewColumn {
	x := make([]TableViewColumn, ActColProductsCount+1)

	type t = TableViewColumn
	x[ActColID] =
		t{Title: "ID", Width: 80}
	x[ActColNumber] =
		t{Title: "№", Width: 80}
	x[ActColCreatedAt] =
		t{Title: "Дата", Width: 150, Format: "02.01.06"}
	x[ActColDoc] =
		t{Title: "Код", Width: 80}
	x[ActColRouteSheet] =
		t{Title: "М/л", Width: 150}
	x[ActColProductsCount] =
		t{Title: "Штук", Width: 100}

	return x
}