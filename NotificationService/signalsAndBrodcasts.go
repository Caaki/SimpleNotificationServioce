package main

import (
	"NotificationService/listener"
	"fmt"
	"sync"
	"time"
)

var beacon = sync.NewCond(&sync.Mutex{})
var notReadMessages = make(map[string][]string)
var addingMessageBeacpm = sync.NewCond(&sync.Mutex{})

func main() {

	go listener.Listen()

	wg := sync.WaitGroup{}

	go addMessageForRecipientsLoop(&wg)
	processMessages(&wg)

	wg.Wait()
}

func processMessages(wg *sync.WaitGroup) {
	wg.Add(1)
	for {
		beacon.L.Lock()
		for len(notReadMessages) == 0 {
			beacon.Wait()
		}

		beacon.L.Unlock()
		sendMessageToRecipiants(beacon)
	}
	wg.Done()
}

func addMessageForRecipientsLoop(wg *sync.WaitGroup) {

	for {
		fmt.Print("Unesite broj primaoca: ")
		countOfRecipients := 0
		fmt.Scan(&countOfRecipients)

		recipients := make([]string, 0)

		fmt.Print("Unesite poruku : ")
		messageForUsers := ""
		fmt.Scanln(&messageForUsers)

		for range countOfRecipients {
			name := ""
			fmt.Print("Unesite ime primaoca: ")
			fmt.Scanln(&name)
			recipients = append(recipients, name)
		}

		addingMessageBeacpm.L.Lock()
		for _, name := range recipients {
			if len(notReadMessages[name]) == 0 {
				notReadMessages[name] = make([]string, 0)
			}
			notReadMessages[name] = append(notReadMessages[name], messageForUsers)
		}
		beacon.Signal()
		addingMessageBeacpm.L.Unlock()
	}
}

func addMessageForRecipiants(message string, recipients []string, sec int, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(time.Duration(sec) * time.Second)
	addingMessageBeacpm.L.Lock()
	for _, value := range recipients {
		if len(notReadMessages[value]) == 0 {
			notReadMessages[value] = make([]string, 0)
		}
		notReadMessages[value] = append(notReadMessages[value], message)
	}
	beacon.Signal()
	addingMessageBeacpm.L.Unlock()
}

func sendMessageToRecipiants(beacon *sync.Cond) {
	beacon.L.Lock()
	for recipient, messages := range notReadMessages {
		for _, message := range messages {
			listener.SendMessage(message, recipient)
		}
	}
	beacon.L.Unlock()
	notReadMessages = make(map[string][]string)
}
