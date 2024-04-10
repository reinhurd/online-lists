package service

import (
	"fmt"
	"os"
	"strings"

	"online-lists/internal/clients/yandex"
	"online-lists/internal/helpers"
	"online-lists/internal/models"
)

var defaultCsvName string

type Service struct {
	yaClient *yandex.Client
}

func (s *Service) GetYaList() []string {
	return s.yaClient.GetYDList()
}

func (s *Service) GetHeaders() string {
	res := helpers.GetCSVHeaders(models.FileFolder + defaultCsvName)
	return strings.Join(res, ", ")
}

func (s *Service) SetDefaultCsv(csvName string) string {
	defaultCsvName = csvName
	return fmt.Sprintf("Set %s as default csv", csvName)
}

func (s *Service) ListCsv() string {
	files, err := helpers.GetCSVFiles()
	if err != nil {
		fmt.Println(err)
	}
	return strings.Join(files, ", ")
}

func (s *Service) Add(header, value string) string {
	var resp string
	if defaultCsvName == "" {
		resp = "Set default csv filename first"
	} else {
		err := helpers.InsertNewValueUnderHeader(models.FileFolder+defaultCsvName, header, value)
		if err != nil {
			fmt.Println(err)
		}
		resp = fmt.Sprintf("Added %s under %s", value, header)
	}
	return resp
}

func (s *Service) YAFile(filename string) string {
	defaultExcelName := "tmp.xlsx"
	s.yaClient.GetYDFileByPath(os.Getenv("YDFILE"), defaultExcelName)
	helpers.ConvertToCSV(defaultExcelName)

	return "File downloaded and converted to CSV"
}

func (s *Service) YAUpload(filename string) string {
	var resp string
	err := s.yaClient.SaveFileToYD(filename)
	if err != nil {
		resp = "Error uploading file to Yandex Disk " + err.Error()
	} else {
		resp = "File uploaded to Yandex Disk"
	}
	return resp
}

func NewService(yaClient *yandex.Client) *Service {
	return &Service{yaClient: yaClient}
}
