package helpers

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/xuri/excelize/v2"
)

func ReadXLSX() {
	f, err := excelize.OpenFile("tmp.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	sl := f.GetSheetList()
	fmt.Println(sl)
	// Get all the rows in the Sheet1.
	rows, err := f.GetRows("СПИСКЕН")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, row := range rows {
		for _, colCell := range row {
			fmt.Print(colCell, "\t")
		}
		fmt.Println()
	}
}

func ConvertToCSV() {
	f := openExcel("internal/repository/tmp.xlsx")

	worksheets := f.GetSheetList()

	for i := range worksheets {
		createCSV(f, worksheets[i])
	}
}

func openExcel(fileName string) *excelize.File {
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer func() {
		// Close the spreadsheet.
		if err = f.Close(); err != nil {
			fmt.Println(err)
		}
		return
	}()

	return f
}

func createCSV(f *excelize.File, worksheet string) {
	allRows, arErr := f.GetRows(worksheet)
	if arErr != nil {
		panic(arErr)
	}
	//don't write if sheet is empty
	if len(allRows) == 0 {
		return
	}

	csvFile, csvErr := os.Create("internal/repository/" + transliterateCyrillicToEnglish(worksheet) + ".csv")
	if csvErr != nil {
		fmt.Println(csvErr)
	}
	defer func() {
		if csvErr = csvFile.Close(); csvErr != nil {
			panic(csvErr)
		}
	}()

	writer := csv.NewWriter(csvFile)

	var writerErr = writer.WriteAll(allRows)
	if writerErr != nil {
		fmt.Println(writerErr)
	}
}
