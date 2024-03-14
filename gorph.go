package gorph

import (
	"fmt"
	"io/fs"

	"github.com/bmatcuk/doublestar/v4"
)

// Gorph is a file system that can do globbing
type Gorph interface {
	fs.GlobFS
	Root() string
	Pattern() string
	Walk() ([]string, error)
}

// gorph implements Gorph
type gorph struct {
	root    string
	pattern string
	backer  fs.FS
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

func (g *gorph) Glob(pattern string) ([]string, error) {
	matches := []string{}
	var err error
	var fn fs.WalkDirFunc = func(path string, d fs.DirEntry, err error) error {
		itMatches, _ := doublestar.PathMatch(pattern, path)
		if itMatches && path != "." {
			matches = append(matches, path)
		}
		return err
	}
	fs.WalkDir(g.backer, ".", fn)
	return matches, err
}

func (g *gorph) Walk() ([]string, error) {
	paths := []string{}
	var fn fs.WalkDirFunc = func(path string, d fs.DirEntry, err error) error {
		paths = append(paths, path)
		return err
	}
	fs.WalkDir(g.backer, ".", fn)
	return paths, nil
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
	g := gorph{root: root, backer: back, pattern: pattern}
	return &g, nil
}
