package rabbit

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

// Input is a base64 encoded string
func WebMToWav(input string) ([]byte, error) {
	webmData, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		log.Println("error decoding base64 string:", err)
		return nil, err
	}

	webmFile, err := ioutil.TempFile("", "input-*.webm")
	if err != nil {
		return nil, err
	}
	defer os.Remove(webmFile.Name())

	if _, err = webmFile.Write(webmData); err != nil {
		return nil, err
	}
	if err = webmFile.Close(); err != nil {
		return nil, err
	}

	// Create a temporary file for the output WAV data
	wavFile, err := ioutil.TempFile("", "output-*.wav")
	if err != nil {
		return nil, err
	}
	defer os.Remove(wavFile.Name())

	wavFileName := wavFile.Name()
	wavFile.Close()

	// Use ffmpeg to convert the WebM file to a WAV file
	cmd := exec.Command("ffmpeg", "-i", webmFile.Name(), wavFileName)
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	// Read the resulting WAV file into a []byte
	wavData, err := ioutil.ReadFile(wavFileName)
	if err != nil {
		return nil, err
	}

	if len(wavData) == 0 {
		return nil, err
	}

	return wavData, nil
}
