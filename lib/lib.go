package lib

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"math/rand"
	"os"
	"strconv"
)

func RandomPortion(soldP []uint8, units string) ([]uint8, error) {
	sold := make(map[uint8]struct{})
	for _, p := range soldP {
		sold[p] = struct{}{}
	}
	all := make(map[uint8]struct{})
	for i := 0; i < 100; i++ {
		all[uint8(i)] = struct{}{}
	}
	var lotP []uint8
	for k := range all {
		_, ok := sold[k]
		if !ok {
			lotP = append(lotP, k)
		}
	}

	u, err := strconv.Atoi(units)
	if err != nil {
		return nil, fmt.Errorf("RandomPortion: %v", err)
	}

	// limit lottery units
	rem := 100 - len(soldP)
	var lim int
	if u < rem {
		lim = u
	} else {
		lim = rem
	}

	var randP []uint8
	for i := 0; i < lim; i++ {
		r := rand.Intn(len(lotP))
		randP = append(randP, lotP[r])
	}
	return randP, nil
}

func ProcessImage(fn string) error {
	path := fmt.Sprintf("upload/%s", fn)
	src, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("ProcessImage_Open: %v", err)
	}
	defer src.Close()

	srcImg, _, err := image.Decode(src)
	if err != nil {
		return fmt.Errorf("ProcessImage_Decode: %v", err)
	}
	srcBounds := srcImg.Bounds()

	dest := image.NewRGBA(srcBounds)

	for v := srcBounds.Min.Y; v < srcBounds.Max.Y; v++ {
		for h := srcBounds.Min.X; h < srcBounds.Max.X; h++ {
			curlPixel := srcImg.At(h, v)
			r, g, b, a := curlPixel.RGBA()
			r, g, b, a = r>>8, g>>8, b>>8, a>>8
			mean := (r + g + b) / 3
			col := color.RGBA{R: uint8(mean), G: uint8(mean), B: uint8(mean), A: uint8(a)}
			dest.Set(h, v, col)
		}
	}

	outPath := fmt.Sprintf("out/%s", fn)
	outfile, _ := os.Create(outPath)
	defer outfile.Close()

	png.Encode(outfile, dest)

	return nil
}
