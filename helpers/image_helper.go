package helpers

import (
	"doctorrank_go/configs"
	"github.com/h2non/bimg"
	"strconv"
	"time"
)

// ImageProcessing
func ProcessAndSaveAvatar(buffer []byte, name string, top int, left int, width int, height int) (string, error) {
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

	writeError := bimg.Write(path+"/user/avatar/"+filename, extracted)
	if writeError != nil {
		return filename, writeError
	}

	writeError = bimg.Write(path+"/user/thumbnail/"+filename, processed)
	if writeError != nil {
		return filename, writeError
	}

	return filename, nil
}
