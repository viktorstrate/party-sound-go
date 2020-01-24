package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

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
