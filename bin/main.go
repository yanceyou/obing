package main

import (
	"fmt"

	"github.com/yanceyou/obing"
)

func main() {
	images, err := obing.GetHPImages(1, 7)
	if err != nil {
		fmt.Println(err)
		return
	}
	for i, image := range images {
		err := image.Download(".")
		if err != nil {
			fmt.Printf("download image-[%d] err: %+v\n", i, err)
		}
		fmt.Printf("download image-[%d] image success: %+v\n", i, image)
	}
}
