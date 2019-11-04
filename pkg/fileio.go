package daimon

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
)

type FileIO interface {
	GetFullPath(string) string
	GetFullNotePath(string) string
	OpenFile(string) (*os.File, error)
	Open(string) (*os.File, error)
	// Different from os MkdirAll since we handle the file perm
	MkdirAll(string) error
	Remove(string) error
	RemoveAll(string) error
	WriteString(string, string) (int, error)
}

type FM struct {
	pathBase string
	mode     int
}

func NewFM(base string) FM {
	return FM{pathBase: base, mode: int(0750)}
}

func hasIllegalPath(p string) bool {
	re := regexp.MustCompile(`^\W`)
	if re.MatchString(p) {
		return true
	}
	return false
}

func (fm FM) GetFullPath(p string) string {
	// A single . is allowed
	if p != "." && hasIllegalPath(p) {
		fmt.Println(p)
		panic("Path starts with an invalid character. Panicking due to expected shenanigans.")
	}
	return filepath.Join(fm.pathBase + path.Clean(p))
}

func (fm FM) GetFullNotePath(p string) string {
	// Notes always have to start with allowed characters. Dots are not allowed
	if hasIllegalPath(path.Base(p)) {
		panic("Path starts with an invalid character. Panicking due to expected shenanigans.")
	}
	return fm.GetFullPath(p) + ".md"
}

func (fm FM) MkdirAll(p string) error {
	return os.MkdirAll(fm.GetFullPath(p), os.FileMode(fm.mode))
}

func (fm FM) Open(p string) (*os.File, error) {
	return os.Open(fm.GetFullNotePath(p))
}

func (fm FM) OpenFile(p string) (*os.File, error) {
	dir := path.Dir(p)
	err := fm.MkdirAll(dir)
	if err != nil {
		return nil, err
	}
	return os.OpenFile(fm.GetFullNotePath(p), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0750)
}

// Rename renames from a file from the old location to new location,
// and creating new projects in the process for the new location
func (fm FM) Rename(ol string, nl string) error {
	dir := path.Dir(nl)
	err := fm.MkdirAll(dir)
	if err != nil {
		return err
	}
	return os.Rename(fm.GetFullNotePath(ol), fm.GetFullNotePath(nl))
}

func (fm FM) Remove(p string) error {
	return os.Remove(fm.GetFullNotePath(p))
}

func (fm FM) RemoveAll(p string) error {
	return os.RemoveAll(fm.GetFullPath(p))
}

func (fm FM) WriteString(p, s string) (int, error) {
	fp, err := fm.OpenFile(p)
	if err != nil {
		return 0, err
	}

	defer fp.Close()

	n, err := fp.WriteString(s)
	return n, err
}
