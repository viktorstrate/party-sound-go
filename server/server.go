package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
)

const sampleRate = 44100
const seconds = 1

func main() {

	f, err := os.Open("music.mp3")
	if err != nil {
		chk(err)
	}
	defer f.Close()

	mp3Decoder, err := mp3.NewDecoder(f)
	var mp3buffer bytes.Buffer

	fmt.Println("Filling mp3 buffer")

	io.Copy(&mp3buffer, mp3Decoder)

	fmt.Println("Starting player")

	player, err := oto.NewPlayer(sampleRate, 2, 2, 8192)
	if err != nil {
		chk(err)
	}
	defer player.Close()

	io.Copy(player, &mp3buffer)

	http.HandleFunc("/audio", func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("Client streaming")

		if err != nil {
			chk(err)
		}

		w.Header().Set("Connection", "Keep-Alive")
		w.Header().Set("Transfer-Encoding", "chunked")
		w.Header().Set("X-Samplerate", fmt.Sprint(mp3Decoder.SampleRate()))
		w.Header().Set("X-Start-Time", "100")

		fmt.Println("Stream ended")

	})

	http.HandleFunc("/samplerate", func(w http.ResponseWriter, r *http.Request) {

	})

	fmt.Println("Listening on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		chk(err)
	}

}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
