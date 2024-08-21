package format

import (
	"bytes"
	"fmt"
	"os"
)

func ExamplePNGImage() {
	pngImg := NewPNGImage()
	imgPath := "../testdata/test.png"
	img, err := os.ReadFile(imgPath)
	if err != nil {
		fmt.Printf("Read image file failed: %s", err.Error())
	}

	file, err := os.Create("../testdata/test_resized.png")
	if err != nil {
		fmt.Printf("Open image file failed: %s", err.Error())
	}
	defer file.Close()

	w, h := pngImg.BestImageSize(100, 100, 1920, 1080)
	fmt.Println(w, h)
	if err = pngImg.ResizeImage(bytes.NewReader(img), w, h, file); err != nil {
		fmt.Printf("Resize image failed: %s", err.Error())
	}
	// Output:1186 667
}

func ExampleToJPG() {
	imgPath := "../testdata/test.png"
	img, err := os.ReadFile(imgPath)
	if err != nil {
		fmt.Printf("Read image file failed: %s", err.Error())
	}

	file, err := os.Create("../testdata/test_jpg.jpg")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if err = ToJPG(bytes.NewReader(img), file); err != nil {
		panic(err)
	}
	// Output:
}
