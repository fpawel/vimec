package main

import (
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	robotoCondensed = "RobotoCondensed"
	titleFontSize   = 16
	datefontsize    = 14
	regularFontSize = 12
)

func savePDF() {
	pdf := gofpdf.New("P", "mm", "A4", "fonts")
	pdf.AddFont(robotoCondensed, "", "RobotoCondensed-Regular.json")
	pdf.AddFont(robotoCondensed, "B", "RobotoCondensed-Bold.json")

	tr := pdf.UnicodeTranslatorFromDescriptor("cp1251")
	//pageWidth, pageHeight := pdf.GetPageSize()

	cellFormat := func(w, h float64, align, str string) {
		pdf.CellFormat(w, h, str, "", 0, align, false, 0, "")
	}

	cellFormatLn := func(w, h float64, align, str string) {
		pdf.CellFormat(w, h, str, "", 1, align, false, 0, "")
	}

	pdf.SetLineWidth(0.3)

	for nRecord, a := range modelTableActs.items {

		if nRecord%4 == 0 {
			pdf.AddPage()
		}

		pdf.SetFont(robotoCondensed, "", titleFontSize)
		str := tr(fmt.Sprintf("Акт визуального и измерительного контроля № %d", a.ActNumber))
		cellFormatLn(0, 10, "C", str)
		pdf.SetFont(robotoCondensed, "B", datefontsize)
		cellFormatLn(0, 8, "R", tr(fmt.Sprintf("от %s", a.Date().Format("02.01.2006"))))

		pdf.SetFont(robotoCondensed, "", regularFontSize)
		str = tr("В соответствии с заявкой ОТКиИ согласно операционным картам  00226247.60103.16125,")
		cellFormatLn(pdf.GetStringWidth(str)+10, 5, "R", str)

		cellFormatLn(0, 5, "L", tr("00226247.60103.16126, 00226247.60103.16127 выполнен визуальный и измерительный контроль (ВиК)"))

		str = tr(fmt.Sprintf("магнитопровода, изделие КЭГ 9721 ИБЯЛ.304135.00%d, партия", a.DocCode))
		cellFormat(pdf.GetStringWidth(str)+1, 5, "L", str)

		pdf.SetFont(robotoCondensed, "BU", regularFontSize)
		str = tr(a.RouteSheet)
		cellFormat(pdf.GetStringWidth(str), 5, "L", str)

		pdf.SetFont(robotoCondensed, "", regularFontSize)
		cellFormatLn(0, 5, "L", tr(fmt.Sprintf(", в количестве %d штук.", a.ProductsCount)))

		pdf.SetFont(robotoCondensed, "B", regularFontSize)

		pdf.Ln(3)

		str = tr("Заключение по результатам ВиК: ")
		cellFormat(pdf.GetStringWidth(str)+10, 5, "R", str)

		pdf.SetFont(robotoCondensed, "", regularFontSize)
		cellFormatLn(0, 5, "L", tr("дефектов не выявлено, объекты контроля соответсвуют"))
		cellFormatLn(0, 5, "L", tr(fmt.Sprintf("требованиям ОСТ 4ГО.070.015, РД 03-606-03 и ИБЯЛ.30413500%d.", a.DocCode)))

		pdf.Ln(5)

		pdf.SetFont(robotoCondensed, "B", regularFontSize)
		str = tr("Контроль выполнил:")
		cellFormat(pdf.GetStringWidth(str)+1, 5, "L", str)

		pdf.SetFont(robotoCondensed, "", regularFontSize)
		str = tr("Филимоненков П.А.")
		cellFormat(pdf.GetStringWidth(str)+30, 5, "L", str)

		pdf.SetFont(robotoCondensed, "B", regularFontSize)
		str = tr("Начальник ОТКиИ:")
		cellFormat(pdf.GetStringWidth(str)+1, 5, "L", str)

		pdf.SetFont(robotoCondensed, "", regularFontSize)
		str = tr("Лемешев В.Л.")
		cellFormat(pdf.GetStringWidth(str), 5, "L", str)

		if nRecord%4 != 3 {
			pdf.Ln(12)
			pdf.MoveTo(pdf.GetX(), pdf.GetY())
			pdf.LineTo(pdf.GetX()+190, pdf.GetY())
			pdf.DrawPath("D")
			pdf.Ln(5)
		}
	}

	dir := filepath.Join(filepath.Dir(os.Args[0]), "pdf")
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(dir, os.ModePerm)
		}
		if err != nil {
			log.Fatal(err)
		}
	}

	tmpFile, err := ioutil.TempFile(dir, "*.pdf")
	if err != nil {
		log.Fatal(err)
	}
	if err := tmpFile.Close(); err != nil {
		log.Fatal(err)
	}
	if err := os.Remove(tmpFile.Name()); err != nil {
		return
	}

	if err := pdf.OutputFileAndClose(tmpFile.Name()); err != nil {
		panic(err)
	}

	if err := exec.Command("explorer.exe", tmpFile.Name()).Start(); err != nil {
		panic(err)
	}
}
