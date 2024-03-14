package gorph_test

import (
	"io/fs"
	"os"
	"reflect"
	"testing"

	"github.com/sean9999/gorph"
)

func TestNewGorph_Root(t *testing.T) {
	type args struct {
		root    string
		back    fs.FS
		pattern string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "self", args: args{root: ".", back: os.DirFS("."), pattern: "*"}, want: ".", wantErr: false},
		{name: "testdata", args: args{root: "testdata", back: os.DirFS("testdata"), pattern: "*"}, want: "testdata", wantErr: false},
		{name: "x", args: args{root: "x", back: os.DirFS("x"), pattern: "*"}, want: "testdata", wantErr: true},
		{name: "go.mod", args: args{root: "go.mod", back: os.DirFS("go.mod"), pattern: "*"}, want: "go.mod", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotGorph, err := gorph.NewGorph(tt.args.root, tt.args.pattern, tt.args.back)
			if err == nil {
				got := gotGorph.Root()
				if got != tt.want {
					t.Errorf("NewGorph() = %v, want %v", got, tt.want)
				}
			} else {
				if !tt.wantErr {
					t.Errorf("NewGorph() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestNewGorph_Walk(t *testing.T) {
	var VIDEOS = []string{".", "VID_1.mp4", "VID_2.mp4", "waiting-for-mommy.mov"}
	var DOWNLOADS = []string{".", "node", "node/node-v20.11.0-linux-x64.tar.xz"}

	type args struct {
		root    string
		back    fs.FS
		pattern string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{name: "Videos", args: args{root: "testdata/Videos", back: os.DirFS("testdata/Videos"), pattern: "*"}, want: VIDEOS, wantErr: false},
		{name: "Downloads", args: args{root: "testdata/Downloads", back: os.DirFS("testdata/Downloads"), pattern: "*"}, want: DOWNLOADS, wantErr: false},
		{name: "folder that doesn't exist", args: args{root: "x", back: os.DirFS("x"), pattern: "*"}, want: []string{""}, wantErr: true},
		{name: "go.mod", args: args{root: "go.mod", back: os.DirFS("go.mod"), pattern: "*"}, want: []string{""}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotGorph, err := gorph.NewGorph(tt.args.root, tt.args.pattern, tt.args.back)
			if err == nil {
				got, _ := gotGorph.Walk()
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("NewGorph() = %v, want %v", got, tt.want)
				}
			} else {
				if !tt.wantErr {
					t.Errorf("NewGorph() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestNewGorph_Glob(t *testing.T) {
	var VIDEOS = []string{"VID_1.mp4", "VID_2.mp4", "waiting-for-mommy.mov"}
	var JUST_NODE = []string{"node"}
	var DOWNLOADS = []string{"node", "node/node-v20.11.0-linux-x64.tar.xz"}
	var MOMMY = []string{"Documents/the-mommy-book.txt", "Pictures/mommy-and-me.jpeg", "Videos/waiting-for-mommy.mov"}

	type args struct {
		root    string
		back    fs.FS
		pattern string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{name: "3 videos", args: args{root: "testdata/Videos", back: os.DirFS("testdata/Videos"), pattern: "*"}, want: VIDEOS, wantErr: false},
		{name: "node", args: args{root: "testdata/Downloads", back: os.DirFS("testdata/Downloads"), pattern: "*"}, want: JUST_NODE, wantErr: false},
		{name: "Downloads", args: args{root: "testdata/Downloads", back: os.DirFS("testdata/Downloads"), pattern: "**"}, want: DOWNLOADS, wantErr: false},
		{name: "mommy files", args: args{root: "testdata", back: os.DirFS("testdata"), pattern: "**/*mommy*"}, want: MOMMY, wantErr: false},
		{name: "folder that doesn't exist", args: args{root: "x", back: os.DirFS("x"), pattern: "*"}, want: []string{""}, wantErr: true},
		{name: "go.mod", args: args{root: "go.mod", back: os.DirFS("go.mod"), pattern: "*"}, want: []string{""}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotGorph, err := gorph.NewGorph(tt.args.root, tt.args.pattern, tt.args.back)
			if err == nil {
				got, _ := gotGorph.Glob(tt.args.pattern)
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("got %v, wanted %v", got, tt.want)
				}
			} else {
				if !tt.wantErr {
					t.Errorf("got error = %v, but wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}
