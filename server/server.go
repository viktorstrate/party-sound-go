package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/hajimehoshi/go-mp3"
)

const sampleRate = 44100
const seconds = 1

func main() {
	f, err := os.Open("music.mp3")
	if err != nil {
		chk(err)
	}

	defer f.Close()

	d, err := mp3.NewDecoder(f)
	if err != nil {
		chk(err)
	}

	http.HandleFunc("/audio", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Connection", "Keep-Alive")
		w.Header().Set("Transfer-Encoding", "chunked")

		io.Copy(w, d)

	})

	http.HandleFunc("/samplerate", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprint(d.SampleRate())))
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
