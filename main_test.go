package main

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/pkg/errors"
)

func makeTempImageDir(imagePaths []string) (string, error) {
	path, err := ioutil.TempDir("", "images-")
	if err != nil {
		return "", err
	}

	for _, imagePath := range imagePaths {
		fullPath := filepath.Join(path, imagePath)
		err := os.MkdirAll(filepath.Dir(fullPath), 0777)
		if err != nil {
			return "", errors.Wrap(err, "failed to create directories")
		}
		err = ioutil.WriteFile(fullPath, []byte{}, 0666)
		if err != nil {
			return "", errors.Wrap(err, "failed to write image file")
		}
	}

	return path, nil
}

func TestGetImagesInDirReturnsErrorWhenDirDoesNotExist(t *testing.T) {
	nonExistentDir := "/thisDirDoesNotExist"
	_, err := GetImagesInDir(nonExistentDir)

	assert.NotNil(t, err)
}

func TestGetImagesInDirWorksWithTwoJpgs(t *testing.T) {
	tmpImagesDir, err := makeTempImageDir([]string{"portrait/image1.jpg", "landscape/image2.jpg"})
	assert.Nil(t, err)
	defer os.RemoveAll(tmpImagesDir)

	imagePaths, err := GetImagesInDir(tmpImagesDir)

	assert.Equal(t, 2, len(imagePaths))
	assert.Contains(t, imagePaths, filepath.Join(tmpImagesDir, "portrait/image1.jpg"))
	assert.Contains(t, imagePaths, filepath.Join(tmpImagesDir, "landscape/image2.jpg"))
	assert.Nil(t, err)
}

func TestGetImagesInDirReturnsNoResultsForNonImageFiles(t *testing.T) {
	tmpImagesDir, err := makeTempImageDir([]string{"portrait/image1.png", "landscape/image2.png", "landscape/readme.md"})
	assert.Nil(t, err)
	defer os.RemoveAll(tmpImagesDir)

	imagePaths, err := GetImagesInDir(tmpImagesDir)

	assert.Equal(t, 0, len(imagePaths))
	assert.Nil(t, err)
}