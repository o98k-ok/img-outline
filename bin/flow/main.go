package main

import (
	"bytes"
	"errors"
	"image"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/o98k-ok/img-outline/format"
	"github.com/o98k-ok/img-outline/merge"
	"github.com/o98k-ok/lazy/v2/alfred"
	"golang.design/x/clipboard"
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
		runScreencapture(conf[SCREEN_SHOT])
		outline(s)
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
	return items
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
	jpgHandler := format.NewJPGImage()
	{
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
		frontdata = clipboard.Read(clipboard.FmtImage)
		_, _, err := image.Decode(bytes.NewReader(frontdata))
		if err != nil {
			alfred.ErrItems("readfront", err).Show()
			return
		}

		var d bytes.Buffer
		if err = format.ToJPG(bytes.NewReader(frontdata), &d); err != nil {
			alfred.ErrItems("toJPG", err).Show()
			return
		}
		frontdata = d.Bytes()
		fw, fh = jpgHandler.ImageSize(bytes.NewReader(frontdata))

		if fw >= bw || fh >= bh {
			fw, fh = jpgHandler.BestImageSize(fw, fh, bw, bh)
			var resizeWriter bytes.Buffer
			if err = jpgHandler.ResizeImage(bytes.NewReader(frontdata), fw, fh, &resizeWriter); err != nil {
				alfred.ErrItems("resize", err).Show()
				return
			}
			frontdata = resizeWriter.Bytes()
		}
	}

	var file bytes.Buffer
	x, y := merge.CenterCoordinate(fw, fh, bw, bh)
	if err = merge.AppendOutline(frontdata, backdata, x, y, &file); err != nil {
		alfred.ErrItems("outline", err).Show()
		return
	}

	clipboard.Write(clipboard.FmtImage, file.Bytes())
}
