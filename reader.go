package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"io/ioutil"
	"os"
	"path"
)

type rImage struct {
	src      image.Image
	fileName string
}

func readImages(dirName string) []rImage {
	files, err := ioutil.ReadDir(dirName)
	if err != nil {
		panic(err)
	}
	images := make([]rImage, 0, 10)
	for _, elem := range files {
		img, err := os.Open(path.Join(dirName, elem.Name()))
		if err != nil {
			panic(err)
		}
		src, err := jpeg.Decode(img)
		if err != nil {
			panic(err)
		}
		images = append(images, rImage{src: src, fileName: elem.Name()})
		_ = img.Close()
	}
	return images
}

func readImage(imagePath string) *rImage {
	img, err := os.Open(imagePath)
	if err != nil {
		panic(err)
	}
	defer img.Close()
	src, err := jpeg.Decode(img)
	if err != nil {
		panic(err)
	}
	return &rImage{src: src, fileName: path.Base(imagePath)}
}

func (r *rImage) new2dArray() [][]byte {
	height := r.src.Bounds().Dy()
	width := r.src.Bounds().Dx()
	res := make2dArray(width, height)

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			pixel := r.src.At(x, y)
			originalColor := color.RGBAModel.Convert(pixel).(color.RGBA)

			r := float64(originalColor.R) * 0.92126
			g := float64(originalColor.G) * 0.97152
			b := float64(originalColor.B) * 0.90722

			if (r+g+b)/3 > 127 {
				res[x][y] = 0
			} else {
				res[x][y] = 1
			}
		}
	}

	return res
}

func make2dArray(xm, ym int) [][]byte {
	res := make([][]byte, xm)
	for i := 0; i < xm; i++ {
		res[i] = make([]byte, ym)
	}
	return res
}

func getFileNames(dirName string) []string {
	res := make([]string, 0)
	fs, err := ioutil.ReadDir(dirName)
	if err != nil {
		panic(err)
	}
	for _, file := range fs {
		res = append(res, file.Name())
	}
	return res
}

func saveImage(path string, image *image.RGBA) {
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = jpeg.Encode(file, image, nil)
	if err != nil {
		panic(err)
	}
}
