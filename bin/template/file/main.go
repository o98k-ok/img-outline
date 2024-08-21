package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path"
	"time"

	"math/rand"

	"github.com/o98k-ok/img-outline/format"
	"github.com/o98k-ok/img-outline/merge"
)

func main() {
	var front string
	var out string
	var style string
	var err error
	flag.StringVar(&front, "front", "", "Front image path")
	flag.StringVar(&out, "out", "", "Output Image path")
	flag.StringVar(&style, "style", "macos", "Merge style: macos/raycast")
	flag.Parse()

	if front == "" || out == "" {
		flag.Usage()
		return
	}
	var back string
	{
		dir := path.Join("./bin/template/imgs", style)
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
	jpgHandler := format.NewJPGImage()
	{
		backdata, err = os.ReadFile(back)
		if err != nil {
			fmt.Println(front, err)
			flag.Usage()
			return
		}

		var d bytes.Buffer
		if err = format.ToJPG(bytes.NewReader(backdata), &d); err != nil {
			fmt.Println(back, err)
			flag.Usage()
			return
		}
		backdata = d.Bytes()
		bw, bh = jpgHandler.ImageSize(bytes.NewReader(backdata))
	}

	{
		frontdata, err = os.ReadFile(front)
		if err != nil {
			fmt.Println(front, err)
			flag.Usage()
			return
		}

		var d bytes.Buffer
		if err = format.ToJPG(bytes.NewReader(frontdata), &d); err != nil {
			fmt.Println(front, err)
			flag.Usage()
			return
		}
		frontdata = d.Bytes()
		fw, fh = jpgHandler.ImageSize(bytes.NewReader(frontdata))

		if fw >= bw || fh >= bh {
			fw, fh = jpgHandler.BestImageSize(fw, fh, bw, bh)
			var resizeWriter bytes.Buffer
			if err = jpgHandler.ResizeImage(bytes.NewReader(frontdata), fw, fh, &resizeWriter); err != nil {
				fmt.Println("resize image: ", err)
				return
			}
			frontdata = resizeWriter.Bytes()
		}
	}

	file, err := os.Create(out)
	if err != nil {
		fmt.Println(out, err)
		flag.Usage()
		return
	}
	defer file.Close()

	x, y := merge.CenterCoordinate(fw, fh, bw, bh)
	if err = merge.AppendOutline(frontdata, backdata, x, y, file); err != nil {
		fmt.Println(out, err)
		flag.Usage()
		return
	}
}
