package usecase

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/1101-1/wildberries_api/json_data"

	"github.com/go-resty/resty/v2"
)

type UserApi interface {
	GetGeoInfo() json_data.GeoData
	SearchItems(keyword string) (json_data.SearchData, error)
}

type UserClient struct {
	Client  *resty.Client
	GeoData json_data.GeoData
}

func (c *UserClient) GetGeoInfo() json_data.GeoData {
	return c.GeoData
}

func (c *UserClient) SearchItems(keyword string) (json_data.SearchData, error) {
	baseURL := "https://search.wb.ru/exactmatch/ru/common/v4/search"

	queryParams := url.Values{}

	queryParams.Set("query", keyword)
	queryParams.Set("resultset", "catalog")
	queryParams.Set("sort", "popular")
	queryParams.Set("suppressSpellcheck", "false")
	queryParams.Set("uclusters", "0")

	fullURL := fmt.Sprintf("%s?%s&%s", baseURL, c.GeoData.XInfo, queryParams.Encode())

	resp, err := c.Client.R().SetHeader("Host", "search.wb.ru").Get(fullURL)

	if err != nil {
		return json_data.SearchData{}, err
	}

	var searchData json_data.SearchData

	if resp.Header().Get("Content-Encoding") == "gzip" {
		body := resp.Body()

		fmt.Println(body)
		gzReader, err := gzip.NewReader(bytes.NewReader(body))
		if err != nil {
			return json_data.SearchData{}, fmt.Errorf("Error creating gzip reader: %w", err)
		}
		defer gzReader.Close()

		decodedBody, err := ioutil.ReadAll(gzReader)
		if err != nil {
			return json_data.SearchData{}, fmt.Errorf("Error reading decompressed data: %w", err)
		}

		err = json.Unmarshal(decodedBody, &searchData)
		if err != nil {
			return json_data.SearchData{}, fmt.Errorf("Error unmarshaling JSON: %w", err)
		}
	} else {
		body := resp.Body()

		err = json.Unmarshal(body, &searchData)
		if err != nil {
			return json_data.SearchData{}, fmt.Errorf("Error unmarshaling JSON: %w", err)
		}
	}

	return searchData, nil
}
