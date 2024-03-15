package main

import (
	"fmt"
	"log"
	"os"

	"github.com/sean9999/gorph"
)

func main() {
	g, err := gorph.NewGorph("../../testdata", "*", os.DirFS("../../testdata"))
	if err != nil {
		panic(err)
	}
	everything, err := g.Walk()
	if err != nil {
		panic(err)
	}
	fmt.Printf("WALK:\t%#v\n\n", everything)

	folders := g.Folders()
	fmt.Printf("FOLDERS:\t%#v\n\n", folders)

	evs, ers := g.Listen()

	for {
		select {
		case ev, ok := <-evs:
			//log.Println("event", ev, ok)
			if ok {
				log.Printf("\nnotify:\t%s\nevent:\t%s\npath:\t%s", ev.NotifyEvent.Op.String(), ev.Op, ev.Path)
				log.Println(g.WatchList())
				log.Println("**************")
			} else {
				log.Println("the channel is dead")
			}

		case er, ok := <-ers:
			log.Println("error", er, ok)
		}
	}

	//g.Close()

}
