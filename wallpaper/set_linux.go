package wallpaper

import (
	"os/exec"
	"github.com/pkg/errors"
)

func Set(path string) error {
	cmd := exec.Command("feh", "--bg-fill", path)
	err := cmd.Run()
	return errors.Wrap(err, "failed to run wallpaper set command")
}
