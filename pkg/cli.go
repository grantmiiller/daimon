package daimon

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/chzyer/readline"
)

var fm FM

const (
	// Directory
	envName = "DAIMON_DIR"
)

type command struct {
	Description string
	Func        func([]string) error
}

// Big ol list of commands and their corresponding functions and descriptions
var commandMap = map[string]command{
	"n": command{
		Description: "Creates a new note. Usage: n <NOTE_NAME> <NOTE_CONTENT>",
		Func:        newNote,
	},
	"e": command{
		Description: "Opens a note in editor. Usage: e <NOTE_NAME>",
		Func:        editNote,
	},
	"p": command{
		Description: "Prints a note. Usage: p <NOTE_NAME>",
		Func:        printNote,
	},
	"mv": command{
		Description: "Renames a note. Usage: mv <NOTE_NAME> <NEW_NAME>",
		Func:        renameNote,
	},
	"l": command{
		Description: "Lists notes, optionally in a project. Usage: l [PROJECT_NAME]",
		Func:        listNotes,
	},
	"d": command{
		Description: "Deletes a note. Usage: d <NOTE_NAME>",
		Func:        deleteNote,
	},
	"np": command{
		Description: "Create a new project. Usage: np <PROJECT_NAME>",
		Func:        createProject,
	},
	"lp": command{
		Description: "Lists projects, by default in root directory or provide project name to list subprojects. Usage: lp [PROJECT_NAME]",
		Func:        listProjects,
	},
	"dp": command{
		Description: "Deletes a project. Usage: lp <PROJECT_NAME>",
		Func:        deleteProject,
	},
	"la": command{
		Description: "Lists all notes and projects. Usage: la",
		Func:        listAll,
	},
	"lan": command{
		Description: "Lists all notes in root, projects, and subprojects. Usage: la",
		Func:        listAllNotes,
	},
	"lap": command{
		Description: "Lists all projects and subprojects. Usage: la",
		Func:        listAllProjects,
	},
	"line": command{
		Description: "Says a line. Usage: line",
		Func:        sayLine,
	},
	"c": command{
		Description: "Clears the terminal. Usage: c",
		Func:        clearTerminal,
	},
	"q": command{
		Description: "Quit the program. Usage: q",
		Func:        quit,
	},
	"exit": command{
		Description: "Quit the program. Usage: exit",
		Func:        quit,
	},
	"quit": command{
		Description: "Quit the program. Usage: quit",
		Func:        quit,
	},
}

func help() {
  intro()
	fmt.Println("SYNOPSIS:")
	fmt.Printf("\t daimon [COMMAND] [ARGS]\n\n")
	fmt.Printf("\t Calling daimon without any arguments starts it in interactive mode\n\n")
	fmt.Println("Commands:")
	fmt.Printf("h, help\t - Displays this help message. Usage: h|help\n\n")
	for k, v := range commandMap {
		fmt.Printf("%s\t- %s\n\n", k, v.Description)
	}
}

func intro() {
	fmt.Println("daimon: A handy little note-taking assistant")
	fmt.Println("============================================")
}

func newNote(args []string) error {
	// new note
	if len(args) < 2 {
		return errors.New("not enough arguments for note")
	}

	if err := QuickNote(fm, args[0], strings.Join(args[1:len(args)], " ")); err != nil {
		return fmt.Errorf("could not create note: %s", err)
	}
	return nil
}

func editNote(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("need to provide name of note")
	}

	err := EditNote(fm, args[0])

	if err != nil {
		return fmt.Errorf("could not open note: %s", err)
	}
	return nil
}

func printNote(args []string) error {
	// print note
	if len(args) < 1 {
		return fmt.Errorf("not enough arguments for note")
	}

	note, err := GetNote(fm, args[0])

	if err != nil {
		return fmt.Errorf("could not read note: %s", err)
	}

	// Add an empty buffer line before printing note
	fmt.Println("")
	fmt.Printf("%s\n", note)
	return nil
}

func renameNote(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("not enough arguments for renaming")
	}
	return fm.Rename(args[0], args[1])
}

func listNotes(args []string) error {
	// list notes
	p := "."
	if len(args) > 0 {
		p = args[0]
	}
	notes, err := ListNotes(fm, p)
	if err != nil {
		return fmt.Errorf("could not list notes: %s", err)
	}
	for _, note := range notes {
		fmt.Println(note)
	}
	return nil
}

func deleteNote(args []string) error {
	// Delete note
	if len(args) < 1 {
		return errors.New("not enough arguments for note")
	}
	err := RemoveNote(fm, args[0])
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func createProject(args []string) error {
	// Creates a new project
	if len(args) < 1 {
		return errors.New("not enough arguments for note")
	}

	if err := NewProject(fm, args[0]); err != nil {
		return fmt.Errorf("could not create project: %s", err)
	}
	return nil
}

func listProjects(args []string) error {
	// List projects
	p := "."
	if len(args) > 0 {
		p = args[0]
	}
	projects, err := ListProjects(fm, p)
	if err != nil {
		return fmt.Errorf("could not list projects: %s", err.Error())
	}
	if len(projects) == 0 {
		fmt.Println("<No Projects>")
	}
	for _, project := range projects {
		fmt.Println(project)
	}
	return nil
}

func deleteProject(args []string) error {
	// Delete Project
	if len(args) < 1 {
		return errors.New("not enough arguments for deleting project")
	}
	return RemoveProject(fm, args[0])
}

func listAll(_ []string) error {
	// List all
	projects, err := ListAll(fm)
	if err != nil {
		return fmt.Errorf("could not list projects: %s", err.Error())
	}
	if len(projects) == 0 {
		fmt.Println("<No Projects>")
	}
	for _, project := range projects {
		fmt.Println(project)
	}
	return nil
}

func listAllNotes(_ []string) error {
	// List all
	projects, err := ListAllNotes(fm)
	if err != nil {
		return fmt.Errorf("could not list projects: %s", err.Error())
	}
	if len(projects) == 0 {
		fmt.Println("<No Projects>")
	}
	for _, project := range projects {
		fmt.Println(project)
	}
	return nil
}

func listAllProjects(_ []string) error {
	// List all projects
	projects, err := ListAllProjects(fm)
	if err != nil {
		return fmt.Errorf("could not list projects: %s", err.Error())
	}
	if len(projects) == 0 {
		fmt.Println("<No Projects>")
	}
	for _, project := range projects {
		fmt.Println(project)
	}
	return nil
}

func sayLine(_ []string) error {
	lines := []string{
		"Boopy stoopy",
		"Shloop",
	}
	rand.Seed(time.Now().Unix())
	n := rand.Int() % len(lines)
	fmt.Println(lines[n])
	return nil
}

func clearTerminal(_ []string) error {
	fmt.Print("\033[H\033[2J")
	return nil
}

func quit(_ []string) error {
	os.Exit(0)
	return nil
}

func getBaseDir() string {
	dir, hasEnv := os.LookupEnv(envName)
	if hasEnv {
		if !strings.HasSuffix(dir, "/") {
			dir = dir + "/"
		}
		return dir
	}
	fmt.Fprintf(os.Stderr, "ERROR: %s env is not set\n", envName)
	os.Exit(1)
	return ""
}

func setup() {
	fm = NewFM(getBaseDir())
}

func getCmdPrompt(text string) (string, []string) {
	text = strings.TrimSuffix(text, "\n")
	line := strings.Split(text, " ")
	command := line[0]
	var args []string
	if len(line) > 1 {
		args = line[1:len(line)]
	}
	return command, args
}

func autocompleteListProjects(fm FileIO, dir string) func(line string) []string {
	return func(_ string) []string {
		dirs, err := ListAllProjects(fm)
		if err != nil {
			panic("Could not list directories")
		}
		return dirs
	}
}

func autocompleteListAll(fm FileIO, dir string) func(line string) []string {
	return func(_ string) []string {
		dirs, err := ListAll(fm)
		if err != nil {
			panic("Could not list all")
		}
		return dirs
	}
}

func autocompleteListAllNotes(fm FileIO, dir string) func(line string) []string {
	return func(_ string) []string {
		dirs, err := ListAllNotes(fm)
		if err != nil {
			panic("Could not list all")
		}
		return dirs
	}
}

func getCompleter(fm FileIO) *readline.PrefixCompleter {
	return readline.NewPrefixCompleter(
		readline.PcItem("n",
			readline.PcItemDynamic(autocompleteListAll(fm, ".")),
		),
		readline.PcItem("e",
			readline.PcItemDynamic(autocompleteListAll(fm, ".")),
		),
		readline.PcItem("p",
			readline.PcItemDynamic(autocompleteListAllNotes(fm, ".")),
		),
		readline.PcItem("mv",
			readline.PcItemDynamic(autocompleteListAllNotes(fm, "."),
				readline.PcItemDynamic(autocompleteListProjects(fm, ".")),
			),
		),
		readline.PcItem("l",
			readline.PcItemDynamic(autocompleteListProjects(fm, ".")),
		),
		readline.PcItem("d",
			readline.PcItemDynamic(autocompleteListAllNotes(fm, ".")),
		),
		readline.PcItem("np",
			readline.PcItemDynamic(autocompleteListProjects(fm, ".")),
		),
		readline.PcItem("lp",
			readline.PcItemDynamic(autocompleteListProjects(fm, ".")),
		),
		readline.PcItem("dp",
			readline.PcItemDynamic(autocompleteListProjects(fm, ".")),
		),
		readline.PcItem("la"),
		readline.PcItem("lan"),
		readline.PcItem("lap"),
		readline.PcItem("line"),
		readline.PcItem("q"),
		readline.PcItem("quit"),
		readline.PcItem("exit"),
	)
}

func runCommand(command string, args []string) error {
	if command == "h" || command == "help" {
		help()
	} else if cmd, ok := commandMap[command]; ok {
		cmd.Func(args)
	} else {
		return errors.New("I don't understand")
	}
	return nil
}

func interactiveLoop() {
	l, err := readline.NewEx(&readline.Config{
		Prompt:            "\033[1;95m»\033[0m ",
		HistoryFile:       "/tmp/readline.tmp",
		InterruptPrompt:   "^C",
		EOFPrompt:         "exit",
		AutoComplete:      getCompleter(fm),
		HistorySearchFold: true,
	})
	if err != nil {
		panic(err)
	}
	for {
		text, err := l.Readline()
		command, args := getCmdPrompt(text)
		err = runCommand(command, args)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
		}
	}
}

// Init is the entry point to run daimon
func Init() {
	setup()
	args := os.Args[1:]
	if len(args) < 1 {
    intro()
		interactiveLoop()
	}
	err := runCommand(args[0], args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
}
