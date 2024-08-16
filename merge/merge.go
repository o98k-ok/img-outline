package merge

import (
	"image/jpeg"
	"io"
	"os"

	gim "github.com/ozankasikci/go-image-merge"
)

func AppendOutline(front, back []byte, x, y int, out io.Writer) error {
	frontFile, err := os.CreateTemp("", "front*.jpg")
	if err != nil {
		return err
	}
	frontFile.Write(front)
	defer frontFile.Close()

	backFile, err := os.CreateTemp("", "back*.jpg")
	if err != nil {
		return err
	}
	backFile.Write(back)
	defer backFile.Close()

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
