package ui

import "fmt"

// UpdateUI listens for caption strings and displays them.
func UpdateUI(uiUpdateChan <-chan string, quitChan <-chan struct{}) {
	for {
		select {
		case caption := <-uiUpdateChan:
			fmt.Println("Caption:", caption)
		case <-quitChan:
			fmt.Println("Stopping UI updater.")
			return
		}
	}
}
