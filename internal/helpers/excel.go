package helpers

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/xuri/excelize/v2"
)

func ReadXLSX(sheetname string) {
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
	rows, err := f.GetRows(sheetname)
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

func ConvertCSVtoXLSX(csvFile, xlsxFile string) error {
	f, err := os.Open(csvFile)
	if err != nil {
		return err
	}
	defer f.Close()

	// Create a new reader for the CSV file
	r := csv.NewReader(f)
	// Adjust the CSV reader settings if necessary (e.g., different delimiter)
	// r.Comma = ';' // If your CSV uses semicolons
	r.FieldsPerRecord = -1
	// Read all records at once
	records, err := r.ReadAll()
	if err != nil {
		return err
	}

	// Create a new Excel file
	xlsx := excelize.NewFile()
	// Create a new sheet named "Sheet1"
	index, _ := xlsx.NewSheet("Sheet1")
	// Set the active sheet of the workbook
	xlsx.SetActiveSheet(index)

	fmt.Println(records)
	// Iterate through records to populate the sheet
	for i, record := range records {
		for j, field := range record {
			cell, _ := excelize.CoordinatesToCellName(j+1, i+1)
			xlsx.SetCellValue("Sheet1", cell, field)
		}
	}

	// Save the XLSX file
	if err = xlsx.SaveAs(xlsxFile); err != nil {
		return err
	}

	return nil
}
