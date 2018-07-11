package utils

import (
	"fmt"
	"github.com/tprasadtp/wsl-git/version"
	"os"
	"os/exec"
)

//CheckwslExists ... Check if wsl.exe is present
func CheckwslExists() {
	path, err := exec.LookPath("wsl")
	if err != nil {
		fmt.Printf("Didn't find 'wsl' executable.\nMake Sure you have WSL enabled and a distro installed?\n")
		os.Exit(7)
	} else {
		fmt.Printf("'wsl' executable is in '%s'\n", path)
	}
}

//Win2Wsl ... Convert path to wsl path
func Win2Wsl(path string) (string, error) {
	var wslpathreturn []byte
	// https://docs.microsoft.com/en-us/windows/wsl/release-notes#build-17046
	/*
			  wslpath usage:
		    -a    force result to absolute path format
		    -u    translate from a Windows path to a WSL path (default)
		    -w    translate from a WSL path to a Windows path
				-m    translate from a WSL path to a Windows path, with ‘/’ instead of ‘\\’
	*/
	wslpathreturn, err := exec.Command("wsl", "wslpath", "-u", path).Output()
	if err != nil {
		fmt.Printf("Some error occurre while converting path to wsl: %v\n", err)
		os.Exit(2)
	}
	return string(wslpathreturn), err
}

//PrintError ... Prints error message and exists with error code passed.
func PrintError(msg string, exitcode int) {
	//fmt.Println(colorstring.Color("[red]" + err.Error()))
	fmt.Printf(msg)
	os.Exit(exitcode)
}

//Usage ... Help Text
func Usage() {
	fmt.Printf("WSL-GIT - A Bridge between git installed in WSL and Windows\n")
	fmt.Printf("----------------\n\n")
	fmt.Printf("Version - %s\n", version.VERSION)
	fmt.Printf("Usage:")
	fmt.Printf("wsl-git <your git commands>\n")
	fmt.Printf("----------------\n\n")
	fmt.Printf("--wsl-git-version     Display Version info.\n")
	fmt.Printf("--wsl-git-help        Display this message.\n")
	fmt.Printf("--wsl-git-print-args  Display all the argumets passed to the program.\n")
	fmt.Printf("Remember if no arguments are passed, it will display git's help, as git would do.\n")
	fmt.Printf("----------------\n\n")
}
