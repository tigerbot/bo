package main

// This file is largely based on the go-bindata-assetfs library from elazarl. We copied and
// modified his code mostly just to prevent the http.FileServer from reporting 500 internal
// server error when it should really be reporting 404 not found.

import (
	"bytes"
	"net/http"
	"os"
	"path"
	"time"
)

var (
	defaultFileTimestamp = time.Now()
)

// FakeFile implements http.File interface for a given path and size. It will error on any
// function that doesn't worth for both files and directories, so the struct that embed this
// don't have to worry about defining anything that they don't actually need.
type FakeFile struct {
	name      string
	dir       bool
	size      int64
	timestamp time.Time
}

func (f *FakeFile) Sys() interface{}   { return nil }
func (f *FakeFile) Name() string       { return f.name }
func (f *FakeFile) IsDir() bool        { return f.dir }
func (f *FakeFile) Size() int64        { return f.size }
func (f *FakeFile) ModTime() time.Time { return f.timestamp }
func (f *FakeFile) Mode() os.FileMode {
	mode := os.FileMode(0644)
	if f.dir {
		return mode | os.ModeDir
	}
	return mode
}

func (f *FakeFile) Read(b []byte) (int, error) {
	return 0, os.ErrInvalid
}
func (f *FakeFile) Seek(offset int64, whence int) (int64, error) {
	return 0, os.ErrInvalid
}
func (f *FakeFile) Close() error {
	return nil
}
func (f *FakeFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, os.ErrInvalid
}
func (f *FakeFile) Stat() (os.FileInfo, error) {
	return f, nil
}

// AssetFile implements http.File interface for a no-directory file with content
type AssetFile struct {
	FakeFile
	*bytes.Reader
}

func (f *AssetFile) Read(b []byte) (int, error) {
	return f.Reader.Read(b)
}
func (f *AssetFile) Seek(offset int64, whence int) (int64, error) {
	return f.Reader.Seek(offset, whence)
}

func NewAssetFile(name string, content []byte, timestamp time.Time) *AssetFile {
	if timestamp.IsZero() {
		timestamp = defaultFileTimestamp
	}
	name = path.Base(name)
	return &AssetFile{
		FakeFile{name, false, int64(len(content)), timestamp},
		bytes.NewReader(content),
	}
}

// AssetDirectory implements http.File interface for a directory
type AssetDirectory struct {
	FakeFile
	ChildrenRead int
	Children     []os.FileInfo
}

func NewAssetDirectory(name string, children []string) *AssetDirectory {
	fileinfos := make([]os.FileInfo, 0, len(children))
	for _, child := range children {
		childInfo := &FakeFile{child, false, 0, defaultFileTimestamp}
		fullpath := path.Join(name, child)
		if info, err := AssetInfo(fullpath); err == nil {
			childInfo.timestamp = info.ModTime()
			childInfo.dir = info.IsDir()
		}
		fileinfos = append(fileinfos, childInfo)
	}
	return &AssetDirectory{
		FakeFile:     FakeFile{name, true, 0, defaultFileTimestamp},
		ChildrenRead: 0,
		Children:     fileinfos,
	}
}

func (f *AssetDirectory) Readdir(count int) ([]os.FileInfo, error) {
	if count <= 0 {
		return f.Children, nil
	}
	if f.ChildrenRead+count > len(f.Children) {
		count = len(f.Children) - f.ChildrenRead
	}
	rv := f.Children[f.ChildrenRead : f.ChildrenRead+count]
	f.ChildrenRead += count
	return rv, nil
}

// AssetFS implements http.FileSystem, allowing the files embedded in this package to be
// served using the net/http package.
type AssetFS struct{}

func (AssetFS) Open(name string) (http.File, error) {
	if len(name) > 0 && name[0] == '/' {
		name = name[1:]
	}
	if b, err := Asset(name); err == nil {
		timestamp := defaultFileTimestamp
		if info, err := AssetInfo(name); err == nil {
			timestamp = info.ModTime()
		}
		return NewAssetFile(name, b, timestamp), nil
	}
	if children, err := AssetDir(name); err == nil {
		return NewAssetDirectory(name, children), nil
	} else {
		return nil, os.ErrNotExist
	}
}
