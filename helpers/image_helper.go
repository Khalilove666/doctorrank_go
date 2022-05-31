package helpers

import (
	"doctorrank_go/configs"
	"github.com/h2non/bimg"
	"strconv"
	"time"
)

var Folders = newFolderRegistry()

func newFolderRegistry() *folderRegistry {
	return &folderRegistry{
		User:     "user",
		Doctor:   "doctor",
		Hospital: "hospital",
	}
}

type folderRegistry struct {
	User     string
	Doctor   string
	Hospital string
}

// ImageProcessing
func ProcessAndSaveAvatar(buffer []byte, name string, directory string, top int, left int, width int, height int) (string, error) {
	path := configs.Env("FILESYSTEM_PATH")
	timeNow := strconv.FormatInt(time.Now().Unix(), 10)
	filename := name + "_" + timeNow + ".webp"

	converted, err := bimg.NewImage(buffer).Convert(bimg.WEBP)
	if err != nil {
		return filename, err
	}
	extracted, err := bimg.NewImage(converted).Extract(top, left, width, height)
	if err != nil {
		return filename, err
	}
	processed, err := bimg.NewImage(extracted).Thumbnail(256)
	if err != nil {
		return filename, err
	}

	writeError := bimg.Write(path+"/"+directory+"/avatar/"+filename, extracted)
	if writeError != nil {
		return filename, writeError
	}

	writeError = bimg.Write(path+"/"+directory+"/thumbnail/"+filename, processed)
	if writeError != nil {
		return filename, writeError
	}

	return filename, nil
}
