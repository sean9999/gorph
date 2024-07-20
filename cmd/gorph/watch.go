package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/sean9999/go-flargs"
	"github.com/sean9999/gorph"
)

var ErrGorph = errors.New("gorph")
var ErrFatal = fmt.Errorf("%w: fatal", ErrGorph)

type outputOptions struct {
	quoted    bool
	base64    bool
	json      bool
	delimiter string
}

func quoted(str string) string {
	return fmt.Sprintf("'%s'", str)
}

func output(gev gorph.GorphEvent, opts outputOptions) string {
	ev := gev.NotifyEvent.Op.String()
	path := gev.Path
	typ := gev.Op.String()

	switch {
	case opts.base64 == true && opts.json == false:
		if opts.quoted {
			ev = quoted(ev)
			path = quoted(ev)
			typ = quoted(typ)
		}
		nakedString := fmt.Sprintf("%s\n%s\n%s\n", ev, path, typ)
		// ex: V1JJVEUKLmJhc2hyYwpGc05vdGlmeUV2ZW50Cg==
		return base64.StdEncoding.EncodeToString([]byte(nakedString))
	case opts.base64 == true && opts.json == true:
		m := map[string]string{
			"event": ev,
			"path":  path,
			"type":  typ,
		}
		j, _ := json.Marshal(m)
		r := base64.StdEncoding.EncodeToString(j)
		if opts.quoted {
			//	ex: 'eyJldmVudCI6IldSSVRFIiwicGF0aCI6Ii5iYXNocmMiLCJ0eXBlIjoiRnNOb3RpZnlFdmVudCJ9'
			r = quoted(r)
		}
		// ex: eyJldmVudCI6IkNITU9EIiwicGF0aCI6Ii5iYXNocmMiLCJ0eXBlIjoiRnNOb3RpZnlFdmVudCJ9
		return r
	case opts.base64 == false && opts.json == true:
		m := map[string]string{
			"event": ev,
			"path":  path,
			"type":  typ,
		}
		j, _ := json.Marshal(m)
		r := string(j)
		if opts.quoted {
			// ex: '{"event":"WRITE","path":".bashrc","type":"FsNotifyEvent"}'
			r = quoted(r)
		}
		//	ex: {"event":"WRITE","path":".bashrc","type":"FsNotifyEvent"}
		return r
	default:
		if opts.quoted {
			//	ex: 'MODIFY'	'.bashrc'	'inotifyEvent'
			ev = quoted(ev)
			path = quoted(ev)
			typ = quoted(typ)
		}
		//	ex:	MODIFY	.bashrc	inotifyEvent
		return fmt.Sprintf("%s%s%s%s%s", ev, opts.delimiter, path, opts.delimiter, typ)
	}
}

func watch(env *flargs.Environment, root string, pat string, opts outputOptions) {
	g, err := gorph.NewGorph(root, pat, os.DirFS(root))
	if err != nil {
		panic(err)
	}
	events, errs := g.Listen()
	timeToDie := false
	for !timeToDie {
		select {
		case ev, ok := <-events:
			if ok {
				fmt.Fprintln(env.OutputStream, output(ev, opts))
				if ev.Op == gorph.TimeToDie {
					timeToDie = true
				}
			} else {
				fmt.Fprintln(env.ErrorStream, "the channel is dead")
			}

		case err, ok := <-errs:
			if errors.Is(err, ErrFatal) {
				timeToDie = true
			}
			fmt.Fprintln(env.ErrorStream, "error", err, ok)
		}
	}
	g.Close()
}
