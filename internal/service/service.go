package service

import "online-lists/internal/clients/yandex"

type Service struct {
	yaClient *yandex.Client
}

func (s *Service) GetYaList() []string {
	return s.yaClient.GetYDList()
}

func NewService(yaClient *yandex.Client) *Service {
	return &Service{yaClient: yaClient}
}
