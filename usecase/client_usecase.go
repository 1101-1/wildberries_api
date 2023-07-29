package usecase

import (
	"github.com/1101-1/wildberries_api/json_data"

	"github.com/go-resty/resty/v2"
)

type APIusecase interface {
	GetGeoInfo() json_data.GeoData
}

type Client struct {
	Client  *resty.Client
	GeoData json_data.GeoData
}

func (c *Client) GetGeoInfo() json_data.GeoData {
	return c.GeoData
}
