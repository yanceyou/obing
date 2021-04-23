package main

import (
	"fmt"

	"github.com/yanceyou/obing"
)

func main() {
	images, err := obing.GetHPImages(1, 7)
	fmt.Println(err)
	for index, image := range images {
		fmt.Printf("%d: %+v\n", index, image)
		// fmt.Printf("%d-download: %+v", index, image.Download("/home/yanceyyang/2021/obing/bin"))
		fmt.Printf("%d-download: %+v", index, image.DownloadResolution("/home/yanceyyang/2021/obing/bin", 580, 800))
	}
}
