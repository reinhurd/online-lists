package yandex

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/go-resty/resty/v2"
	"online-lists/internal/models"
)

type Client struct {
	resty *resty.Client
	token string
}

func (c *Client) GetYDToken() {
	//TODO implement
}

func (c *Client) GetYDList() {
	list := models.YADISKList{}
	headers := map[string]string{
		"Accept":        "application/json",
		"Authorization": "OAuth " + c.token,
	}
	res, err := c.resty.R().SetHeaders(headers).Get("https://cloud-api.yandex.net/v1/disk/resources/files?limit=1000")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(res.Body(), &list)
	for _, item := range list.Items {
		if strings.Contains(item.Name, "МАЙН") {
			fmt.Printf("%+v\n", item)
		}
	}
}

func (c *Client) GetYDFileByPath(path string) {
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
	//download file by link
	res, err = c.resty.R().SetHeaders(headers).Get(item.File)
	err = os.WriteFile("internal/repository/tmp.xlsx", res.Body(), 0644)
	if err != nil {
		panic(err)
	}
}

func NewClient(
	r *resty.Client,
	token string,
) *Client {
	return &Client{
		resty: r,
		token: token,
	}
}
