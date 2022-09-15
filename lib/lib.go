package lib

import (
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"

	"crypto/sha512"
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
		lotP = append(lotP[:r], lotP[r+1:]...)
	}
	return randP, nil
}

func ProcessImage(fn string, getP map[uint8]struct{}) error {
	path := fmt.Sprintf("uploaded/%s", fn)
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

	// fill Portions
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

func SplitImage(fn string) error {
	path := filepath.Join("uploaded", fn)
	src, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("SplitImage_Open: %v", err)
	}
	defer src.Close()

	srcImg, _, err := image.Decode(src)
	if err != nil {
		return fmt.Errorf("SplitImage_Decode: %v", err)
	}
	srcBounds := srcImg.Bounds()

	height, width := srcBounds.Max.Y, srcBounds.Max.X
	hUnit := int(math.Ceil(float64(height) / 10))
	wUnit := int(math.Ceil(float64(width) / 10))
	var hEdge, wEdge int

	if hUnit/2 >= (height - hUnit*9) {
		hUnit -= 1
	}
	if wUnit/2 >= (width - wUnit*9) {
		wUnit -= 1
	}
	wEdge = width - wUnit*9
	hEdge = height - hUnit*9

	dirName := filepath.Base(fn[:len(fn)-len(filepath.Ext(fn))])
	dirPath := filepath.Join("out", "original", dirName)

	if err := os.RemoveAll(dirPath); err != nil {
		log.Fatal(err)
	}

	if err := os.Mkdir(dirPath, 0777); err != nil {
		log.Fatal(err)
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
		dst := image.NewRGBA(image.Rect(0, 0, wUnit, hUnit))
		draw.Draw(dst, srcBounds, srcImg, image.Pt((p%10)*wUnit, (p/10)*hUnit), draw.Src)

		fileName := fmt.Sprintf("No_%s.png", strconv.Itoa(p))
		outPath := filepath.Join(dirPath, fileName)
		outfile, err := os.Create(outPath)
		if err != nil {
			log.Fatal(err)
		}

		png.Encode(outfile, dst)
		outfile.Close()
	}

	// fill right side
	for _, p := range []int{9, 19, 29, 39, 49, 59, 69, 79, 89} {
		dst := image.NewRGBA(image.Rect(0, 0, wEdge, hUnit))
		draw.Draw(dst, srcBounds, srcImg, image.Pt((p%10)*wUnit, (p/10)*hUnit), draw.Src)

		fileName := fmt.Sprintf("No_%s.png", strconv.Itoa(p))
		outPath := filepath.Join(dirPath, fileName)
		outfile, err := os.Create(outPath)
		if err != nil {
			log.Fatal(err)
		}

		png.Encode(outfile, dst)
		outfile.Close()
	}

	// fill bottom side
	for _, p := range []int{90, 91, 92, 93, 94, 95, 96, 97, 98} {
		dst := image.NewRGBA(image.Rect(0, 0, wUnit, hEdge))
		draw.Draw(dst, srcBounds, srcImg, image.Pt((p%10)*wUnit, (p/10)*hUnit), draw.Src)

		fileName := fmt.Sprintf("No_%s.png", strconv.Itoa(p))
		outPath := filepath.Join(dirPath, fileName)
		outfile, err := os.Create(outPath)
		if err != nil {
			log.Fatal(err)
		}

		png.Encode(outfile, dst)
		outfile.Close()
	}

	// fill corner
	dst := image.NewNRGBA(image.Rect(0, 0, wEdge, hEdge))
	draw.Draw(dst, srcBounds, srcImg, image.Pt(9*wUnit, 9*hUnit), draw.Src)

	fileName := fmt.Sprintf("No_%s.png", strconv.Itoa(99))
	outPath := filepath.Join(dirPath, fileName)
	outfile, err := os.Create(outPath)
	if err != nil {
		log.Fatal(err)
	}

	png.Encode(outfile, dst)
	outfile.Close()

	return nil
}

//エンコード
func Encode(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fi, err := file.Stat() //FileInfo interface
	if err != nil {
		return "", err
	}
	size := fi.Size() //ファイルサイズ

	data := make([]byte, size)
	file.Read(data)

	return base64.StdEncoding.EncodeToString(data), err
}

func GetWidth(path string) (int, error) {
	src, err := os.Open(path)
	if err != nil {
		return 0, fmt.Errorf("GetWidth_Open: %v", err)
	}
	defer src.Close()

	srcImg, _, err := image.Decode(src)
	if err != nil {
		return 0, fmt.Errorf("GetWidth_Decode: %v", err)
	}
	srcBounds := srcImg.Bounds()
	return srcBounds.Max.X, nil
}

func MakeHash(p string) string {
	pbyte := []byte(p)
	pHash := sha512.Sum512(pbyte)
	return fmt.Sprintf("%x", pHash)
}
