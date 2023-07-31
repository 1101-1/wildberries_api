package usecase_test

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/1101-1/wildberries_api"
	"github.com/1101-1/wildberries_api/json_data"
	"github.com/go-resty/resty/v2"
)

type UserApi interface {
	GetGeoInfo() json_data.GeoData
	SearchItems(keyword string) (json_data.SearchDataTest, error)
}

type UserClient struct {
	Client  *resty.Client
	GeoData json_data.GeoData
}

func (c *UserClient) GetGeoInfo() json_data.GeoData {
	return c.GeoData
}

func (c *UserClient) SearchItems(keyword string) (json_data.SearchDataTest, error) {
	baseURL := "https://search.wb.ru/exactmatch/ru/common/v4/search"

	queryParams := url.Values{}

	queryParams.Set("query", keyword)
	queryParams.Set("resultset", "catalog")
	queryParams.Set("sort", "popular")
	queryParams.Set("suppressSpellcheck", "false")
	queryParams.Set("uclusters", "0")

	fullURL := fmt.Sprintf("%s?%s&%s", baseURL, c.GeoData.XInfo, queryParams.Encode())

	resp, err := c.Client.SetHeader("Host", "search.wb.ru").R().Get(fullURL)

	if err != nil {
		return json_data.SearchDataTest{}, err
	}

	var searchData json_data.SearchDataTest

	if resp.Header().Get("Content-Encoding") == "gzip" {
		body := resp.Body()

		fmt.Println(body)
		gzReader, err := gzip.NewReader(bytes.NewReader(body))
		if err != nil {
			return json_data.SearchDataTest{}, fmt.Errorf("Error creating gzip reader: %w", err)
		}
		defer gzReader.Close()

		decodedBody, err := ioutil.ReadAll(gzReader)
		if err != nil {
			return json_data.SearchDataTest{}, fmt.Errorf("Error reading decompressed data: %w", err)
		}

		err = json.Unmarshal(decodedBody, &searchData)
		if err != nil {
			return json_data.SearchDataTest{}, fmt.Errorf("Error unmarshaling JSON: %w", err)
		}
	} else {
		body := resp.Body()

		err = json.Unmarshal(body, &searchData)
		if err != nil {
			return json_data.SearchDataTest{}, fmt.Errorf("Error unmarshaling JSON: %w", err)
		}
	}

	return searchData, nil
}

func TestSearch(t *testing.T) {
	client := wildberries_api.CreateClient()

	var testClient UserApi
	var test_err error
	for retries := 0; retries < 3; retries++ {
		geo_data, err := GetGeo(client)
		if err == nil {
			testClient =
				&UserClient{
					Client:  client,
					GeoData: geo_data,
				}
			resp, err := testClient.SearchItems("подушка")

			if err != nil {
				t.Errorf("failed to search item. Caused by %v", err)
			}

			assert_err := AssertationData(resp)

			if assert_err != nil {
				t.Errorf("Wrong search data. Caused by %v", err)
			}

			return
		}
		test_err = err
		time.Sleep(1 * time.Second)
	}

	t.Errorf("failed to get geo data: max retries exceeded. Caused by %v", test_err)
}

func AssertationData(data json_data.SearchDataTest) error {
	var expectedResultData json_data.SearchDataTest
	expected_data := `{
        "metadata": {
            "name": "подушка",
            "catalog_type": "preset",
            "catalog_value": "preset=10213203",
            "normquery": "подушка"
        },
        "state": 0,
        "version": 2,
        "params": {
            "version": 1,
            "curr": "rub",
            "spp": 0
        }
    }`

	err := json.Unmarshal([]byte(expected_data), &expectedResultData)

	if err != nil {
		return fmt.Errorf("Error unmarshaling JSON: %v", err)
	}

	if !reflect.DeepEqual(expectedResultData, data) {
		return fmt.Errorf("Expected: %v. Got: %v", expectedResultData, data)
	}

	return nil
}

func GetGeo(client *resty.Client) (json_data.GeoData, error) {
	url := "https://user-geo-data.wildberries.ru/get-geo-info?currency=RUB&latitude=55.753737&longitude=37.6201&locale=ru&address=%D0%9C%D0%BE%D1%81%D0%BA%D0%B2%D0%B0&dt=0"

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
