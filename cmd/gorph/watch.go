package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/sean9999/gorph"
)

var ErrGorph = errors.New("gorph")
var ErrFatal = fmt.Errorf("%w: fatal", ErrGorph)

type outputOptions struct {
	quotes    bool
	base64    bool
	json      bool
	delimiter string
}

func output(gev gorph.GorphEvent, opts outputOptions) string {

	ev := gev.NotifyEvent.Op.String()
	path := gev.Path
	typ := gev.Op.String()

	switch {
	case opts.base64 == true && opts.json == false:
		nakedString := fmt.Sprintf("%s\n%s\n%s\n", ev, path, typ)
		return base64.StdEncoding.EncodeToString([]byte(nakedString))
	case opts.base64 == true && opts.json == true:
		m := map[string]string{
			"event": ev,
			"path":  path,
			"type":  typ,
		}
		j, _ := json.Marshal(m)
		return base64.StdEncoding.EncodeToString(j)
	case opts.base64 == false && opts.json == true:
		m := map[string]string{
			"event": ev,
			"path":  path,
			"type":  typ,
		}
		j, _ := json.Marshal(m)
		return string(j)
	case opts.quotes == true:
		ev = fmt.Sprintf("'%s'", ev)
		path = fmt.Sprintf("'%s'", path)
		typ = fmt.Sprintf("'%s'", typ)
		return fmt.Sprintf("%s%s%s%s%s", ev, opts.delimiter, path, opts.delimiter, typ)
	default:
		return fmt.Sprintf("%s%s%s%s%s", ev, opts.delimiter, path, opts.delimiter, typ)
	}

}

func watch() {
	g, err := gorph.NewGorph("testdata", "*", os.DirFS("testdata"))
	if err != nil {
		panic(err)
	}

	opts := outputOptions{
		delimiter: "\t",
	}

	evs, ers := g.Listen()

	timeToDie := false

	for !timeToDie {
		select {
		case ev, ok := <-evs:
			if ok {
				fmt.Println(output(ev, opts))
				if ev.Op == gorph.TimeToDie {
					timeToDie = true
				}
			} else {
				//log.Println("the channel is dead")
				fmt.Fprintln(os.Stderr, "the channel is dead")
			}

		case er, ok := <-ers:
			if errors.Is(er, ErrFatal) {
				timeToDie = true
			}
			log.Println("error", er, ok)
		}
	}

	g.Close()

}
