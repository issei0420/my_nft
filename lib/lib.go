package lib

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"math"
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
	hUnit := int(math.Ceil(float64(height) / 10))
	wUnit := int(math.Ceil(float64(width) / 10))

	if hUnit/2 >= (height - hUnit*9) {
		hUnit -= 1
	}
	if wUnit/2 >= (width - wUnit*9) {
		wUnit -= 1
	}

	pEdge := map[int]struct{}{
		9: {}, 19: {}, 29: {}, 39: {}, 49: {}, 59: {}, 69: {}, 79: {}, 89: {},
		90: {}, 91: {}, 92: {}, 93: {}, 94: {}, 95: {}, 96: {}, 97: {}, 98: {}, 99: {},
	}

	var portions []int
	for i := 0; i < 100; i++ {
		if _, ok := pEdge[i]; !ok {
			portions = append(portions, i)
		}
	}

	// fill portions
	for _, p := range portions {
		_, ok := getP[uint8(p)]
		if !ok {
			for w := (p % 10) * wUnit; w < (p%10+1)*wUnit; w++ {
				for h := (p / 10) * hUnit; h < (p/10+1)*hUnit; h++ {
					// c := color.GrayModel.Convert(srcImg.At(w, h))
					// gray, _ := c.(color.Gray)
					pixel := srcImg.At(w, h)
					r, g, b, _ := pixel.RGBA()
					r, g, b = r>>8, g>>8, b>>8
					mean := (r + g + b) / 3
					col := color.RGBA{R: uint8(mean), G: uint8(mean), B: uint8(mean), A: 100}
					dest.Set(w, h, col)
				}
			}
		} else {
			for w := (p % 10) * wUnit; w < (p%10+1)*wUnit; w++ {
				for h := (p / 10) * hUnit; h < (p/10+1)*hUnit; h++ {
					pixel := srcImg.At(w, h)
					r, g, b, a := pixel.RGBA()
					r, g, b, a = r>>8, g>>8, b>>8, a>>8
					col := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
					dest.Set(w, h, col)
				}
			}
		}
	}

	// fill right side
	for _, p := range []int{9, 19, 29, 39, 49, 59, 69, 79, 89} {
		_, ok := getP[uint8(p)]
		if !ok {
			for w := (p % 10) * wUnit; w < width; w++ {
				for h := (p / 10) * hUnit; h < (p/10+1)*hUnit; h++ {
					pixel := srcImg.At(w, h)
					r, g, b, _ := pixel.RGBA()
					r, g, b = r>>8, g>>8, b>>8
					mean := (r + g + b) / 3
					col := color.RGBA{R: uint8(mean), G: uint8(mean), B: uint8(mean), A: 100}
					dest.Set(w, h, col)
				}
			}
		} else {
			for w := (p % 10) * wUnit; w < width; w++ {
				for h := (p / 10) * hUnit; h < (p/10+1)*hUnit; h++ {
					pixel := srcImg.At(w, h)
					r, g, b, a := pixel.RGBA()
					r, g, b, a = r>>8, g>>8, b>>8, a>>8
					col := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
					dest.Set(w, h, col)
				}
			}
		}
	}

	// fill bottom side
	for _, p := range []int{90, 91, 92, 93, 94, 95, 96, 97, 98} {
		_, ok := getP[uint8(p)]
		if !ok {
			for w := (p % 10) * wUnit; w < (p%10+1)*wUnit; w++ {
				for h := (p / 10) * hUnit; h < height; h++ {
					pixel := srcImg.At(w, h)
					r, g, b, _ := pixel.RGBA()
					r, g, b = r>>8, g>>8, b>>8
					mean := (r + g + b) / 3
					col := color.RGBA{R: uint8(mean), G: uint8(mean), B: uint8(mean), A: 100}
					dest.Set(w, h, col)
				}
			}
		} else {
			for w := (p % 10) * wUnit; w < (p%10+1)*wUnit; w++ {
				for h := (p / 10) * hUnit; h < height; h++ {
					pixel := srcImg.At(w, h)
					r, g, b, a := pixel.RGBA()
					r, g, b, a = r>>8, g>>8, b>>8, a>>8
					col := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
					dest.Set(w, h, col)
				}
			}
		}
	}

	// fill corner
	_, ok := getP[99]
	if !ok {
		for w := 9 * wUnit; w < width; w++ {
			for h := 9 * hUnit; h < height; h++ {
				pixel := srcImg.At(w, h)
				r, g, b, _ := pixel.RGBA()
				r, g, b = r>>8, g>>8, b>>8
				mean := (r + g + b) / 3
				col := color.RGBA{R: uint8(mean), G: uint8(mean), B: uint8(mean), A: 100}
				dest.Set(w, h, col)
			}
		}
	} else {
		for w := 9 * wUnit; w < width; w++ {
			for h := 9 * hUnit; h < height; h++ {
				pixel := srcImg.At(w, h)
				r, g, b, a := pixel.RGBA()
				r, g, b, a = r>>8, g>>8, b>>8, a>>8
				col := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
				dest.Set(w, h, col)
			}
		}

	}

	outPath := fmt.Sprintf("out/%s", fn)
	outfile, _ := os.Create(outPath)
	defer outfile.Close()

	png.Encode(outfile, dest)

	return nil
}
