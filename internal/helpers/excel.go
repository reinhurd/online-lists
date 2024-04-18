package helpers

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/xuri/excelize/v2"
	"online-lists/internal/config"
)

func ReadXLSX(sheetname string) error {
	f, err := excelize.OpenFile("tmp.xlsx")
	if err != nil {
		return err
	}
	defer func() {
		if err = f.Close(); err != nil {
			panic(err)
		}
	}()
	sl := f.GetSheetList()
	log.Info().Msgf("Sheets: %v", sl)

	rows, err := f.GetRows(sheetname)
	if err != nil {
		return err
	}
	for _, row := range rows {
		for _, colCell := range row {
			log.Info().Msgf("%s\t", colCell)
		}
	}
	return nil
}

func ConvertToCSV(excelName string) error {
	f, err := openExcel(config.FileFolder + excelName)
	if err != nil {
		return err
	}

	worksheets := f.GetSheetList()

	for i := range worksheets {
		err = createCSV(f, worksheets[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func openExcel(fileName string) (*excelize.File, error) {
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = f.Close(); err != nil {
			log.Err(err).Msg("Error closing file")
		}
	}()

	return f, nil
}

func createCSV(f *excelize.File, worksheet string) error {
	allRows, arErr := f.GetRows(worksheet)
	if arErr != nil {
		return arErr
	}
	//don't write if sheet is empty
	if len(allRows) == 0 {
		return errors.New("sheet is empty")
	}

	csvFile, csvErr := os.Create(config.FileFolder + transliterateCyrillicToEnglish(worksheet) + ".csv")
	if csvErr != nil {
		return csvErr
	}
	defer func() {
		if csvErr = csvFile.Close(); csvErr != nil {
			panic(csvErr)
		}
	}()

	writer := csv.NewWriter(csvFile)

	var writerErr = writer.WriteAll(allRows)
	if writerErr != nil {
		return writerErr
	}
	return nil
}

func GetCSVHeaders(csvFile string) ([]string, error) {
	f, err := os.Open(csvFile)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = f.Close()
		if err != nil {
			log.Err(err).Msg("Error closing file")
		}
	}()

	r := csv.NewReader(f)
	r.FieldsPerRecord = -1
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	return records[0], nil
}

func ConvertCSVtoXLSX(csvFile, xlsxFile string) error {
	f, err := os.Open(csvFile)
	if err != nil {
		return err
	}
	defer func() {
		err = f.Close()
		if err != nil {
			log.Err(err).Msg("Error closing file")
		}
	}()

	r := csv.NewReader(f)
	// Adjust the CSV reader settings if necessary (e.g., different delimiter)
	// r.Comma = ';' // If your CSV uses semicolons
	r.FieldsPerRecord = -1
	records, err := r.ReadAll()
	if err != nil {
		return err
	}

	xlsx := excelize.NewFile()
	index, _ := xlsx.NewSheet("Sheet1")
	xlsx.SetActiveSheet(index)

	for i, record := range records {
		for j, field := range record {
			cell, _ := excelize.CoordinatesToCellName(j+1, i+1)
			err = xlsx.SetCellValue("Sheet1", cell, field)
			if err != nil {
				return err
			}
		}
	}

	if err = xlsx.SaveAs(xlsxFile); err != nil {
		return err
	}

	return nil
}

func InsertNewValueUnderHeader(csvFile, header, value string) error {
	tempFileName := "tempfile.csv"

	file, err := os.Open(csvFile)
	if err != nil {
		log.Err(err).Msg("Error opening file")
		return err
	}
	defer func() {
		err = file.Close()
		if err != nil {
			log.Err(err).Msg("Error closing file")
		}
	}()

	tempFile, err := os.Create(tempFileName)
	if err != nil {
		log.Err(err).Msg("Error creating temp file")
		return err
	}
	defer func() {
		err = tempFile.Close()
		if err != nil {
			log.Err(err).Msg("Error closing temp file")
		}
	}()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	writer := csv.NewWriter(tempFile)

	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("error reading headers: %w", err)
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
		return fmt.Errorf("header not found: %s", header)
	}

	if err = writer.Write(headers); err != nil {
		return fmt.Errorf("error writing headers: %w", err)
	}

	records, err := reader.ReadAll()
	if err != nil {
		log.Err(err).Msg("Error reading records")
		return err
	}

	isSaved := false
	for _, record := range records {
		if strings.TrimSpace(record[targetIndex]) == "" && !isSaved {
			record[targetIndex] = value
			isSaved = true
		}
		if err = writer.Write(record); err != nil {
			return fmt.Errorf("error writing record: %w", err)
		}
	}

	writer.Flush()
	if err = writer.Error(); err != nil {
		return fmt.Errorf("error flushing writer: %w", err)
	}

	err = tempFile.Close()
	if err != nil {
		return err
	}
	err = file.Close()
	if err != nil {
		return err
	}

	if err = os.Rename(tempFileName, csvFile); err != nil {
		return fmt.Errorf("error renaming file: %w", err)
	}
	return nil
}

func GetCSVFiles() ([]string, error) {
	files, err := os.ReadDir(config.FileFolder)
	if err != nil {
		return nil, err
	}

	var csvFiles []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".csv") {
			csvFiles = append(csvFiles, file.Name())
		}
	}

	return csvFiles, nil
}
