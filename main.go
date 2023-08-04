package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"flag"
	"fmt"
	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	vaudio "github.com/voicedock/audio"
	ttsv1 "github.com/voicedock/go-text-to-wav/internal/api/grpc/gen/voicedock/core/tts/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"os"
)

var addr string
var lang string
var speaker string
var cmd string

func init() {
	flag.StringVar(&addr, "a", "0.0.0.0:9999", "TTS gRPC server host:port")
	flag.StringVar(&lang, "l", "en", "Language code")
	flag.StringVar(&speaker, "s", "", "Speaker name")
	flag.StringVar(&cmd, "c", "", "Command: `list` or `download` voices")
	flag.Parse()
}

func main() {
	conn, err := grpc.DialContext(
		context.TODO(), addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock(),
	)
	chk(err)
	ttsClient := ttsv1.NewTtsAPIClient(conn)

	switch cmd {
	case "list":
		ret, err := ttsClient.GetVoices(context.TODO(), &ttsv1.GetVoicesRequest{})
		chk(err)

		fmt.Print("Lang\tSpeaker\tDownloaded\n")
		for _, v := range ret.Voices {
			fmt.Printf("%s\t%s\t%t\n", v.Lang, v.Speaker, v.Downloaded)
		}
		return
	case "download":
		fmt.Printf("Staring download (lang: %s, speaker: %s)\n", lang, speaker)
		_, err := ttsClient.DownloadVoice(context.TODO(), &ttsv1.DownloadVoiceRequest{
			Lang: lang,
			Speaker: speaker,
		})
		chk(err)
		fmt.Print("Download complete\n")
		return
	}

	// convert text to wav
	scanner := bufio.NewScanner(os.Stdin)
	inputText := ""
	for scanner.Scan() {
		inputText += scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	srv, err := ttsClient.TextToSpeech(context.TODO(), &ttsv1.TextToSpeechRequest{
		Text:    inputText,
		Lang:    lang,
		Speaker: speaker,
	})
	chk(err)

	var enc *wav.Encoder
	var sampleRate int
	var channels int
	initialized := false

	for {
		resp, err := srv.Recv()
		if err == io.EOF {
			break
		}
		chk(err)

		if !initialized {
			initialized = true
			sampleRate = int(resp.Audio.SampleRate)
			channels = int(resp.Audio.Channels)
			enc = wav.NewEncoder(os.Stdout, sampleRate, 16, channels, 1)
		}

		// convert []byte to []int
		r := new(bytes.Buffer)
		r.Write(resp.Audio.Data)
		out := make([]int16, len(resp.Audio.Data)/2)
		err = binary.Read(r, binary.LittleEndian, out)
		chk(err)
		intData := vaudio.ConvertNumbers[int](out)

		// write data
		enc.Write(&audio.IntBuffer{
			Format:         &audio.Format{
				NumChannels: 1,
				SampleRate:  sampleRate,
			},
			Data:           intData,
			SourceBitDepth: 16,
		})
	}

	enc.Close()
	os.Stdout.Sync()
	os.Stdout.Close()
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}