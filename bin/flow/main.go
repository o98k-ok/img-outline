package main

import (
	"bytes"
	"errors"
	"image"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/o98k-ok/img-outline/format"
	"github.com/o98k-ok/img-outline/merge"
	"github.com/o98k-ok/lazy/v2/alfred"
	"golang.design/x/clipboard"
	"golang.org/x/exp/rand"
)

var (
	BACK_IMAGES = "background_path"
	SCREEN_SHOT = "screen_shot_cmd"
)

func main() {
	conf, err := alfred.FlowVariables()
	if err != nil {
		alfred.ErrItems("get flow variables", err).Show()
	}

	cli := alfred.NewApp("vscode util toools")
	cli.Bind("images", func(s []string) { backgroundImages(conf[BACK_IMAGES]) })
	cli.Bind("outline", func(s []string) {
		cmd := conf[SCREEN_SHOT]
		if len(cmd) != 0 {
			runScreencapture(cmd)
		}

		if len(s) == 0 {
			images := listBackgroundImages(conf[BACK_IMAGES])
			if len(images) != 0 {
				s = []string{images[0].Arg}
			}
		}
		outline(s)
		alfred.Log("shot out finish")
	})
	cli.Run(os.Args)
}

func runScreencapture(cmd string) {
	running := exec.Command("bash", "-c", cmd)
	out, err := running.CombinedOutput()
	alfred.Log("shot out msg=%s err=%v", strings.TrimSpace(string(out)), err)
}

func backgroundImages(imagePath string) {
	items := alfred.NewItems()
	items.Items = listBackgroundImages(imagePath)
	items.Show()
}

func listBackgroundImages(imagePath string) []*alfred.Item {
	fs, err := os.ReadDir(imagePath)
	if err != nil {
		alfred.ErrItems("get background images", err).Show()
		return []*alfred.Item{}
	}

	var items []*alfred.Item
	for _, f := range fs {
		full := path.Join(imagePath, f.Name())
		switch {
		case f.IsDir():
			items = append(items, listBackgroundImages(full)...)
		case strings.HasSuffix(full, ".png") || strings.HasSuffix(full, ".jpg"):
			item := alfred.NewItem("", "", full)
			item.Icon = &alfred.Icon{Path: full}
			items = append(items, item)
		default:

		}
	}
	Shuffle(items)
	return items
}

func Shuffle[T any](slice []T) {
	r := rand.New(rand.NewSource(uint64(time.Now().Unix())))
	for n := len(slice); n > 0; n-- {
		randIndex := r.Intn(n)
		slice[n-1], slice[randIndex] = slice[randIndex], slice[n-1]
	}
}

func outline(s []string) {
	if len(s) <= 0 {
		alfred.ErrItems("arg", errors.New("empty args")).Show()
		return
	}

	var back string = s[0]
	var frontdata, backdata []byte
	var fw, fh, bw, bh int
	var err error

	{
		jpgHandler := format.NewJPGImage()
		backdata, err = os.ReadFile(back)
		if err != nil {
			alfred.ErrItems("readback", err).Show()
			return
		}

		var d bytes.Buffer
		if err = format.ToJPG(bytes.NewReader(backdata), &d); err != nil {
			alfred.ErrItems("toJPG", err).Show()
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
			alfred.ErrItems("readfront", err).Show()
			return
		}

		var d bytes.Buffer
		if err = format.ToPNG(bytes.NewReader(frontdata), &d); err != nil {
			alfred.ErrItems("toPNG", err).Show()
			return
		}
		frontdata = d.Bytes()
		fw, fh = pngHandler.ImageSize(bytes.NewReader(frontdata))

		if fw >= bw || fh >= bh {
			fw, fh = pngHandler.BestImageSize(fw, fh, bw, bh)
			var resizeWriter bytes.Buffer
			if err = pngHandler.ResizeImage(bytes.NewReader(frontdata), fw, fh, &resizeWriter); err != nil {
				alfred.ErrItems("resize", err).Show()
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
		alfred.ErrItems("outline", err).Show()
		return
	}

	clipboard.Write(clipboard.FmtImage, file.Bytes())
}
