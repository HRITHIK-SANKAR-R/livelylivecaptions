package multiplexer

import "fmt"

// MultiplexAudio merges audio from mic and system channels into a single outgoing channel.
func MultiplexAudio(micAudioChan, systemAudioChan <-chan []byte, outgoingChan chan<- []byte, quitChan <-chan struct{}) {
	for {
		select {
		case audioData := <-micAudioChan:
			outgoingChan <- audioData
		case audioData := <-systemAudioChan:
			outgoingChan <- audioData
		case <-quitChan:
			fmt.Println("Stopping audio multiplexer.")
			return
		}
	}
}
