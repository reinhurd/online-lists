package transport

type ListService interface {
	GetYaList() []string
	GetHeaders() string
	SetDefaultCsv(csvName string) string
	ListCsv() string
	Add(header, value string) string
	YAFile(filename, path string) string
	YAUpload(filename string) string
}
