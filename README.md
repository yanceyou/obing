# obing

Download HP images from https://www.bing.com

## Example

```golang
package main

import (
	"log"

	"github.com/yanceyou/bing"
)

func main() {
	images, err := bing.GetMarketHPImages(bing.HostGlobal, "zh-CN", 0, 7)
	if err != nil {
		log.Println(err)
	}
	for i, image := range images {
		log.Printf("image-[%d]: %+v", i, image)
		filename, err := bing.DownloadHPImage(bing.HostGlobal+image.URL, ".")
		if err != nil {
			log.Printf("download image-[%d] err: %+v", i, err)
		}
		log.Printf("download image-[%d] into %s", i, filename)
	}
}
```
