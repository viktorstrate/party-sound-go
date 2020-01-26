package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/hajimehoshi/oto"
)

func main() {

	resp, err := http.Get("http://localhost:8080/audio")
	if err != nil {
		chk(err)
	}

	sampleStr := resp.Header.Get("X-Samplerate")
	if len(sampleStr) == 0 {
		panic("Sample rate header was not set")
	}
	sampleRate, err := strconv.Atoi(sampleStr)
	if err != nil {
		chk(err)
	}

	fmt.Printf("Samplerate: %d\n", sampleRate)

	startStr := resp.Header.Get("X-Start-Time")
	if len(startStr) == 0 {
		panic("Start-Time header was not set")
	}
	startTimeRaw, err := strconv.Atoi(startStr)
	if err != nil {
		chk(err)
	}
	startTimeSec := int64(time.Duration(startTimeRaw) / time.Second)
	startTime := time.Unix(startTimeSec, int64(startTimeRaw)-startTimeSec*int64(time.Second))

	startDelay := startTime.Sub(time.Now())
	startTimer := time.NewTimer(startDelay)
	fmt.Printf("Start delay %d ms", startDelay/time.Millisecond)

	<-startTimer.C
	p, err := oto.NewPlayer(sampleRate, 2, 2, 8192)
	if err != nil {
		chk(err)
	}
	defer p.Close()

	io.Copy(p, resp.Body)

	fmt.Println("Done")
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
