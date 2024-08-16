package merge

import (
	"fmt"
	"os"
)

func ExampleAppend() {
	front := "../testdata/test_jpg.jpg"
	back := "../testdata/macos-big-sur-dark.jpg"

	frontData, err := os.ReadFile(front)
	if err != nil {
		panic(err)
	}
	backData, err := os.ReadFile(back)
	if err != nil {
		panic(err)
	}

	file, err := os.Create("../testdata/test_append.jpg")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	x, y := CenterCoordinate(1310, 832, 1920, 1080)
	fmt.Println(x, y)
	if err = AppendOutline(frontData, backData, x, y, file); err != nil {
		panic(err)
	}
	// Output:305 124
}
