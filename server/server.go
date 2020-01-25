package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/hajimehoshi/go-mp3"
)

const sampleRate = 44100
const seconds = 1

type SongData struct {
	buffer    bytes.Buffer
	decoder   mp3.Decoder
	timestamp time.Time
	lock      sync.Mutex
}

func loadSongData() *SongData {
	f, err := os.Open("music.mp3")
	if err != nil {
		chk(err)
	}
	defer f.Close()

	mp3decoder, err := mp3.NewDecoder(f)
	var mp3buffer bytes.Buffer

	fmt.Println("Filling mp3 buffer")

	io.Copy(&mp3buffer, mp3decoder)

	song_data := SongData{
		buffer:    mp3buffer,
		decoder:   *mp3decoder,
		timestamp: time.Now(),
	}

	return &song_data
}

func main() {

	song_data := loadSongData()

	http.HandleFunc("/audio", func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("Client streaming")

		w.Header().Set("Connection", "Keep-Alive")
		w.Header().Set("Transfer-Encoding", "chunked")
		w.Header().Set("X-Samplerate", fmt.Sprint(song_data.decoder.SampleRate()))
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
