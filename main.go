package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/tprasadtp/wsl-git/utils"
	"github.com/tprasadtp/wsl-git/version"
)

func main() {
	var (
		arguments []string
	)

	//Get all the arguments
	arguments = append(os.Args[1:])

	// Edge case: no arguments
	if len(arguments) < 1 {
		arguments = append([]string{"help"}, arguments...)
	}

	//
	if strings.ToLower(arguments[0]) == "--wsl-git-version" {
		fmt.Printf("\nWSL-GIT Version - %s\n\n", version.VERSION)
		return
	}
	// help
	if strings.ToLower(arguments[0]) == "--wsl-git-help" {
		utils.Usage()
		return
	}

	if strings.ToLower(arguments[0]) == "--wsl-git-print-args" {
		fmt.Printf("Arguments passed to %s are\n%v", os.Args[0], arguments)
		return
	}

	//check if not on windows
	if runtime.GOOS != "windows" {
		utils.PrintError("Not running on Windows.\n", 1)
	}

	// This is just a hack which should work most of the times
	// loop over arguments to check any have :\\ or :\ or \
	// If So, they are paths. convert them to WSL Path

	// Only do it if arguments are greater than or equal to 1
	if len(arguments) >= 1 {
		for index, arg := range arguments {
			if strings.Contains(arg, "\\") {

				// convert path to wsl
				wslpath, err := utils.Win2Wsl(arg)
				if err != nil {
					fmt.Printf(err.Error())
					utils.PrintError("Something went wront while getting path is wsl format.\n", 6)
				}
				// Assign wsl path to argument element
				arguments[index] = wslpath
			}
		}
	} else {
		utils.PrintError("Invalid Nummer of arguments.\n", 8)
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	//Pprepend git to arguments
	arguments = append([]string{"git"}, arguments...)
	cmd := exec.Command("wsl", arguments...)

	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()

	var errStdout, errStderr error
	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderr := io.MultiWriter(os.Stderr, &stderrBuf)
	err := cmd.Start()
	if err != nil {
		fmt.Printf("cmd.Start() failed with '%s'\n", err)
		os.Exit(9)
	}

	go func() {
		_, errStdout = io.Copy(stdout, stdoutIn)
	}()

	go func() {
		_, errStderr = io.Copy(stderr, stderrIn)
	}()

	err = cmd.Wait()
	if err != nil {
		fmt.Printf("cmd.Run() failed with %s\n", err)
		os.Exit(9)
	}
	if errStdout != nil || errStderr != nil {
		utils.PrintError("Failed to capture stdout or stderr\n", 10)
	}
	//outStr, errStr := string(stdoutBuf.Bytes()), string(stderrBuf.Bytes())
	// Display Output
	//fmt.Printf("\nOutput:\n%s\n", outStr)
	//fmt.Printf("\nErrors:\n%s\n", errStr)
	//
}
