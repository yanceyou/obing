package bing

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

// bing.com hosts
const (
	HostGlobal = "https://global.bing.com"
	HostCN     = "https://cn.bing.com"
)

// MarketCodes are Microsoft Market Codes
// https://docs.microsoft.com/en-us/rest/api/cognitiveservices/bing-web-api-v7-reference#market-codes
var MarketCodes = []string{
	"es-AR", "en-AU", "de-AT", "nl-BE", "fr-BE", "pt-BR", "en-CA", "fr-CA",
	"es-CL", "da-DK", "fi-FI", "fr-FR", "de-DE", "zh-HK", "en-IN", "en-ID",
	"it-IT", "ja-JP", "ko-KR", "en-MY", "es-MX", "nl-NL", "en-NZ", "no-NO",
	"zh-CN", "pl-PL", "en-PH", "ru-RU", "en-ZA", "es-ES", "sv-SE", "fr-CH",
	"de-CH", "zh-TW", "tr-TR", "en-GB", "en-US", "es-US",
}

// HPImage image details
type HPImage struct {
	Copyright string `json:"copyright"`
	URL       string `json:"url"`
	URLBase   string `json:"urlbase"`

	FullStartDate string `json:"fullstartdate"`
	StartDate     string `json:"startdate"`
	EndDate       string `json:"enddate"`
}

// ToString HP Image info with date, copyright and name
func (i *HPImage) ToString() string {
	return fmt.Sprintf("Startdate: %s | Copyright: %s | Filename: %s\n", i.FullStartDate, i.Copyright, filepath.Base(i.URL))
}

// GetMarketHPImages get HP images from target host and market
// `index` means days before today, and -1 <= index <= 7
// `n` means images number before the `index` day, and n <= 7
// so, we can get images of 7 + 7 days ago, once we set `index = 7` and `n = 7`
func GetMarketHPImages(host, market string, index, n int) (images []*HPImage, err error) {
	url := fmt.Sprintf("%s/HPImageArchive.aspx?format=js&setmkt=%s&idx=%d&n=%d", host, market, index, n)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var res struct {
		Images []*HPImage `json:"images"`
	}
	err = json.Unmarshal(body, &res)
	return res.Images, err
}

// DownloadHPImage download HP image from url to destination path
// return final image file path
func DownloadHPImage(url, dest string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	filename := filepath.Join(dest, path.Base(url))
	file, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return filename, err
}
