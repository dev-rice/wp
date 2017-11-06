package main

import (
	"fmt"
	"github.com/pkg/errors"
	_ "image/jpeg"
	"os"
	"path/filepath"
	"image"
	"os/exec"

	"github.com/urfave/cli"
)

func GetImagesInDir(dir string) ([]string, error) {
	images := []string{}
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrap(err, "unable to walk file path")
		}
		if filepath.Ext(path) == ".jpg" {
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

func setWallpaper(path string) error {
	cmd := exec.Command("osascript", "-e", fmt.Sprintf("tell application \"Finder\" to set desktop picture to POSIX file \"%s\"", path))
	err := cmd.Run()
	return errors.Wrap(err, "failed to run wallpaper set command")
}

func setWallpaperLinux(path string) error {
    cmd := exec.Command("feh", "--bg-fill", path)
    err := cmd.Run()
    return errors.Wrap(err, "failed to run wallpaper set command")
}

func PrintListOfImages(paths []string) error {
	for i, path := range paths {
		fmt.Printf("%d: %s\n", i, filepath.Base(path))
	}

	return nil
}

// wp ls - will give a numbered list of all jpegs in the specified directory
// wp set -n <n> - will set wallpaper to the image at wp_ls[n]

const wallpaperDir = "/home/chris.rice/Pictures/Wallpapers"

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Usage: "give a numbered list of all jpegs in the specified directory",
			Action:  func(ctx *cli.Context) error {
				imagePaths, err := GetImagesInDir(wallpaperDir)
				if err != nil {
					return err
				}
				PrintListOfImages(imagePaths)
				return nil
			},
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
			Action: func(ctx *cli.Context) error {
				imagePaths, err := GetImagesInDir(wallpaperDir)
				if err != nil {
					return errors.Wrap(err, "failed to get list of images")
				}

				err = setWallpaperLinux(imagePaths[ctx.Int("n")])
				return errors.Wrap(err, "failed to set wallpaper")
			},
		},
	}

	app.Run(os.Args)
}
