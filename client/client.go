package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/hajimehoshi/oto"
)

func main() {

	host := "127.0.0.1"

	if len(os.Args) == 2 {
		host = os.Args[1]
	}

	url := fmt.Sprintf("http://%s:8080/audio", host)
	fmt.Printf("Connecting to: %s\n", url)

	resp, err := http.Get(url)
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
	fmt.Printf("Start delay %f ms\n", float64(startDelay)/float64(time.Millisecond))

	p, err := oto.NewPlayer(sampleRate, 2, 2, 8192)
	if err != nil {
		chk(err)
	}
	defer p.Close()

	<-startTimer.C
	io.Copy(p, resp.Body)

	fmt.Println("Done")
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
