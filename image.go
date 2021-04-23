package obing

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// HPImage image details
type HPImage struct {
	Market string
	Host   string

	URL     string `json:"url"`
	URLBase string `json:"urlbase"`
	Hash    string `json:"hsh"`

	Title string `json:"title"`
	Quiz  string `json:"quiz"`

	Copyright     string `json:"copyright"`
	CopyrightLink string `json:"copyrightlink"`

	FullStartDate string `json:"fullstartdate"`
	StartDate     string `json:"startdate"`
	EndDate       string `json:"enddate"`
}

const (
	R1920x1080 = "1920x1080"
	R360x480 = "360x480"
)

// Filename return the image filename extract from url.
// url like: "/th?id=OHR.RootBridge_ZH-CN5173953292_1920x1080.jpg&rf=LaDigue_1920x1080.jpg&pid=hp"
// filename: "OHR.RootBridge_ZH-CN5173953292_1920x1080.jpg"
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

// Name return the name (like OHR.RootBridge) of image filename.
// filename: "OHR.RootBridge_ZH-CN5173953292_1920x1080.jpg"
// name: "OHR.RootBridge"
func (i *HPImage) Name() string {
	items := strings.Split(i.Filename(), "_")
	if len(items) < 2 {
		return ""
	}
	return items[0]
}

func (i *HPImage) download(url, target string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(target)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

// Download HP image to destination folder.
func (i *HPImage) Download(folder string) error {
	return i.download(i.Host+i.URL, filepath.Join(folder, i.Filename()))
}

// DownloadResolution HP image with resolution to destination folder.
func (i *HPImage) DownloadResolution(folder string, w, h int) error {
	resolution := fmt.Sprintf("%dx%d", w, h)
	url := strings.Replace(i.Host+i.URL, "1920x1080", resolution, 1)
	filename := strings.Replace(i.Filename(), "1920x1080", resolution, 1)
	return i.download(url, filepath.Join(folder, filename))
}
