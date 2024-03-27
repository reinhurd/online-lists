package helpers

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

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

func GetCSVHeaders(csvFile string) []string {
	f, err := os.Open(csvFile)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = -1
	records, err := r.ReadAll()
	if err != nil {
		fmt.Println(err)
	}

	return records[0]
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

func InsertNewValueUnderHeader(csvFile, header, value string) error {
	// TODO remove the temp file logic
	tempFileName := "tempfile.csv"

	// Open the original CSV file
	file, err := os.Open(csvFile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err
	}
	defer file.Close()

	// Create a new temp file to write modifications
	tempFile, err := os.Create(tempFileName)
	if err != nil {
		fmt.Println("Error creating temp file:", err)
		return err
	}
	defer tempFile.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	writer := csv.NewWriter(tempFile)

	// Read the headers
	headers, err := reader.Read()
	if err != nil {
		fmt.Println("Error reading headers:", err)
		return err
	}

	// Find the index of the target header
	targetIndex := -1
	for i, header1 := range headers {
		if header1 == header {
			targetIndex = i
			break
		}
	}

	if targetIndex == -1 {
		fmt.Println("Header not found")
		return err
	}

	// Write headers to the temp file
	if err = writer.Write(headers); err != nil {
		fmt.Println("Error writing headers:", err)
		return err
	}

	// Read and modify rows
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading records:", err)
		return err
	}

	isSaved := false
	for _, record := range records {
		if strings.TrimSpace(record[targetIndex]) == "" && !isSaved {
			record[targetIndex] = value
			isSaved = true
		}
		if err = writer.Write(record); err != nil {
			fmt.Println("Error writing record:", err)
			return err
		}
	}

	writer.Flush()
	if err = writer.Error(); err != nil {
		fmt.Println("Error flushing writer:", err)
		return err
	}

	// Close the temp file and original file before renaming
	tempFile.Close()
	file.Close()

	// Replace the original file with the modified file
	if err = os.Rename(tempFileName, csvFile); err != nil {
		fmt.Println("Error replacing the original file:", err)
		return err
	}
	return nil
}
