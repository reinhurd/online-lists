package yandex

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
	"online-lists/internal/models"
)

const YAFileCountLimit = "1000"

type Client struct {
	resty      *resty.Client
	token      string
	fileFolder string
}

func (c *Client) GetYDToken() {
	//TODO implement
}

func (c *Client) GetYDList() []string {
	list := models.YADISKList{}
	headers := map[string]string{
		"Accept":        "application/json",
		"Authorization": "OAuth " + c.token,
	}
	res, err := c.resty.R().SetHeaders(headers).Get("https://cloud-api.yandex.net/v1/disk/resources/files?limit=" + YAFileCountLimit)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(res.Body(), &list)
	if err != nil {
		panic(err)
	}
	var names []string
	for _, item := range list.Items {
		names = append(names, item.Name)
	}
	return names
}

func (c *Client) GetYDFileByPath(path, defaultExcelName string) {
	item := models.YDItem{}
	headers := map[string]string{
		"Accept":        "application/json",
		"Authorization": "OAuth " + c.token,
	}
	res, err := c.resty.R().SetHeaders(headers).Get("https://cloud-api.yandex.net/v1/disk/resources?path=" + path)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(res.Body(), &item)
	if err != nil {
		panic(err)
	}
	//download file by link
	res, err = c.resty.R().SetHeaders(headers).Get(item.File)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(c.fileFolder+defaultExcelName, res.Body(), 0644)
	if err != nil {
		panic(err)
	}
}

func (c *Client) SaveFileToYD(filename string) error {
	fileData, err := os.Open(c.fileFolder + filename)
	if err != nil {
		panic(err)
	}
	headers := map[string]string{
		"Accept":        "application/json",
		"Authorization": "OAuth " + c.token,
	}
	pathToUpload := "disk:/" + filename
	respUrl := models.YDUploadResponse{}
	urlToUpload, err := c.resty.R().SetHeaders(headers).Get("https://cloud-api.yandex.net/v1/disk/resources/upload?path=" + pathToUpload)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(urlToUpload.Body(), &respUrl)
	if err != nil {
		panic(err)
	}
	//todo deal with unsupported protocol error when uploading file
	put, err := c.resty.R().SetHeaders(headers).SetBody(fileData).Put(respUrl.Href)
	if err != nil {
		fmt.Println(err)
		fmt.Printf("%+v", put)
	}
	return err
}

func NewClient(
	r *resty.Client,
	token string,
	fileFolder string,
) *Client {
	return &Client{
		resty:      r,
		token:      token,
		fileFolder: fileFolder,
	}
}
