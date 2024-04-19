package service

import (
	"fmt"
	"os"
	"strings"

	"online-lists/internal/clients/yandex"
	"online-lists/internal/helpers"

	"github.com/rs/zerolog/log"
)

var defaultCsvName string

type Service struct {
	yaClient   *yandex.Client
	fileFolder string
}

func (s *Service) GetYaList() []string {
	//todo show paths and names in response, not all info
	res, err := s.yaClient.GetYDList()
	if err != nil {
		log.Error().Err(err).Msg("Error getting list from Yandex Disk")
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
		log.Error().Err(err).Msg("Error getting headers")
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
		log.Error().Err(err).Msg("Error getting csv files")
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
			log.Error().Err(err).Msg("Error adding value")
			return fmt.Sprintf("Error adding value %s", err)
		}
		resp = fmt.Sprintf("Added %s under %s", value, header)
	}
	return resp
}

func (s *Service) YAFile(filename, path string) string {
	// todo add possibility to save not only xlsx files
	if filename == "" {
		filename = "tmp.xlsx"
	}
	if path == "" {
		path = os.Getenv("YDFILE")
	}
	err := s.yaClient.GetYDFileByPath(path, filename)
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
