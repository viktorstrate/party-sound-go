package main

import (
	"net/http"

	"github.com/hajimehoshi/oto"
)

func main() {

	resp, err := http.Get("localhost:8080/samplerate")
	if err != nil {
		chk(err)
	}
	defer resp.Body.Close()

	samplerate := string(resp)

	p, err := oto.NewPlayer(d.SampleRate(), 2, 2, 8192)
	if err != nil {
		return err
	}
	defer p.Close()
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
