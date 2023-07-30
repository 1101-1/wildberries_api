package json_data

type Metadata struct {
	Name         string `json:"name"`
	CatalogType  string `json:"catalog_type"`
	CatalogValue string `json:"catalog_value"`
	NormQuery    string `json:"normquery"`
}

type Color struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type Size struct {
	Name       string `json:"name"`
	OrigName   string `json:"origName"`
	Rank       int    `json:"rank"`
	OptionID   int    `json:"optionId"`
	ReturnCost int    `json:"returnCost"`
	WH         int    `json:"wh"`
	Sign       string `json:"sign"`
}

type Product struct {
	Time1           int     `json:"time1"`
	Time2           int     `json:"time2"`
	Dist            int     `json:"dist"`
	ID              int     `json:"id"`
	Root            int     `json:"root"`
	KindID          int     `json:"kindId"`
	SubjectID       int     `json:"subjectId"`
	SubjectParentID int     `json:"subjectParentId"`
	Name            string  `json:"name"`
	Brand           string  `json:"brand"`
	BrandID         int     `json:"brandId"`
	SiteBrandID     int     `json:"siteBrandId"`
	SupplierID      int     `json:"supplierId"`
	Sale            int     `json:"sale"`
	PriceU          int     `json:"priceU"`
	SalePriceU      int     `json:"salePriceU"`
	LogisticsCost   int     `json:"logisticsCost"`
	SaleConditions  int     `json:"saleConditions"`
	ReturnCost      int     `json:"returnCost"`
	Pics            int     `json:"pics"`
	Rating          int     `json:"rating"`
	ReviewRating    float64 `json:"reviewRating"`
	Feedbacks       int     `json:"feedbacks"`
	Volume          int     `json:"volume"`
	Colors          []Color `json:"colors"`
	Sizes           []Size  `json:"sizes"`
	DiffPrice       bool    `json:"diffPrice"`
	Log             struct {
		Cpm           int    `json:"cpm"`
		Promotion     int    `json:"promotion"`
		PromoPos      int    `json:"promoPosition"`
		Position      int    `json:"position"`
		PromoTextCard string `json:"promoTextCard"`
		PromoTextCat  string `json:"promoTextCat"`
	} `json:"log"`
}

type Data struct {
	Products []Product `json:"products"`
}

type Params struct {
	Version int    `json:"version"`
	Curr    string `json:"curr"`
	Spp     int    `json:"spp"`
}

type SearchData struct {
	Metadata Metadata `json:"metadata"`
	State    int      `json:"state"`
	Version  int      `json:"version"`
	Params   Params   `json:"params"`
	Data     Data     `json:"data"`
}
