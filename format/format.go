package format

import (
	"image"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/nfnt/resize"
)

type ImageFormater interface {
	ResizeImage(reader io.Reader, width, height int, imgWriter io.Writer) error
	BestImageSize(int, int) (int, int)
	ImageSize(read io.Reader) (int, int)
}

type JPGImage struct {
	Noise float32
}

func NewJPGImage() ImageFormater {
	return &JPGImage{Noise: 1.618}
}

type PNGImage struct {
	Noise float32
}

func NewPNGImage() ImageFormater {
	return &PNGImage{Noise: 1.618}
}

func (i *JPGImage) BestImageSize(width, height int) (int, int) {
	var resWidth, resHeight int
	if width > height {
		resWidth = int(float32(width) / i.Noise)
		rate := float32(resWidth) / float32(width)
		resHeight = int(float32(height) * rate)
	} else {
		resHeight = int(float32(height) / i.Noise)
		rate := float32(resHeight) / float32(height)
		resWidth = int(float32(width) * rate)
	}
	return resWidth, resHeight
}

func (i *JPGImage) ImageSize(read io.Reader) (int, int) {
	img, _, err := image.Decode(read)
	if err != nil {
		return 0, 0
	}
	return img.Bounds().Dx(), img.Bounds().Dy()
}

func (i *PNGImage) BestImageSize(width, height int) (int, int) {
	var resWidth, resHeight int
	if width > height {
		resWidth = int(float32(width) / i.Noise)
		rate := float32(resWidth) / float32(width)
		resHeight = int(float32(height) * rate)
	} else {
		resHeight = int(float32(height) / i.Noise)
		rate := float32(resHeight) / float32(height)
		resWidth = int(float32(width) * rate)
	}
	return resWidth, resHeight
}

func (i *JPGImage) ResizeImage(reader io.Reader, width, height int, imgWriter io.Writer) error {
	img, _, err := image.Decode(reader)
	if err != nil {
		return err
	}

	imgRes := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)
	if err = jpeg.Encode(imgWriter, imgRes, nil); err != nil {
		return err
	}
	return nil
}
func (i *PNGImage) ResizeImage(reader io.Reader, width, height int, imgWriter io.Writer) error {
	img, _, err := image.Decode(reader)
	if err != nil {
		return err
	}

	imgRes := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)
	if err = png.Encode(imgWriter, imgRes); err != nil {
		return err
	}
	return nil
}

func (i *PNGImage) ImageSize(read io.Reader) (int, int) {
	img, _, err := image.Decode(read)
	if err != nil {
		return 0, 0
	}
	return img.Bounds().Dx(), img.Bounds().Dy()
}

func ToJPG(reader io.Reader, writer io.Writer) error {
	img, _, err := image.Decode(reader)
	if err != nil {
		return err
	}

	return jpeg.Encode(writer, img, nil)
}
