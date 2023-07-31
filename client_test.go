package wildberries_api

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"
)

type GeoDataTest struct {
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Address      string  `json:"address"`
	XInfo        string  `json:"xinfo"`
	UserDataSign string  `json:"userDataSign"`
	Destinations []int64 `json:"destinations"`
	Locale       string  `json:"locale"`
	Shard        int     `json:"shard"`
	Currency     string  `json:"currency"`
	IP           string  `json:"ip"`
	Dt           int     `json:"dt"`
}

func TestGeoResp(t *testing.T) {

	expectedResult := `{"latitude":55.753737,"longitude":37.6201,"address":"Москва","xinfo":"appType=1&curr=rub&dest=-1257786&regions=80,38,83,4,64,33,68,70,30,40,86,75,69,22,1,31,66,110,48,71,114&spp=0","userDataSign":"version=1&uid=0&spp=0&timestamp=1690737224&sign=0445235822bebebc692f084f33d3cc82423ea19dcd479fcc54e7134705fffa06","destinations":[-1029256,-102269,-2162196,-1257786],"locale":"ru","shard":0,"currency":"RUB","ip":"37.252.91.4","dt":0}`

	var expectedResultData GeoDataTest

	err := json.Unmarshal([]byte(expectedResult), &expectedResultData)

	if err != nil {
		t.Errorf("Error unmarshaling JSON: %v", err)
	}

	geoData, err := GetGeoData()

	if err != nil {
		t.Errorf("Error getting geo data: %v", err)
	}

	expectedResultData.UserDataSign = ""
	geoData.UserDataSign = ""

	if !reflect.DeepEqual(expectedResultData, geoData) {
		t.Errorf("Expected: %v. Got: %v", expectedResultData, geoData)
	}

}

func GetGeoData() (GeoDataTest, error) {
	url := "https://user-geo-data.wildberries.ru/get-geo-info?currency=RUB&latitude=55.753737&longitude=37.6201&locale=ru&address=%D0%9C%D0%BE%D1%81%D0%BA%D0%B2%D0%B0&dt=0"

	client := CreateClient()

	client.SetHeader("Host", "user-geo-data.wildberries.ru")

	resp, err := client.R().Get(url)

	if err != nil {
		return GeoDataTest{}, fmt.Errorf("Error to send request: %w", err)
	}

	if resp.StatusCode() != 200 {
		return GeoDataTest{}, fmt.Errorf("Status code err: %d", resp.StatusCode())
	}

	var geoData GeoDataTest

	if resp.Header().Get("Content-Encoding") == "gzip" {
		body := resp.Body()
		gzReader, err := gzip.NewReader(bytes.NewReader(body))
		if err != nil {
			return GeoDataTest{}, fmt.Errorf("Error creating gzip reader: %w", err)
		}
		defer gzReader.Close()

		decodedBody, decode_err := ioutil.ReadAll(gzReader)

		if decode_err != nil {
			return GeoDataTest{}, fmt.Errorf("Error decoding: %w", err)
		}

		err = json.Unmarshal(decodedBody, &geoData)
		if err != nil {
			return GeoDataTest{}, fmt.Errorf("Error unmarshaling JSON: %w", err)
		}
	} else {
		body := resp.Body()

		err = json.Unmarshal(body, &geoData)
		if err != nil {
			return GeoDataTest{}, fmt.Errorf("Error unmarshaling JSON: %w", err)
		}

	}

	return geoData, nil
}
