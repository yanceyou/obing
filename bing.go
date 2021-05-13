package obing

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// bing.com hosts
const (
	HostDefault = "https://www.bing.com"
	HostGlobal  = "https://global.bing.com"
	HostCN      = "https://cn.bing.com"
)

// Microsoft Market Codes
// https://docs.microsoft.com/en-us/rest/api/cognitiveservices/bing-web-api-v7-reference#market-codes
// "es-AR", "en-AU", "de-AT", "nl-BE", "fr-BE", "pt-BR", "en-CA", "fr-CA",
// "es-CL", "da-DK", "fi-FI", "fr-FR", "de-DE", "zh-HK", "en-IN", "en-ID",
// "it-IT", "ja-JP", "ko-KR", "en-MY", "es-MX", "nl-NL", "en-NZ", "no-NO",
// "zh-CN", "pl-PL", "en-PH", "ru-RU", "en-ZA", "es-ES", "sv-SE", "fr-CH",
// "de-CH", "zh-TW", "tr-TR", "en-GB", "en-US", "es-US",

// RowMarketCodes has no special market HP images.
// HPImage URL which has `ROW` keywords
var RowMarketCodes = []string{
	"es-AR", "en-AU", "de-AT", "nl-BE", "fr-BE", "pt-BR", "es-CL", "da-DK", "fi-FI",
	"zh-HK", "en-ID", "it-IT", "ko-KR", "en-MY", "es-MX", "nl-NL", "en-NZ", "no-NO",
	"pl-PL", "en-PH", "ru-RU", "en-ZA", "es-ES", "sv-SE", "fr-CH", "de-CH", "zh-TW",
	"tr-TR", "es-US",
}

// MarketCodes only has valid market codes, the others are row market codes.
var MarketCodes = []string{
	"en-CA", "fr-CA", "fr-FR", "de-DE", "en-IN", "ja-JP", "zh-CN", "en-GB", "en-US",
}

func getHPImages(host, market string, index, num int) (images []*HPImage, err error) {
	url := fmt.Sprintf("%s/HPImageArchive.aspx?format=js&setmkt=%s&idx=%d&n=%d", host, market, index, num)
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
	for _, img := range res.Images {
		img.Market = market
		img.Host = host
	}
	return res.Images, err
}

// GetMarketHPImages get HP images from global host and target market
// `index` means days before today, and -1 <= index <= 7
// `num` means images number before the `index` day, and num <= 7
// so, we can get images of 7 + 7 days ago, once we set `index = 7` and `num = 7`
func GetMarketHPImages(market string, index, num int) (images []*HPImage, err error) {
	return getHPImages(HostGlobal, market, index, num)
}

// GetAllMarketHPImages get HP images from all market codes.
func GetAllMarketHPImages(index, num int) (images []*HPImage, err error) {
	imgs := make([]*HPImage, 0)
	for _, mkt := range MarketCodes {
		mktImgs, err := GetMarketHPImages(mkt, index, num)
		if err != nil {
			return nil, err
		}
		imgs = append(imgs, mktImgs...)
	}
	return imgs, nil
}

// GetHPImages get HP images with index and num
func GetHPImages(index, num int) (images []*HPImage, err error) {
	return getHPImages(HostCN, "", index, num)
}
