package utils

import (
  "fmt"
  // "image"
  "os"
  // "os/exec"
  // "github.com/nfnt/resize"
  "bytes"
  "github.com/disintegration/imaging"
  ffmpeg "github.com/u2takey/ffmpeg-go"
)

func GetSnapShot(inFileName, imageName string, frameNum int) (string, error) {
  outpath := "./public/covers/" + imageName
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(inFileName).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		return "", err
	}
  
  img, err := imaging.Decode(buf)
  if err != nil {
      return "", err
  }
  err = imaging.Save(img, outpath)
  if err != nil {
      return "", err
  }
	return outpath, nil
}

