package main

import (
	"fmt"
	"github.com/fpawel/vimec/internal/data"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type mainWindow struct {
	*walk.MainWindow
}

var (
	mw             mainWindow
	modelTableActs = new(ModelTableActs)
)

func runMainWindow() error {

	var (
		tblActs                   *walk.TableView
		cbYear,
		cbMonth *walk.ComboBox
		allowComboBoxHandleChange bool
	)

	modelTableActsSelectByDate := func() {
		if !allowComboBoxHandleChange {
			return
		}
		y, _ := strconv.ParseInt(years[cbYear.CurrentIndex()], 10, 64)
		modelTableActs.SelectByDate(int(y), cbMonth.CurrentIndex())
		tblActs.EnsureItemVisible(modelTableActs.RowCount() - 1)
	}

	if err := (MainWindow{
		AssignTo:   &mw.MainWindow,
		Title:      "Акты ВИК",
		Name:       "MainWindow",
		Font:       Font{PointSize: 12, Family: "Segoe UI"},
		Background: SolidColorBrush{Color: walk.RGB(255, 255, 255)},
		Size:       Size{800, 600},
		Layout:     HBox{},

		Children: []Widget{
			TableView{
				AssignTo:                 &tblActs,
				NotSortableByHeaderClick: true,
				LastColumnStretched:      true,
				CheckBoxes:               false,
				Model:                    modelTableActs,
				Columns:                  actColumns(),
			},
			ScrollView{
				HorizontalFixed: true,
				Layout:          VBox{},
				Children: []Widget{
					GroupBox{
						Title:  "Дата",
						Layout: VBox{},
						Children: []Widget{
							Label{Text: "Год:"},
							ComboBox{
								Model:                 years,
								AssignTo:              &cbYear,
								OnCurrentIndexChanged: modelTableActsSelectByDate,
							},
							Label{Text: "Месяц:"},
							ComboBox{
								Model:                 months,
								AssignTo:              &cbMonth,
								CurrentIndex:          int(time.Now().Month()),
								OnCurrentIndexChanged: modelTableActsSelectByDate,
							},
						},
					},
					PushButton{
						Text:      "Добавить",
						OnClicked: runNewActDialog,
					},
					PushButton{
						Text:      "PDF",
						OnClicked: savePDF,
					},
					PushButton{
						Text:"Удалить PDF",
						OnClicked: func() {
							_ = os.RemoveAll(filepath.Join(filepath.Dir(os.Args[0]), "pdf"))
						},
					},
				},
			},
		},
	}).Create(); err != nil {
		return err
	}

	for i, y := range years {
		n, _ := strconv.ParseInt(y, 10, 64)
		if time.Now().Year() == int(n) {
			_ = cbYear.SetCurrentIndex(i)
			break
		}
	}
	allowComboBoxHandleChange = true
	modelTableActsSelectByDate()

	mw.Run()

	return nil
}

func runNewActDialog() {
	var (
		edProductsCount,
		ed1, ed2, ed3,
		edAct *walk.NumberEdit
		cbDocCode *walk.ComboBox
		dlg       *walk.Dialog
		btn       *walk.PushButton
		edDate    *walk.DateEdit
		lblInfo   *walk.Label
	)

	d := Dialog{
		Title:         "Добавление актов",
		Font:          Font{PointSize: 12, Family: "Segoe UI"},
		Background:    SolidColorBrush{Color: walk.RGB(255, 255, 255)},
		Layout:        Grid{Columns: 2},
		AssignTo:      &dlg,
		DefaultButton: &btn,
		CancelButton:  &btn,
		Children: []Widget{
			Label{Text: "М/л:", TextAlignment: AlignFar},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					NumberEdit{
						AssignTo: &ed1,
						Decimals: 0,
					},
					Label{Text: "-"},
					NumberEdit{
						AssignTo: &ed2,
						Decimals: 0,
					},
					Label{Text: "/"},
					NumberEdit{
						AssignTo: &ed3,
						Decimals: 0,
					},
				},
			},

			Label{Text: "Исполнение:", TextAlignment: AlignFar},
			ComboBox{
				AssignTo:     &cbDocCode,
				Model:        []string{"2", "4", "6"},
				CurrentIndex: 0,
			},

			Label{Text: "Кол-во:", TextAlignment: AlignFar},
			NumberEdit{
				AssignTo: &edProductsCount,
				Decimals: 0,
			},

			Label{Text: "Дата:", TextAlignment: AlignFar},
			DateEdit{
				Date:     time.Now(),
				AssignTo: &edDate,
			},

			Label{Text: "Акт:", TextAlignment: AlignFar},
			NumberEdit{
				AssignTo: &edAct,
				Decimals: 0,
				Value:    float64(data.NextActNumber()),
			},


			Composite{},
			PushButton{
				AssignTo: &btn,
				Text:     "Добавить",
				OnClicked: func() {
					docCode, _ := strconv.Atoi(cbDocCode.Text())
					act := data.Act{
						ProductsCount: int(edProductsCount.Value()),
						ActNumber:     int(edAct.Value()),
						DocCode:       docCode,
						Year:          edDate.Date().Year(),
						Month:         edDate.Date().Month(),
						Day:           edDate.Date().Day(),
						RouteSheet: fmt.Sprintf("%v-%v/%v",
							ed1.Value(), ed2.Value(), ed3.Value()),
					}
					err := data.DB.Save(&act)
					if err == nil {
						_ = lblInfo.SetText("")
						_ = edAct.SetValue(edAct.Value() + 1)
						modelTableActs.SelectByDate(edDate.Date().Year(), int(edDate.Date().Month()))

					} else {
						_ = lblInfo.SetText(err.Error())
					}
				},
			},

			Label{
				ColumnSpan: 2,
				AssignTo:   &lblInfo,
				TextColor:  0xFF,
			},
		},
	}
	if err := d.Create(mw.MainWindow); err != nil {
		panic(err)
		return
	}
	dlg.Run()
}

var (
	months = []string{
		"",
		"01 январь",
		"02 февраль",
		"03 март",
		"04 апрель",
		"05 май",
		"06 июнь",
		"07 июль",
		"08 август",
		"09 сентябрь",
		"10 октябрь",
		"11 ноябрь",
		"12 декабрь",
	}

	years = func() (years []string) {
		years = append(years, "")
		for _, y := range data.Years() {
			if y != time.Now().Year() {
				years = append(years, strconv.Itoa(y))
			}
		}
		years = append(years, strconv.Itoa(time.Now().Year()))
		return
	}()
)
