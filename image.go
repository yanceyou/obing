package bing

import (
	"net/url"
	"strings"
)

const hpImageNameSeparator = "_"

// HPImage image details
type HPImage struct {
	URL     string `json:"url"`
	URLBase string `json:"urlbase"`

	Market    string
	Title     string `json:"title"`
	Copyright string `json:"copyright"`

	FullStartDate string `json:"fullstartdate"`
	StartDate     string `json:"startdate"`
	EndDate       string `json:"enddate"`
}

// Filename return the image filename extract from url.
// url like: "/th?id=OHR.RootBridge_ZH-CN5173953292_1920x1080.jpg&rf=LaDigue_1920x1080.jpg&pid=hp"
func (i *HPImage) Filename() string {
	if len(i.URL) == 0 {
		return ""
	}
	values, err := url.ParseQuery(i.URL)
	if err != nil {
		return ""
	}
	if len(values["/th?id"]) == 0 {
		return ""
	}
	return values["/th?id"][0]
}

// MarketID return the market info (like ZH-CN5173953292) of the image filename.
// filename like: OHR.RootBridge_ZH-CN5173953292_1920x1080.jpg
func (i *HPImage) MarketID() string {
	items := strings.Split(i.Filename(), hpImageNameSeparator)
	if len(items) < 2 {
		return ""
	}
	return items[1]
}

// Name return the name (like OHR.RootBridge) of image filename.
// filename like: OHR.RootBridge_ZH-CN5173953292_1920x1080.jpg
func (i *HPImage) Name() string {
	items := strings.Split(i.Filename(), hpImageNameSeparator)
	if len(items) < 2 {
		return ""
	}
	return items[0]
}

// Content return the copyright or the title of the image.
func (i *HPImage) Content() string {
	if len(i.Title) == 0 {
		return i.Copyright
	}
	return i.Title
}
