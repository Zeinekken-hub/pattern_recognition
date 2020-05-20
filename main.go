package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"text/template"
	"time"
)

//RecognitionOutput template
type RecognitionOutput struct {
	ExecTime              time.Duration
	FigureFound           int
	PathToProccessedImage string
	FileName              string
}

var (
	currPath = "data/1_0001.jpg"
	fileNum  = 1
)

func main() {
	http.HandleFunc("/rec", drawHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	http.ListenAndServe(":8080", nil)
}

func drawHandler(w http.ResponseWriter, r *http.Request) {
	prev, next := r.FormValue("p"), r.FormValue("n")
	getPathToFile(prev, next)

	recOut := &RecognitionOutput{}
	err := scanImage(currPath, recOut)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	tmp := template.Must(template.ParseFiles("index.html"))
	tmp.Execute(w, *recOut)
}

func scanImage(imagePath string, recOut *RecognitionOutput) error {
	file, err := os.Open(currPath)
	if err != nil {
		return fmt.Errorf("error for reading file")
	}
	defer file.Close()

	start := time.Now()
	img, err := jpeg.Decode(file)
	if err != nil {
		return fmt.Errorf("error for decoding image")
	}
	byteImg, err := getByteImage(img)
	if err != nil {
		return fmt.Errorf("get binary slice error")
	}
	log.Printf("Image name: %s\n", path.Base(currPath))
	rectangles := scan(byteImg.arr)
	byteImg.SetCenters(rectangles)
	for _, rect := range rectangles {
		byteImg.SetRectangle(rect.start.x, rect.start.y, rect.end.x, rect.end.y)
	}
	imgRgba := byteImg.NewRGBAImage()
	elapsed := time.Since(start)
	log.Printf("Time of scan executing: %v\n", elapsed)

	recOut.ExecTime = elapsed
	recOut.FigureFound = len(rectangles)

	saveImage(imgRgba, recOut)

	return nil
}

func saveImage(image *image.RGBA, recOut *RecognitionOutput) {
	fName := strings.Split(currPath, "/")[1]
	fPath := "static/data/" + fName
	file, err := os.Create(fPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	jpeg.Encode(file, image, nil)

	recOut.PathToProccessedImage = fPath
	recOut.FileName = fName
}

func formatPathByInt(x int) string {
	s := "data/1_00"
	if x < 10 {
		s += "0" + strconv.Itoa(x)
	} else {
		s += strconv.Itoa(x)
	}
	s += ".jpg"
	return s
}

func getPathToFile(prevPost, nextPost string) {
	if prevPost != "" {
		fileNum--
		if fileNum == 0 {
			currPath = formatPathByInt(32)
			fileNum = 32
		} else {
			currPath = formatPathByInt(fileNum)
		}
	} else if nextPost != "" {
		fileNum++
		if fileNum == 33 {
			currPath = formatPathByInt(1)
			fileNum = 1
		} else {
			currPath = formatPathByInt(fileNum)
		}
	}
}
