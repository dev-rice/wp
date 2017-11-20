package wallpaper

import (
	"os/exec"
	"fmt"
	"github.com/pkg/errors"
)

func Set(path string) error {
	cmd := exec.Command("osascript", "-e", fmt.Sprintf("tell application \"Finder\" to set desktop picture to POSIX file \"%s\"", path))
	err := cmd.Run()
	return errors.Wrap(err, "failed to run wallpaper set command")
}