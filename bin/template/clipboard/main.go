package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"os"
	"path"
	"time"

	"math/rand"

	"github.com/o98k-ok/img-outline/format"
	"github.com/o98k-ok/img-outline/merge"
	"golang.design/x/clipboard"
)

func main() {
	var style, dir string
	var err error
	flag.StringVar(&style, "style", "macos", "Merge style: macos/raycast")
	flag.StringVar(&dir, "dir", "./bin/template/imgs", "Directory of background images")
	flag.Parse()

	var back string
	{
		dir := path.Join(dir, style)
		files, err := os.ReadDir(dir)
		if err != nil {
			fmt.Println("Error reading directory:", dir, err)
			return
		}
		if len(files) == 0 {
			fmt.Println("No files found in the directory.", dir)
			return
		}
		rand.Seed(time.Now().UnixNano())
		randomIndex := rand.Intn(len(files))
		back = path.Join(dir, files[randomIndex].Name())
	}

	var frontdata, backdata []byte
	var fw, fh, bw, bh int

	{
		jpgHandler := format.NewJPGImage()
		backdata, err = os.ReadFile(back)
		if err != nil {
			fmt.Println("Error reading background image from file:", err)
			return
		}

		var d bytes.Buffer
		if err = format.ToJPG(bytes.NewReader(backdata), &d); err != nil {
			fmt.Println("Error toJPG: ", err)
			flag.Usage()
			return
		}
		backdata = d.Bytes()
		bw, bh = jpgHandler.ImageSize(bytes.NewReader(backdata))
	}

	{
		pngHandler := format.NewPNGImage()
		frontdata = clipboard.Read(clipboard.FmtImage)
		_, _, err := image.Decode(bytes.NewReader(frontdata))
		if err != nil {
			fmt.Println("Error reading image from clipboard:", err)
			return
		}

		var d bytes.Buffer
		if err = format.ToPNG(bytes.NewReader(frontdata), &d); err != nil {
			fmt.Println("Error toJPG: ", err)
			return
		}
		frontdata = d.Bytes()
		fw, fh = pngHandler.ImageSize(bytes.NewReader(frontdata))

		if fw >= bw || fh >= bh {
			fw, fh = pngHandler.BestImageSize(fw, fh, bw, bh)
			var resizeWriter bytes.Buffer
			if err = pngHandler.ResizeImage(bytes.NewReader(frontdata), fw, fh, &resizeWriter); err != nil {
				fmt.Println("resize image: ", err)
				return
			}
			frontdata = resizeWriter.Bytes()
		}

		// 将frontdata的直角图片信息转变成圆角图片信息
		frontdata = merge.RoundCorner(frontdata, fw, fh)
	}

	var file bytes.Buffer
	x, y := merge.CenterCoordinate(fw, fh, bw, bh)
	if err = merge.AppendOutline(frontdata, backdata, x, y, &file); err != nil {
		fmt.Println("Error appending outline:", err)
		return
	}

	clipboard.Write(clipboard.FmtImage, file.Bytes())
}
