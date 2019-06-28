package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/yanceyou/bing"
)

var (
	dataFile *os.File
)

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
	days := flag.Int("days", 0, "days before today")
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
	logFilepath := filepath.Join(conf.folder, "bing.log")
	logFile, err := os.OpenFile(logFilepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("Open log file err: %v", err))
	}
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
}

func init() {
	initConfig()
	initLogger()
}

func loadData() {

	year := time.Now().Year()
	flag := os.O_RDWR | os.O_CREATE | os.O_APPEND

	dataFilename := filepath.Join(conf.folder, fmt.Sprintf("bing.%d.json", year))
	dataFile, err := os.OpenFile(dataFilename, flag, 0666)
	if err != nil {
		panic(fmt.Sprintf("Open csv file err: %+v", err))
	}

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

func start() error {
	dataBytes, err := ioutil.ReadAll(dataFile)
	if err != nil {
		return fmt.Errorf("Read data file err: %+v", err)
	}
	var oldImgs []*bing.HPImage
	if err := json.Unmarshal(dataBytes, &oldImgs); err != nil {
		return fmt.Errorf("Unmarshal data file err: %+v", err)
	}

	newImgs, err := bing.GetAllMarketHPImages(conf.host, conf.days, conf.num)
	if err != nil {
		return fmt.Errorf("Get market HP images err: %+v", err)
	}

	dedupedNewImgs := deduplication(newImgs, oldImgs)
	log.Printf("Get deduped market HP images num: %d", len(dedupedNewImgs))

	for i, img := range dedupedNewImgs {
		filename := filepath.Join(conf.folder, img.Name())
		if err := bing.DownloadHPImage(conf.host+img.URL, filename); err != nil {
			log.Printf("[ERROR] Download [image-%d] err: %+v", i, err)
			continue
		}
		log.Printf("Download [image-%d] into [%s]: %+v", i, filename, img)
	}

	oldImgs = append(oldImgs, dedupedNewImgs...)

}

func main() {
	log.Println("Start get today's backgrounds...")
	log.Printf("Current config: %+v", conf)

	if err := start(); err != nil {
		log.Printf("[ERROR] %+v", err)
	}

	log.Println("End download today's backgrounds...")
}
