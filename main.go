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

	//Check if not on Windows
	if runtime.GOOS != "windows" {
		utils.PrintError("Not running on Windows.\n", 1)
	}

	//Get all the arguments
	arguments = append(os.Args[1:])

	// Edge case: no arguments, append help to args
	if len(arguments) < 1 {
		arguments = append([]string{"help"}, arguments...)
	}

	//Version Info
	if strings.ToLower(arguments[0]) == "--wsl-git-version" {
		fmt.Printf("\nWSL-GIT Version - %s\n\n", version.VERSION)
		return
	}
	// Help
	if strings.ToLower(arguments[0]) == "--wsl-git-help" {
		utils.Usage()
		return
	}

	// Check if wsl is present, if not Exit
	utils.CheckwslExists()

	//if strings.ToLower(arguments[0]) == "--wsl-git-print-args" {
	//	fmt.Printf("Arguments passed to %s are\n%v", os.Args[0], arguments)
	//		return
	//	}

	// VS Code hack
	// VS code runs git rev-parse --show-toplevel
	// every time to determine if directrory is git repository
	// We need to convert wsl path returned to windows path
	// TODO for V 1.0

	if arguments[0] == "rev-parse" && arguments[1] == "--show-toplevel" {
		toplevelcmd := exec.Command("wsl", "git", "rev-parse", "--show-toplevel")
		out, err := toplevelcmd.CombinedOutput()
		if err != nil {
			fmt.Printf("toplevelcmd.Run() failed with %s\n", err)
			utils.PrintError("Cannot continue\n", 11)
		}

		windowspath, err := utils.Wsl2Win(string(out))
		if err != nil {
			fmt.Printf(err.Error())
			utils.PrintError("Something went wont while converting path to windows format.\n", 12)
		} else {
			fmt.Print(windowspath)
			return
		}
	}

	// This is just a hack which should work most of the times
	// loop over arguments to check any have :\\ or :\ or \
	// If So, they are paths. convert them to WSL Path

	// Only do it if arguments are greater than or equal to 1
	if len(arguments) >= 1 {
		for index, arg := range arguments {
			if strings.Contains(arg, "\\") {

				// Convert path to wsl
				wslpath, err := utils.Win2Wsl(arg)
				if err != nil {
					fmt.Printf(err.Error())
					utils.PrintError("Something went wrong while getting path is wsl format.\n", 11)
				}
				// Assign wsl path to argument element
				arguments[index] = wslpath
			}
		}
	} else {
		utils.PrintError("Invalid Number of arguments.\n", 15)
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	//Prepends git to arguments
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
		os.Exit(11)
	}

	go func() {
		_, errStdout = io.Copy(stdout, stdoutIn)
	}()

	go func() {
		_, errStderr = io.Copy(stderr, stderrIn)
	}()

	err = cmd.Wait()
	if err != nil {
		fmt.Printf("git %s failed with %s\n", os.Args[1], err)
		os.Exit(11)
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
