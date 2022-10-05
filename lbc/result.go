package lbc

import "fmt"

type Result struct {
	Total          int  `json:"total,omitempty"`
	TotalAll       int  `json:"total_all,omitempty"`
	TotalPro       int  `json:"total_pro,omitempty"`
	TotalPrivate   int  `json:"total_private,omitempty"`
	TotalActive    int  `json:"total_active,omitempty"`
	TotalInactive  int  `json:"total_inactive,omitempty"`
	TotalShippable int  `json:"total_shippable,omitempty"`
	MaxPages       int  `json:"max_pages,omitempty"`
	Ads            []Ad `json:"ads,omitempty"`
}

type Ad struct {
	AdType               string       `json:"ad_type,omitempty"`
	Attributes           []Attributes `json:"attributes,omitempty"`
	Body                 string       `json:"body,omitempty"`
	CategoryID           string       `json:"category_id,omitempty"`
	CategoryName         string       `json:"category_name,omitempty"`
	FirstPublicationDate string       `json:"first_publication_date,omitempty"`
	HasPhone             bool         `json:"has_phone,omitempty"`
	Images               Images       `json:"images,omitempty"`
	IndexDate            string       `json:"index_date,omitempty"`
	ListID               int64        `json:"list_id,omitempty"`
	Location             Location     `json:"location,omitempty"`
	Options              Options      `json:"options,omitempty"`
	Owner                Owner        `json:"owner,omitempty"`
	Price                []int        `json:"price,omitempty"`
	PriceCents           int          `json:"price_cents,omitempty"`
	PriceCalendar        interface{}  `json:"price_calendar,omitempty"`
	Status               string       `json:"status,omitempty"`
	Subject              string       `json:"subject,omitempty"`
	URL                  string       `json:"url,omitempty"`
}

func (a Ad) GetPrice() string {
	return fmt.Sprintf("%d€", a.PriceCents/100)
}

type Attributes struct {
	Key        string   `json:"key,omitempty"`
	Value      string   `json:"value,omitempty"`
	Values     []string `json:"values,omitempty"`
	ValueLabel string   `json:"value_label,omitempty"`
	Generic    bool     `json:"generic,omitempty"`
	KeyLabel   string   `json:"key_label,omitempty"`
}

type Images struct {
	ThumbURL  string   `json:"thumb_url,omitempty"`
	SmallURL  string   `json:"small_url,omitempty"`
	NbImages  int      `json:"nb_images,omitempty"`
	Urls      []string `json:"urls,omitempty"`
	UrlsThumb []string `json:"urls_thumb,omitempty"`
	UrlsLarge []string `json:"urls_large,omitempty"`
}

type Geometry struct {
	Type        string    `json:"type,omitempty"`
	Coordinates []float64 `json:"coordinates,omitempty"`
}

type Feature struct {
	Type       string      `json:"type,omitempty"`
	Geometry   Geometry    `json:"geometry,omitempty"`
	Properties interface{} `json:"properties,omitempty"`
}

type Location struct {
	CountryID      string  `json:"country_id,omitempty"`
	RegionID       string  `json:"region_id,omitempty"`
	RegionName     string  `json:"region_name,omitempty"`
	DepartmentID   string  `json:"department_id,omitempty"`
	DepartmentName string  `json:"department_name,omitempty"`
	CityLabel      string  `json:"city_label,omitempty"`
	City           string  `json:"city,omitempty"`
	Zipcode        string  `json:"zipcode,omitempty"`
	Lat            float64 `json:"lat,omitempty"`
	Lng            float64 `json:"lng,omitempty"`
	Source         string  `json:"source,omitempty"`
	Provider       string  `json:"provider,omitempty"`
	IsShape        bool    `json:"is_shape,omitempty"`
	Feature        Feature `json:"feature,omitempty"`
}

type Options struct {
	HasOption  bool `json:"has_option,omitempty"`
	Booster    bool `json:"booster,omitempty"`
	Photosup   bool `json:"photosup,omitempty"`
	Urgent     bool `json:"urgent,omitempty"`
	Gallery    bool `json:"gallery,omitempty"`
	SubToplist bool `json:"sub_toplist,omitempty"`
}

type Owner struct {
	StoreID        string `json:"store_id,omitempty"`
	UserID         string `json:"user_id,omitempty"`
	Type           string `json:"type,omitempty"`
	Name           string `json:"name,omitempty"`
	Siren          string `json:"siren,omitempty"`
	NoSalesmen     bool   `json:"no_salesmen,omitempty"`
	ActivitySector string `json:"activity_sector,omitempty"`
}
