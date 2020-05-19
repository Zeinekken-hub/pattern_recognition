package main

import (
	"fmt"
	"image/jpeg"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"
)

//RecognitionOutput ll
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
	http.HandleFunc("/", drawHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	http.ListenAndServe(":8080", nil)
}

func drawHandler(w http.ResponseWriter, r *http.Request) {
	prev, next := r.FormValue("p"), r.FormValue("n")
	getPathToFile(prev, next)
	recOut, err := scanAndSaveImage(currPath)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	tmp := template.Must(template.ParseFiles("index.html"))
	tmp.Execute(w, *recOut)
}

func scanAndSaveImage(imagePath string) (*RecognitionOutput, error) {
	file, err := os.Open(currPath)
	if err != nil {
		return nil, fmt.Errorf("error for reading file")
	}
	defer file.Close()

	start := time.Now()
	img, err := jpeg.Decode(file)
	slc, err := getByteImage(img)
	if err != nil {
		return nil, fmt.Errorf("get binary slice error")
	}
	rectangles := scan(slc.arr)
	slc.setCenters(rectangles)
	for _, rect := range rectangles {
		slc.SetRectangle(rect.start.x, rect.start.y, rect.end.x, rect.end.y, 2)
	}
	imgRgba := slc.convertOwnBytesToImage()
	elapsed := time.Since(start)

	fName := strings.Split(currPath, "/")[1]
	fPath := "static/data/" + fName
	file, err = os.Create(fPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	jpeg.Encode(file, imgRgba, nil)

	return &RecognitionOutput{
		ExecTime:              elapsed,
		FigureFound:           len(rectangles),
		PathToProccessedImage: fPath,
		FileName:              fName,
	}, nil
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
