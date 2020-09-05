package main

import (
	"image"
	"log"
	"net/http"
	"path"
	"text/template"
	"time"
)

var (
	fileNames       []string
	fileNum         int
	allImages       []rImage
	dataFolder      = "data"
	processedFolder = "static/output"
	tmp             *template.Template
)

//RecognitionOutput template
type RecognitionOutput struct {
	ExecTime              time.Duration
	FigureFound           int
	PathToProccessedImage string
	FileName              string
	MaxSquare             int
	MinSquare             int
	Figures               []Figure
}

//Figure template
type Figure struct {
	Center point
	Square int
}

func main() {
	http.HandleFunc("/rec", drawHandler)
	allImages = readImages(dataFolder)
	fileNames = getFileNames(dataFolder)
	tmp = template.Must(template.ParseFiles("index.html"))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	log.Printf("Server started on localhost:8080\n")
	log.Printf("Go to localhost:8080/rec to get started\n")
	http.ListenAndServe(":8080", nil)
}

func drawHandler(w http.ResponseWriter, r *http.Request) {
	p := r.FormValue("p")
	if p != "" {
		if fileNum-1 < 0 {
			fileNum = len(fileNames) - 1
		} else {
			fileNum--
		}
	} else {
		if fileNum+1 == len(fileNames) {
			fileNum = 0
		} else {
			fileNum++
		}
	}
	log.Printf("file name: %s\n", fileNames[fileNum])
	img := allImages[fileNum]
	recOut := &RecognitionOutput{}
	rgba := scanImage(allImages[fileNum], recOut)
	saveImage(path.Join(processedFolder, img.fileName), rgba)

	_ = tmp.Execute(w, recOut)
}

func scanImage(img rImage, rec *RecognitionOutput) *image.RGBA {
	start := time.Now()
	b := bImage{src: img.new2dArray()}
	b.setBounds()
	b.scan()
	log.Print("scan done")
	b.SetCenters()
	log.Print("set centers done")
	b.SetRectangles()
	rec.ExecTime = time.Since(start)
	rec.PathToProccessedImage = path.Join(processedFolder, img.fileName)
	rec.FigureFound = len(b.cache)
	rec.FileName = img.fileName
	rec.MaxSquare, rec.MinSquare = b.maxMinFigureSquare()
	log.Print("max min figure square done")
	rec.Figures = make([]Figure, 0)
	for _, elem := range b.cache {
		rec.Figures = append(rec.Figures, Figure{Center: elem.center(), Square: b.getFigureSquare(elem)})
	}
	return b.NewRGBAImage()
}
