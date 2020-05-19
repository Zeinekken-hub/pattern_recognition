package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"

	"golang.org/x/image/bmp"
)

func main() {
	http.HandleFunc("/draw", drawHandler)
	http.HandleFunc("/main", mainHandler)
	http.HandleFunc("/upload", loadHandler)

	staticHandler := http.StripPrefix("/data/", http.FileServer(http.Dir("./data")))
	http.Handle("/data/", staticHandler)

	fmt.Println("start listen server on port :8080")
	http.ListenAndServe(":8080", nil)
}

func drawHandler(w http.ResponseWriter, r *http.Request) {
	path := "data/1_0001.jpg"
	key := r.FormValue("imgP")
	if key != "" {
		path = key
	}

	file, err := os.Open(path)
	if err != nil {
		w.Write([]byte("error for reading file"))
		return
	}
	defer file.Close()

	typ := r.FormValue("type")
	var img image.Image
	switch typ {
	case "bmp":
		img, err = bmp.Decode(file)
		if err != nil {
			return
		}
	default:
		img, err = jpeg.Decode(file)
		if err != nil {
			return
		}
	}

	slc, err := getByteImage(img)
	if err != nil {
		w.Write([]byte("getBinarySlice error"))
		return
	}

	rectangles := scan(slc.arr)

	for _, rect := range rectangles {
		slc.SetRectangle(rect.start.x, rect.start.y, rect.end.x, rect.end.y, 2)
	}

	imgRgba := slc.convertOwnBytesToImage()

	png.Encode(w, imgRgba)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`
	<html>
	<body>
		<form action="/upload" method="post" enctype="multipart/form-data">
			Image: <input type="file" name="my_file">
			<input type="submit" value="Upload">
		</form>
	</body>
	</html>
	`))
}

func loadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 * 1024 * 1024)
	file, handler, err := r.FormFile("my_file")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer file.Close()

	newFile, err := os.Create("load_data/" + handler.Filename)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer newFile.Close()
	io.Copy(newFile, file)
}
