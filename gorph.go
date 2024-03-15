package gorph

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/fsnotify/fsnotify"
)

// Gorph is a file system that can do globbing
type Gorph interface {
	fs.GlobFS
	Root() string
	Pattern() string
	Walk() ([]string, error)
	Folders() []string
	AddFolder(GorphEvent)
	RemoveFolder(GorphEvent)
	Children(path string) ([]string, error)
	Listen() (chan GorphEvent, chan error)
	Close() error
	WatchList() []string
}

// gorph implements Gorph
type gorph struct {
	root         string
	pattern      string
	backer       fs.FS
	events       chan GorphEvent
	Watcher      *fsnotify.Watcher
	knownFolders map[string]bool
}

func (g *gorph) Open(name string) (fs.File, error) {
	return g.backer.Open(name)
}

func (g *gorph) Pattern() string {
	return g.pattern
}

func (g *gorph) Root() string {
	return g.root
}

func (g *gorph) WatchList() []string {
	return g.Watcher.WatchList()
}

// walk the tree returning everything that matches pattern
func (g *gorph) Glob(pattern string) ([]string, error) {
	return doublestar.Glob(g.backer, pattern)
}

// return all child folders of parent
func (g *gorph) Children(parent string) ([]string, error) {
	var pattern string
	switch parent {
	case "":
		pattern = "**"
	case ".":
		pattern = "**"
	default:
		pattern = fmt.Sprintf("%s/**", parent)
	}

	matches := []string{}
	var err error
	var fn fs.WalkDirFunc = func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			itMatches, _ := doublestar.PathMatch(pattern, path)
			if itMatches && path != "." && path != parent {
				matches = append(matches, path)
			}
		}
		return err
	}
	fs.WalkDir(g.backer, ".", fn)
	return matches, err
}

// the file tree, but just folders
//
//	@note:	this does NOT filter using glob, because some folders may fail the glob, but files _in_ those folders would pass
//	@todo:	this could be optimised. There are some cases where we know we can safely omit the folder
func (g *gorph) Folders() []string {
	paths := []string{}
	var fn fs.WalkDirFunc = func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			paths = append(paths, path)
			g.knownFolders[path] = true
		}
		return err
	}
	fs.WalkDir(g.backer, ".", fn)
	return paths
}

func (g *gorph) AddFolder(gevent GorphEvent) {
	g.knownFolders[gevent.Path] = true
	g.Watcher.Add(gevent.NotifyEvent.Name)
}

func (g *gorph) RemoveFolder(gevent GorphEvent) {
	delete(g.knownFolders, gevent.Path)
	g.Watcher.Remove(gevent.NotifyEvent.Name)
}

func (g *gorph) Walk() ([]string, error) {
	paths := []string{}
	pattern := g.Pattern()

	var fn fs.WalkDirFunc = func(path string, d fs.DirEntry, err error) error {
		itMatches, _ := doublestar.PathMatch(pattern, path)
		if itMatches {
			paths = append(paths, path)
		}
		return err
	}
	fs.WalkDir(g.backer, ".", fn)
	return paths, nil
}

func (g *gorph) shortPath(longpath string) string {
	short, er := filepath.Rel(g.root, longpath)
	if er != nil {
		panic(er)
	}
	return short
}

func (g *gorph) longPath(shortPath string) string {
	long := filepath.Join(g.root, shortPath)
	return long
}

func IsDir(filesystem fs.FS, path string) bool {
	fyle, err := filesystem.Open(path)
	if err != nil {
		return false
	}
	stat, err := fyle.Stat()
	if err != nil {
		return false
	}
	return stat.IsDir()
}

func (g *gorph) Listen() (chan GorphEvent, chan error) {
	go func() {
		for {
			select {
			case fevent, ok := <-g.Watcher.Events:
				if !ok {
					return
				}
				gevent := GorphEventFromNotifyEvent(g, &fevent)
				switch gevent.Op {
				case FolderAdded:
					g.AddFolder(gevent)
				case FolderRemoved:
					g.RemoveFolder(gevent)
				}
				g.events <- gevent
			case err, ok := <-g.Watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()
	go g.Watcher.Add(g.root)
	return g.events, g.Watcher.Errors
}

func (g *gorph) Close() error {
	return g.Watcher.Close()
}

func NewGorph(root string, pattern string, back fs.FS) (Gorph, error) {
	self, err := back.Open(".")
	if err != nil {
		return nil, err
	}
	stat, err := self.Stat()
	if err != nil {
		return nil, err
	}
	if !stat.IsDir() {
		return nil, fmt.Errorf("%q is not a dir", root)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	events := make(chan GorphEvent)

	g := gorph{root: root, backer: back, pattern: pattern, Watcher: watcher, events: events, knownFolders: map[string]bool{}}

	for _, shortFolder := range g.Folders() {
		g.knownFolders[shortFolder] = true
		longFolder := g.longPath(shortFolder)
		g.Watcher.Add(longFolder)
	}

	return &g, nil
}
