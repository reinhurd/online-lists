package service

import (
	"fmt"
	"log"
	"os"
	"strings"

	"online-lists/internal/clients/yandex"
	"online-lists/internal/helpers"
)

var defaultCsvName string

type Service struct {
	logger     *log.Logger
	yaClient   *yandex.Client
	fileFolder string
}

func (s *Service) GetYaList() []string {
	res, err := s.yaClient.GetYDList()
	if err != nil {
		s.logger.Println(err)
		return nil
	}
	return res
}

func (s *Service) GetHeaders() string {
	if defaultCsvName == "" {
		return "Set default csv filename first"
	}
	res, err := helpers.GetCSVHeaders(s.fileFolder + defaultCsvName)
	if err != nil {
		s.logger.Println(err)
		return fmt.Sprintf("Error getting headers %s", err)
	}
	return strings.Join(res, ", ")
}

func (s *Service) SetDefaultCsv(csvName string) string {
	defaultCsvName = csvName
	return fmt.Sprintf("Set %s as default csv", csvName)
}

func (s *Service) ListCsv() string {
	files, err := helpers.GetCSVFiles()
	if err != nil {
		s.logger.Println(err)
		return fmt.Sprintf("Error getting csv files %s", err)
	}
	return strings.Join(files, ", ")
}

func (s *Service) Add(header, value string) string {
	var resp string
	if defaultCsvName == "" {
		resp = "Set default csv filename first"
	} else {
		err := helpers.InsertNewValueUnderHeader(s.fileFolder+defaultCsvName, header, value)
		if err != nil {
			s.logger.Println(err)
			return fmt.Sprintf("Error adding value %s", err)
		}
		resp = fmt.Sprintf("Added %s under %s", value, header)
	}
	return resp
}

func (s *Service) YAFile(filename string) string {
	if filename == "" {
		filename = "tmp.xlsx"
	}
	err := s.yaClient.GetYDFileByPath(os.Getenv("YDFILE"), filename)
	if err != nil {
		return fmt.Sprintf("Error downloading file from Yandex Disk %s", err)
	}
	err = helpers.ConvertToCSV(filename)
	if err != nil {
		return fmt.Sprintf("Error converting file to CSV %s", err)
	}

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

func NewService(yaClient *yandex.Client, fileFolder string) *Service {
	return &Service{
		yaClient:   yaClient,
		fileFolder: fileFolder,
	}
}
