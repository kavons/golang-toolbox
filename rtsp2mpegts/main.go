package main

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func main()  {
	svc := &Rtsp2Mpegts{
		"rtsp://wowzaec2demo.streamlock.net/vod/mp4:BigBuckBunny_115k.mov",
	}

	err := svc.Convert()
	if err != nil {
		log.Fatal(err)
	}
	return
}

type Rtsp2Mpegts struct {
	RtspURL string
}

func (svc *Rtsp2Mpegts) Convert() (err error) {
	simpleString := strings.Replace(svc.RtspURL, "//", "/", 1)
	splitList := strings.Split(simpleString, "/")

	if splitList[0] != "rtsp:" && len(splitList) < 2 {
		err = errors.New("invalid RTSP address")
		return
	}

	err = FFMPEG(svc.RtspURL)
	return
}

func FFMPEG(rtsp string) (err error) {
	params := []string{
		"-rtsp_transport",
		"tcp",
		"-i",
		rtsp,
		"-f",
		"mpegts",
		"-fs",
		"128K",
		"./demo.ts",
	}

	fmt.Printf("FFMPEG cmd: ffmpeg %v\n", strings.Join(params, " "))
	cmd := exec.Command("ffmpeg", params...)

	err = cmd.Start()
	if err != nil {
		return
	}

	fmt.Printf("waiting for cmd to finish...\n")
	err = cmd.Wait()
	return
}
