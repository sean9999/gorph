package gorph

import (
	"encoding/json"

	"github.com/fsnotify/fsnotify"
)

type GorphOp uint8

const (
	UndefinedOp GorphOp = iota
	FsNotifyEvent
	FolderAdded
	FolderRemoved
	FolderRenamed
	FolderMoved
	FolderModified
	FolderUnknown
)

func (gop GorphOp) String() string {
	return []string{"UndefinedOp", "FsNotifyEvent", "FolderAdded", "FolderRemoved", "FolderRenamed", "FolderMoved", "FolderModified", "FolderUnknown"}[gop]
}

type GorphEvent struct {
	NotifyEvent *fsnotify.Event
	Op          GorphOp
	Path        string
}

func (gevent GorphEvent) String() string {
	m := map[string]any{
		"Op":          gevent.Op.String(),
		"Path":        gevent.Path,
		"NotifyEvent": gevent.NotifyEvent.String(),
	}
	j, _ := json.Marshal(m)
	return string(j)
}

func GorphEventFromNotifyEvent(g *gorph, fevent *fsnotify.Event) GorphEvent {

	gop := UndefinedOp
	shortPath := ""

	wasDir := func(longPath string) bool {
		return g.knownFolders[longPath]
	}

	shortPath = g.shortPath(fevent.Name)

	if IsDir(g.backer, shortPath) {
		switch fevent.Op {
		case fsnotify.Create:
			gop = FolderAdded
		case fsnotify.Remove:
			// @note: this should not happen, because events are fired after the file is removed
			gop = FolderRemoved
		case fsnotify.Rename:
			gop = FolderRenamed
		default:
			gop = FolderUnknown
		}
	} else if wasDir(shortPath) {
		switch fevent.Op {
		case fsnotify.Remove:
			gop = FolderRemoved
		case fsnotify.Rename:
			gop = FolderRenamed
		default:
			gop = FsNotifyEvent
		}
	} else {
		gop = FsNotifyEvent
	}

	if gop == UndefinedOp {
		panic("invalid GorphEvent")
	}
	return GorphEvent{
		NotifyEvent: fevent,
		Op:          gop,
		Path:        shortPath,
	}
}
