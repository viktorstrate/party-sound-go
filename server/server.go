package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/hajimehoshi/go-mp3"
)

const CHUNK_SIZE = 4608
const TIME_DELAY = 2 * time.Second

type SongData struct {
	decoder          mp3.Decoder
	timestamp        time.Time
	timestamp_update time.Time
	lock             sync.Mutex
}

func loadSongData() *SongData {
	f, err := os.Open("music.mp3")
	if err != nil {
		chk(err)
	}

	mp3decoder, err := mp3.NewDecoder(f)

	song_data := SongData{
		decoder:          *mp3decoder,
		timestamp:        time.Now(),
		timestamp_update: time.Now(),
	}

	return &song_data
}

func start_streamer(song_data *SongData, data_channels *[]chan []byte, decode_mutex *sync.Mutex) {
	fmt.Println("Starting streaming")

	for true {
		for song_data.timestamp.Sub(time.Now()) < 0 {
			decode_mutex.Lock()

			// Send chunk
			chunk := make([]byte, CHUNK_SIZE)
			bytes, err := song_data.decoder.Read(chunk)
			if err != nil {
				chk(err)
			}

			for _, client := range *data_channels {
				client <- chunk
			}

			var duration time.Duration = time.Duration(float64(time.Second) * float64(bytes) / float64(song_data.decoder.SampleRate()) / 4.0)

			// fmt.Printf("Read %d bytes, duration %d ms\n", bytes, duration/time.Millisecond)
			song_data.timestamp = song_data.timestamp.Add(duration)
			song_data.timestamp_update = time.Now()

			decode_mutex.Unlock()
		}

		wait_timer := time.NewTimer(song_data.timestamp.Sub(time.Now()))
		<-wait_timer.C
	}
}

func main() {

	song_data := loadSongData()

	data_channels := make([]chan []byte, 0, CHUNK_SIZE)
	var decode_mutex sync.Mutex

	go start_streamer(song_data, &data_channels, &decode_mutex)

	http.HandleFunc("/audio", func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("Client streaming")

		w.Header().Set("Connection", "Keep-Alive")
		w.Header().Set("Transfer-Encoding", "chunked")

		// Lock to initialize client
		decode_mutex.Lock()

		w.Header().Set("X-Samplerate", fmt.Sprint(song_data.decoder.SampleRate()))

		timeDrift := time.Now().Sub(song_data.timestamp_update)
		fmt.Printf("Time drift: %d", timeDrift)

		startTime := song_data.timestamp
		startTime = startTime.Add(TIME_DELAY)
		startTime = startTime.Add(timeDrift)

		w.Header().Set("X-Start-Time", fmt.Sprint(startTime.UnixNano()))

		client_chan := make(chan []byte)
		data_channels = append(data_channels, client_chan)

		decode_mutex.Unlock()

		defer func() {
			decode_mutex.Lock()
			defer decode_mutex.Unlock()

			index := -1
			for i, channel := range data_channels {
				if channel == client_chan {
					index = i
					break
				}
			}

			if index == -1 {
				fmt.Println("ERROR: Disconnecting client's channel not found")
				return
			}

			// Remove channel from array
			data_channels[index] = data_channels[len(data_channels)-1]
			data_channels = data_channels[:len(data_channels)-1]

			fmt.Println("Client disconnect cleanup succeeded")

		}()

		for true {
			chunk := <-client_chan
			_, err := w.Write(chunk)
			if err != nil {
				fmt.Println("Could not write chunk to client:")
				fmt.Println(err.Error())
				fmt.Println("Closing connection")
				break
			}
		}

		fmt.Println("Stream ended")

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
