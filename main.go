package main

import (
	"fmt"
	"github.com/pkg/errors"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"image"

	"github.com/donutmonger/wp/wallpaper"
	"github.com/urfave/cli"
	"math/rand"
	"time"
)


func GetImagesInDir(dir string) ([]string, error) {
	var images []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrap(err, "unable to walk file path")
		}
		if filepath.Ext(path) == ".jpg" || filepath.Ext(path) == ".png" {
			images = append(images, filepath.Join(path))
		}
		return nil
	})

	return images, errors.Wrap(err, "Failed to list all wallpapers")
}


func PrintListOfImages(paths []string) error {
	for i, path := range paths {
		fmt.Printf("%d: %s\n", i,filepath.Base(path))
	}

	return nil
}

const wallpaperDirEnvName = "WP_DIR"

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Usage: "give a numbered list of all jpegs and pngs in the specified directory",
			Action:  listWallpapers,
		},
		{
			Name: "set",
			Aliases: []string{"s"},
			Usage: "will set wallpaper to the image at wp_ls[n]",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name: "n",
				},
				cli.BoolFlag{
					Name: "r",
				},
			},
			Action: setWallpaper,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type Orientation int

const (
	Vertical Orientation = iota + 1
	Horizontal
	Square
)

func GetOrientationFromAspectRatio(a float32) Orientation {
	if a < 1 {
		return Vertical
	} else if a > 1 {
		return Horizontal
	} else {
		return Square
	}
}

func GetImagesWithOrientation(images []string, orientation Orientation) []string {
	var filteredImages []string
	for _, i := range images {
		o := GetOrientationFromAspectRatio(GetImageAspectRatio(i))
		if o == orientation {
			filteredImages = append(filteredImages, i)
		}
	}
	return filteredImages
}

func GetImageAspectRatio(path string) float32 {
	f, err := os.Open(path)
	if err != nil {
		panic(errors.Wrapf(err, "failed to open file %s", path))
	}
	img, _, err := image.DecodeConfig(f)
	if err != nil {
		panic(errors.Wrapf(err, "failed to decode config %s", path))
	}

	return float32(img.Width) / float32(img.Height)
}

func listWallpapers(ctx *cli.Context) error {
	wallpaperDir := os.Getenv(wallpaperDirEnvName)
	if wallpaperDir == "" {
		return fmt.Errorf("%s env var must be provided", wallpaperDirEnvName)
	}
	fmt.Println(wallpaperDir)
	imagePaths, err := GetImagesInDir(wallpaperDir)
	if err != nil {
		return err
	}

	filteredImagePaths := GetImagesWithOrientation(imagePaths, Horizontal)

	PrintListOfImages(filteredImagePaths)
	return nil
}

func setWallpaper(ctx *cli.Context) error {
	wallpaperDir := os.Getenv(wallpaperDirEnvName)
	if wallpaperDir == "" {
		return fmt.Errorf("%s env var must be provided", wallpaperDirEnvName)
	}
	imagePaths, err := GetImagesInDir(wallpaperDir)
	if err != nil {
		return errors.Wrap(err, "failed to get list of images")
	}

	filteredImagePaths := GetImagesWithOrientation(imagePaths, Horizontal)

	wallpaperNum := 0
	if ctx.IsSet("r") {
		rand.Seed(time.Now().UnixNano())
		wallpaperNum = rand.Intn(len(filteredImagePaths))
	} else if ctx.IsSet("n") {
		wallpaperNum = ctx.Int("n")
	} else {
		return errors.New("-n or -r must be set")
	}

	err = wallpaper.Set(filteredImagePaths[wallpaperNum])
	return errors.Wrap(err, "failed to set wallpaper")
}
