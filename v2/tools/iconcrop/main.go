package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"code.olapie.com/sugar/v2/must"
	"github.com/disintegration/imaging"
)

func main() {
	pSize := flag.String("size", "20,29,40,60,76,167", "specify image size list")
	flag.Parse()

	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s [filename]", os.Args[0])
	}

	name := os.Args[1]
	f, err := os.Open(name)
	if err != nil {
		log.Fatalf("Cannot open file: %v", err)
	}
	defer f.Close()
	img, t, err := image.Decode(f)
	if err != nil {
		log.Fatalf("Cannot decode image: %v", err)
	}
	sizes := must.ToIntSlice(strings.Split(strings.TrimSpace(*pSize), ","))
	fmt.Println(sizes)
	name = name[:len(name)-len(path.Ext(name))]
	for _, size := range sizes {
		if err = crop(img, fmt.Sprintf("%s_%d.%s", name, size, t), t, size); err != nil {
			fmt.Printf("Error: size=%d, %v\n", size, err)
		}
		if err = crop(img, fmt.Sprintf("%s_%d@2x.%s", name, size, t), t, 2*size); err != nil {
			fmt.Printf("Error: size=%d, %v\n", 2*size, err)
		}
		if err = crop(img, fmt.Sprintf("%s_%d@3x.%s", name, size, t), t, 3*size); err != nil {
			fmt.Printf("Error: size=%d, %v\n", 3*size, err)
		}
	}
}

func crop(img image.Image, name, t string, size int) error {
	dx := img.Bounds().Dx()
	dy := img.Bounds().Dy()

	var tImg image.Image
	if dx < size || dy < size {
		return fmt.Errorf("size must <= %d and %d", dx, dy)
	}
	if dx > dy {
		tImg = imaging.Resize(img, 0, size, imaging.Lanczos)
	} else {
		tImg = imaging.Resize(img, size, 0, imaging.Lanczos)
	}
	buf := bytes.NewBuffer(nil)
	if t == "png" {
		if err := png.Encode(buf, tImg); err != nil {
			return fmt.Errorf("encode image to png: %w", err)
		}
	} else {
		if err := jpeg.Encode(buf, tImg, &jpeg.Options{
			Quality: 100,
		}); err != nil {
			return fmt.Errorf("encode image to jpeg: %w", err)
		}
	}
	return ioutil.WriteFile(name, buf.Bytes(), 0644)
}
