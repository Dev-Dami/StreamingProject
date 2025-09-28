package video

import (
	"bufio"
	"os/exec"
)

func SaveToDAT(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for frame := range QuantizedFrameChan {
		writer.Write(frame)
	}
	return nil
}

func ReconstructMP4(datFile, mp4File string) {

	// to be added split .dat into jpegs the run
	log.Println("Reconstruction tobe added")
}


