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

func ProcessImage(fn string, getP map[uint8]struct{}) error {
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

	height, width := srcBounds.Max.Y, srcBounds.Max.X
	hUnit := height / 10
	wUnit := width / 10

	for p := 0; p < 100; p++ {
		_, ok := getP[uint8(p)]
		if ok {
			for w := (p % 10) * wUnit; w < ((p%10)+1)*wUnit; w++ {
				for h := (p / 10) * hUnit; h < ((p/10)+1)*hUnit; h++ {
					pixel := srcImg.At(w, h)
					r, g, b, a := pixel.RGBA()
					r, g, b, a = r>>8, g>>8, b>>8, a>>8
					mean := (r + g + b) / 3
					a = a / 5 * 4
					col := color.RGBA{R: uint8(mean), G: uint8(mean), B: uint8(mean), A: uint8(a)}
					dest.Set(w, h, col)
				}
			}
		} else {
			for w := (p % 10) * wUnit; w < ((p%10)+1)*wUnit; w++ {
				for h := (p / 10) * hUnit; h < ((p/10)+1)*hUnit; h++ {
					pixel := srcImg.At(w, h)
					r, g, b, a := pixel.RGBA()
					r, g, b, a = r>>8, g>>8, b>>8, a>>8
					col := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
					dest.Set(w, h, col)
				}
			}
		}
	}

	outPath := fmt.Sprintf("out/%s", fn)
	outfile, _ := os.Create(outPath)
	defer outfile.Close()

	png.Encode(outfile, dest)

	return nil
}
