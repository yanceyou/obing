package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/yanceyou/bing"
)

const csvComma = ';'
const csvFilename = "bing.dat"
const logFilename = "bing.log"

type config struct {
	folder string
	host   string
	days   int
	num    int
}

var conf config

func initConfig() {
	u, err := user.Current()
	if err != nil {
		panic(err)
	}
	folder := flag.String("folder", filepath.Join(u.HomeDir, "pictures"), "target folder")
	days := flag.Int("days", 7, "days before today")
	num := flag.Int("num", 7, "image numbers of days")

	flag.Parse()

	if err := os.MkdirAll(*folder, os.ModePerm); err != nil {
		panic(err)
	}

	conf = config{
		folder: *folder,
		days:   *days,
		num:    *num,
		host:   bing.HostGlobal,
	}
}

func initLogger() {
	logFilepath := filepath.Join(conf.folder, logFilename)
	f, err := os.OpenFile(logFilepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("error opening file: %v", err))
	}
	log.SetOutput(io.MultiWriter(os.Stdout, f))
}

func init() {
	initConfig()
	initLogger()
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
	reader.Comma = csvComma
	reader.LazyQuotes = true
	lines, err := reader.ReadAll()
	if err != nil {
		return imgs, err
	}

	csvMapper := toCSVMapper()
	for _, line := range lines {
		if len(line) < len(csvMapper) {
			continue
		}
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
		mktImages, err := bing.GetMarketHPImages(conf.host, mkt, conf.days, conf.num)
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

func start() error {
	csvFilepath := filepath.Join(conf.folder, csvFilename)
	csvFile, err := os.OpenFile(csvFilepath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return fmt.Errorf("Open csv file err: %+v", err)
	}
	defer csvFile.Close()

	oldImgs, err := readCSV(csvFile)
	if err != nil {
		return fmt.Errorf("Read csv file err: %+v", err)
	}
	if len(oldImgs) == 0 {
		csvFile.WriteString("\xEF\xBB\xBF")
	}

	newImgs, err := getAllMarketHPImages()
	if err != nil {
		return fmt.Errorf("Get market HP images err: %+v", err)
	}

	writer := csv.NewWriter(csvFile)
	writer.Comma = csvComma
	dedupedImgs := deduplication(newImgs, oldImgs)
	log.Printf("Get deduped market HP images num: %d", len(dedupedImgs))
	for i, img := range dedupedImgs {
		filename, err := bing.DownloadHPImage(conf.host+img.URL, conf.folder)
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
	return writer.Error()
}

func main() {
	log.Println("Start get today's backgrounds...")
	log.Printf("Current config: %+v", conf)

	if err := start(); err != nil {
		log.Printf("[ERROR] %+v", err)
	}

	log.Println("End download today's backgrounds...")
}
