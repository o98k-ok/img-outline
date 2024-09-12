package merge

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"os"

	gim "github.com/ozankasikci/go-image-merge"
)

func AppendOutline(front, back []byte, x, y int, out io.Writer) error {
	var frontFile, backFile *os.File
	var err error
	{
		_, format, err := image.Decode(bytes.NewReader(front))
		if err != nil {
			return err
		}

		fileFormat := "front*.png"
		switch format {
		case "png":
			fileFormat = "front*.png"
		case "jpg", "jpeg":
			fileFormat = "front*.jpg"
		default:
		}

		frontFile, err = os.CreateTemp("", fileFormat)
		if err != nil {
			return err
		}
		frontFile.Write(front)
		frontFile.Close()
	}

	{

		_, format, err := image.Decode(bytes.NewReader(back))
		if err != nil {
			return err
		}

		fileFormat := "back*.jpg"
		switch format {
		case "png":
			fileFormat = "back*.png"
		case "jpg", "jpeg":
			fileFormat = "back*.jpg"
		default:
		}

		backFile, err = os.CreateTemp("", fileFormat)
		if err != nil {
			return err
		}
		backFile.Write(back)
		backFile.Close()
	}

	grids := []*gim.Grid{
		{
			ImageFilePath: backFile.Name(),
			Grids: []*gim.Grid{
				{
					ImageFilePath: frontFile.Name(),
					OffsetX:       x, OffsetY: y,
				},
			},
		},
	}

	rgba, err := gim.New(grids, 1, 1).Merge()
	if err != nil {
		return err
	}
	return jpeg.Encode(out, rgba, nil)
}

func CenterCoordinate(fw, fh, bw, bh int) (int, int) {
	x := (bw - fw) / 2
	y := (bh - fh) / 2
	return x, y
}

func RoundCorner(front []byte, fw, fh int) []byte {
	img, err := png.Decode(bytes.NewReader(front))
	if err != nil {
		return nil
	}

	radius := 40 // 圆角半径，可以根据需要调整
	mask := image.NewRGBA(img.Bounds())

	// 绘制圆角遮罩
	for y := 0; y < fh; y++ {
		for x := 0; x < fw; x++ {
			if (x < radius && y < radius && (x-radius)*(x-radius)+(y-radius)*(y-radius) > radius*radius) ||
				(x < radius && y >= fh-radius && (x-radius)*(x-radius)+(y-(fh-radius))*(y-(fh-radius)) > radius*radius) ||
				(x >= fw-radius && y < radius && (x-(fw-radius))*(x-(fw-radius))+(y-radius)*(y-radius) > radius*radius) ||
				(x >= fw-radius && y >= fh-radius && (x-(fw-radius))*(x-(fw-radius))+(y-(fh-radius))*(y-(fh-radius)) > radius*radius) {
				mask.Set(x, y, color.Transparent)
			} else {
				mask.Set(x, y, img.At(x, y))
			}
		}
	}

	var buf bytes.Buffer
	png.Encode(&buf, mask)
	return buf.Bytes()
}
