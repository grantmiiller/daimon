package daimon

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

// GetNote shows note at a given path (path is turned into full path)
// MAYBE: set up to use a writer instead of ioutil.ReadAll
func GetNote(fm FileIO, np string) ([]byte, error) {
	fp, err := fm.Open(np)
	if err != nil {
		return []byte{}, err
	}
	defer fp.Close()

	note, err := ioutil.ReadAll(fp)

	if err != nil {
		return []byte{}, err
	}
	return note, nil
}

func QuickNote(fm FileIO, name string, s string) error {
	if _, err := fm.WriteString(name, s+"\n"); err != nil {
		return err
	}
	return nil
}

func EditNote(fm FileIO, name string) error {
	editor, hasEnv := os.LookupEnv("EDITOR")
	if (!hasEnv) {
		editor = "vim"
	}
	cmd := exec.Command(editor, fm.GetFullNotePath(name))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func NewProject(fm FileIO, name string) error {
	return fm.MkdirAll(name)
}

func RemoveNote(fm FileIO, name string) error {
	return fm.Remove(name)
}

func RemoveProject(fm FileIO, name string) error {
	return fm.RemoveAll(name)
}

func ListNotes(fm FileIO, p string) ([]string, error) {
	files, err := ioutil.ReadDir(fm.GetFullPath(p))
	if err != nil {
		return []string{}, err
	}
	var f []string
	for _, file := range files {
		name := file.Name()
		if strings.HasSuffix(name, ".md") {
			f = append(f, strings.TrimSuffix(name, ".md"))
		}
	}
	return f, nil
}

func ListProjects(fm FileIO, p string) ([]string, error) {
	files, err := ioutil.ReadDir(fm.GetFullPath(p))
	if err != nil {
		return []string{}, err
	}
	var f []string
	for _, file := range files {
		mode := file.Mode()
		if mode.IsDir() {
			f = append(f, file.Name())
		}
	}

	return f, nil
}

// This actually does all the work of ListAllProjects
func listAllProjectsWorker(fm FileIO, d string) ([]string, error) {
	var rp []string
	projects, err := ListProjects(fm, d)
	if err != nil {
		return []string{}, err
	}
	// We haved to reset basepath when it's a '.' becuase I have logic
	// to stop shenanigans when things start with .
	basepath := d + "/"
	if d == "." {
		basepath = ""
	}
	for _, p := range projects {
		rp = append(rp, p+"/")
		// Get sub-projects
		subprojects, err := listAllProjectsWorker(fm, basepath+p)
		if err != nil {
			return []string{}, err
		}
		for _, subp := range subprojects {
			rp = append(rp, p+"/"+subp)
		}
	}
	return rp, nil
}

// This actually does all the work of ListAllProjects
func listAllNotesWorker(fm FileIO, d string) ([]string, error) {
	var rp []string
	files, err := ioutil.ReadDir(fm.GetFullPath(d))
	if err != nil {
		return []string{}, err
	}

	// We haved to reset basepath when it's a '.' becuase I have logic
	// to stop shenanigans when things start with .
	basepath := d + "/"
	if d == "." {
		basepath = ""
	}

	for _, file := range files {
		mode := file.Mode()
		name := file.Name()
		if strings.HasSuffix(name, ".md") {
			rp = append(rp, strings.TrimSuffix(name, ".md"))
		}
		if mode.IsDir() {
			subprojects, err := listAllNotesWorker(fm, basepath+name)
			if err != nil {
				return []string{}, err
			}
			for _, subp := range subprojects {
				rp = append(rp, name+"/"+subp)
			}

		}
	}
	return rp, nil
}

// This actually does all the work of ListAllProjects
func listAllWorker(fm FileIO, d string) ([]string, error) {
	var rp []string
	files, err := ioutil.ReadDir(fm.GetFullPath(d))
	if err != nil {
		return []string{}, err
	}

	// We haved to reset basepath when it's a '.' becuase I have logic
	// to stop shenanigans when things start with .
	basepath := d + "/"
	if d == "." {
		basepath = ""
	}

	for _, file := range files {
		mode := file.Mode()
		name := file.Name()
		if strings.HasSuffix(name, ".md") {
			rp = append(rp, strings.TrimSuffix(name, ".md"))
		}
		if mode.IsDir() {
			rp = append(rp, name+"/")

			subprojects, err := listAllWorker(fm, basepath+name)
			if err != nil {
				return []string{}, err
			}
			for _, subp := range subprojects {
				rp = append(rp, name+"/"+subp)
			}

		}
	}
	return rp, nil
}

// ListAllProjects tries to grab all the project names. May
// be a tad heavy handed, but easiest thing to do for autocomplete
func ListAllProjects(fm FileIO) ([]string, error) {
	// Get initial list of projects
	return listAllProjectsWorker(fm, ".")
}

// ListAllNotes tries to grab all the note names. May
// be a tad heavy handed, but easiest thing to do for autocomplete
func ListAllNotes(fm FileIO) ([]string, error) {
	// Get initial list of projects
	return listAllNotesWorker(fm, ".")
}

// ListAll returns all pojects and notes. Great for autocomplete
func ListAll(fm FileIO) ([]string, error) {
	return listAllWorker(fm, ".")
}
