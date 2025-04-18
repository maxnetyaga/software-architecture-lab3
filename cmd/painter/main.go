package main

import (
	"log"
	"net/http"

	"github.com/maxnetyaga/software-architecture-lab3/painter"
	"github.com/maxnetyaga/software-architecture-lab3/painter/lang"
	"github.com/maxnetyaga/software-architecture-lab3/ui"
)

func main() {
	var (
		pv ui.Visualizer

		opLoop painter.Loop
		parser lang.Parser
	)

	pv.Debug = true
	pv.Title = "Simple painter (Variant 5)"

	pv.OnScreenReady = opLoop.Start
	opLoop.Receiver = &pv

	opLoop.Done = make(chan struct{})

	go func() {
		http.Handle("/", lang.HttpHandler(&opLoop, &parser))
		log.Fatal(http.ListenAndServe("localhost:17000", nil))
	}()

	pv.Main()
	opLoop.StopAndWait()
}
