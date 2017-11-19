package main

import (
	"fmt"
	"github.com/pkg/errors"
	_ "image/jpeg"
	"os"
	"path/filepath"
	"image"

	"github.com/urfave/cli"
	"github.com/donutmonger/wp/wallpaper"
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

func GetImageAspectRatio(path string) float32 {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	img, _, err := image.DecodeConfig(f)
	if err != nil {
		panic(err)
	}

	return float32(img.Width) / float32(img.Height)
}

func PrintListOfImages(paths []string) error {
	for i, path := range paths {
		fmt.Printf("%d: %s\n", i,filepath.Base(path))
	}

	return nil
}

// wp ls - will give a numbered list of all jpegs in the specified directory
// wp set -n <n> - will set wallpaper to the image at wp_ls[n]

const wallpaperDir = "/home/chris/Pictures/Wallpapers"

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
			},
			Action: setWallpaper,
		},
	}

	app.Run(os.Args)
}

func listWallpapers(ctx *cli.Context) error {
	imagePaths, err := GetImagesInDir(wallpaperDir)
	if err != nil {
		return err
	}
	PrintListOfImages(imagePaths)
	return nil
}

func setWallpaper(ctx *cli.Context) error {
	imagePaths, err := GetImagesInDir(wallpaperDir)
	if err != nil {
		return errors.Wrap(err, "failed to get list of images")
	}

	err = wallpaper.Set(imagePaths[ctx.Int("n")])
	return errors.Wrap(err, "failed to set wallpaper")
}