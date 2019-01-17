package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/yanceyou/bing"
)

const folder = "/mnt/f/images"
const csvFilename = "bing.csv"
const logFilename = "bing.log"

func init() {
	logFilepath := filepath.Join(folder, logFilename)
	f, err := os.OpenFile(logFilepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)
}

func toCSVMapper() map[string]int {
	csvMapper := make(map[string]int, 0)
	types := reflect.TypeOf(bing.HPImage{})
	for i := 0; i < types.NumField(); i++ {
		csvMapper[types.Field(i).Name] = i
	}
	return csvMapper
}

func readCSV(csvFile *os.File) ([]*bing.HPImage, error) {
	imgs := make([]*bing.HPImage, 0)

	reader := csv.NewReader(csvFile)
	reader.Comma = ';'
	reader.LazyQuotes = true
	lines, err := reader.ReadAll()
	if err != nil {
		return imgs, err
	}

	csvMapper := toCSVMapper()
	for _, line := range lines {
		img := &bing.HPImage{}
		for k, v := range csvMapper {
			reflect.ValueOf(img).Elem().FieldByName(k).SetString(line[v])
		}
		imgs = append(imgs, img)
	}
	return imgs, nil
}

func getAllMarketHPImages() ([]*bing.HPImage, error) {
	imgs := make([]*bing.HPImage, 0)
	for _, mkt := range bing.MarketCodes {
		mktImages, err := bing.GetMarketHPImages(bing.HostGlobal, mkt, 0, 1)
		if err != nil {
			return nil, err
		}
		imgs = append(imgs, mktImages...)
	}
	return imgs, nil
}

func getHPImageName(URLBase string) string {
	const HPImageNameSeparator = "_"
	base := path.Base(URLBase)
	names := strings.Split(base, HPImageNameSeparator)
	if len(names) == 0 {
		return fmt.Sprintf("unknown-%d", time.Now().UnixNano())
	}
	return names[0]
}

func deduplication(newImgs []*bing.HPImage, oldImgs []*bing.HPImage) []*bing.HPImage {
	dedupedImgs := make([]*bing.HPImage, 0)
	for _, nImg := range newImgs {
		nImgName := getHPImageName(nImg.URLBase)
		has := false
		for _, oImg := range oldImgs {
			if nImgName == getHPImageName(oImg.URLBase) {
				has = true
			}
		}
		for _, dImg := range dedupedImgs {
			if nImgName == getHPImageName(dImg.URLBase) {
				has = true
			}
		}
		if !has {
			dedupedImgs = append(dedupedImgs, nImg)
		}
	}
	return dedupedImgs
}

func toCSVLine(image *bing.HPImage) []string {
	values := reflect.ValueOf(image).Elem()
	line := make([]string, values.NumField())
	for k, v := range toCSVMapper() {
		val, ok := values.FieldByName(k).Interface().(string)
		if !ok {
			line[v] = ""
			continue
		}
		line[v] = val
	}
	return line
}

func main() {
	log.Println("Start get today's backgrounds...")

	csvFilepath := filepath.Join(folder, csvFilename)
	csvFile, err := os.OpenFile(csvFilepath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		log.Printf("[ERROR] Open csv file err: %+v", err)
		return
	}
	defer csvFile.Close()

	oldImgs, err := readCSV(csvFile)
	if err != nil {
		log.Printf("[ERROR] Read csv file err: %+v", err)
		return
	}
	if len(oldImgs) == 0 {
		// Write UTF-8 BOM, 防止乱码
		csvFile.WriteString("\xEF\xBB\xBF")
	}

	newImgs, err := getAllMarketHPImages()
	if err != nil {
		log.Printf("[ERROR] Get market HP images err: %+v", err)
		return
	}

	writer := csv.NewWriter(csvFile)
	writer.Comma = ';'
	dedupedImgs := deduplication(newImgs, oldImgs)
	for i, img := range dedupedImgs {
		filename, err := bing.DownloadHPImage(bing.HostGlobal+img.URL, folder)
		if err != nil {
			log.Printf("[ERROR] Download [image-%d] err: %+v", i, err)
			continue
		}
		log.Printf("Download [image-%d] into [%s]: %+v", i, filename, img)

		if err := writer.Write(toCSVLine(img)); err != nil {
			log.Printf("[ERROR] Writing img to csv err: %+v", err)
		}
	}
	writer.Flush()

	if err := writer.Error(); err != nil {
		log.Printf("[ERROR] Flush csv err: %+v", err)
	}
	log.Println("End download today's backgrounds...")
}
