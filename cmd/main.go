package main

import (
	"bivrost"
	"log"
	"time"
)

type event struct {
	after   int
	message string
}

func main() {

	log.Println("start")

	s, in, out := bivrost.Init()
	go s.Serve()

	var events = []event{
		{1, "one second"},
		{2, "two seconds"},
		{3, "three seconds"},
	}

	go func() {
		for _, event := range events {
			in <- &bivrost.Event{
				When:   time.Now().Add(time.Duration(event.after) * time.Second),
				Entity: interface{}(event.message),
			}
		}
	}()

	go func() {
		for event := range out {
			log.Println(event)
		}
	}()

	time.Sleep(4 * time.Second)
	s.Cancel()
	time.Sleep(1 * time.Millisecond)
}
