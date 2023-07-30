package wildberries_api

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/1101-1/wildberries_api/json_data"
	"github.com/1101-1/wildberries_api/usecase"

	"github.com/go-resty/resty/v2"
)

func NewClient() (usecase.UserApi, error) {
	client := resty.New()
	client.SetTimeout(5 * time.Second)
	client.SetHeaders(map[string]string{
		"User-Agent":     "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36",
		"Accept":         "*/*",
		"Referer":        "https://www.wildberries.ru/catalog/",
		"Origin":         "https://www.wildberries.ru",
		"Connection":     "keep-alive",
		"Sec-Fetch-Dest": "empty",
		"Sec-Fetch-Mode": "cors",
		"Sec-Fetch-Site": "same-site",
		"TE":             "trailers",
	})

	for retries := 0; retries < 3; retries++ {
		geo_data, err := get_geo(client)
		if err == nil {
			return &usecase.UserClient{
				Client:  client,
				GeoData: geo_data,
			}, nil
		}

		time.Sleep(1 * time.Second)
	}

	return nil, fmt.Errorf("failed to get geo data: max retries exceeded")
}

func get_geo(client *resty.Client) (json_data.GeoData, error) {
	url := "https://user-geo-data.wildberries.ru/get-geo-info?currency=RUB&locale=ru"

	resp, err := client.SetHeader("Host", "user-geo-data.wildberries.ru").R().Get(url)
	if err != nil {
		return json_data.GeoData{}, err
	}

	var geoData json_data.GeoData
	if resp.Header().Get("Content-Encoding") == "gzip" {
		body := resp.Body()
		gzReader, err := gzip.NewReader(bytes.NewReader(body))
		if err != nil {
			return json_data.GeoData{}, fmt.Errorf("Error creating gzip reader: %w", err)
		}
		defer gzReader.Close()

		decodedBody, err := ioutil.ReadAll(gzReader)
		if err != nil {
			return json_data.GeoData{}, fmt.Errorf("Error reading decompressed data: %w", err)
		}

		err = json.Unmarshal(decodedBody, &geoData)
		if err != nil {
			return json_data.GeoData{}, fmt.Errorf("Error unmarshaling JSON: %w", err)
		}
	} else {
		body := resp.Body()

		err = json.Unmarshal(body, &geoData)
		if err != nil {
			return json_data.GeoData{}, fmt.Errorf("Error unmarshaling JSON: %w", err)
		}
	}

	return geoData, nil
}
