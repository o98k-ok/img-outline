package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"

	"github.com/o98k-ok/img-outline/format"
	"github.com/o98k-ok/img-outline/merge"
)

func main() {
	var front string
	var back string
	var out string
	var err error
	flag.StringVar(&front, "front", "", "Front image path")
	flag.StringVar(&back, "back", "", "background Image path")
	flag.StringVar(&out, "out", "", "Output Image path")
	flag.Parse()

	if front == "" || back == "" || out == "" {
		flag.Usage()
		return
	}

	var frontdata, backdata []byte
	var fw, fh, bw, bh int
	jpgHandler := format.NewJPGImage()
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
	}

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
