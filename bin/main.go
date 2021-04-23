package main

import (
	"fmt"
	"strings"

	"github.com/yanceyou/obing"
)

func main() {
	images, err := obing.GetHPImages(1, 7)
	fmt.Println(err)
	for index, image := range images {
		if !strings.Contains(image.URL, "_ROW") {
			fmt.Printf("%d: %+v, download: %+v\n", index, image, image.Download("/home/yanceyyang/2021/obing/bin"))
		}
	}
}
