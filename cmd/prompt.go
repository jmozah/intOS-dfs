/*
Copyright © 2020 intOS Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jmozah/intOS-dfs/pkg/utils"

	prompt "github.com/c-bata/go-prompt"
	"github.com/spf13/cobra"

	"github.com/jmozah/intOS-dfs/pkg/dfs"
	"github.com/jmozah/intOS-dfs/pkg/pod"
)

const (
	DefaultPrompt   = "dfs >>"
	PodSeperator    = ">>"
	PromptSeperator = "> "
)

var (
	currentPrompt  string
	currentPodInfo *pod.Info
	dfsAPI         *dfs.DfsAPI
)

// promptCmd represents the prompt command
var promptCmd = &cobra.Command{
	Use:   "prompt",
	Short: "a REPL to interact with intOS's dfs",
	Long: `A command prompt where you can interact with the distributed
file system of the intOS.`,
	Run: func(cmd *cobra.Command, args []string) {
		dfsAPI = dfs.NewDfsAPI(dataDir, beeHost, beePort)
		if !dfsAPI.IsInitialized(dataDir) {
			fmt.Println("Looks like you have not initialised dfs. Please run \"init\" command to start using dfs.")
		} else {
			err := dfsAPI.Init(dataDir, "")
			if err != nil {
				fmt.Println("error starting dfs: ", err)
				os.Exit(1)
			}
		}
		initPrompt()
	},
}

func init() {
	rootCmd.AddCommand(promptCmd)
}

func initPrompt() {
	currentPrompt = DefaultPrompt
	p := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix(currentPrompt),
		prompt.OptionLivePrefix(changeLivePrefix),
		prompt.OptionTitle("dfs"),
	)
	p.Run()
}

func changeLivePrefix() (string, bool) {
	return currentPrompt, true
}

var suggestions = []prompt.Suggest{
	{Text: "init", Description: "initialize the dfs file system"},
	{Text: "pod new", Description: "create a new pod"},
	{Text: "pod del", Description: "delete a existing pod"},
	{Text: "pod login", Description: "login to a existing pod"},
	{Text: "pod logout", Description: "logout from a logged in pod"},
	{Text: "pod ls", Description: "list all the existing pods"},
	{Text: "pod stat", Description: "show the metadata of a pod"},
	{Text: "pod sync", Description: "sync the pod from swarm"},
	{Text: "cd", Description: "change path"},
	{Text: "copyToLocal", Description: "copy file from dfs to local machine"},
	{Text: "copyFromLocal", Description: "copy file from local machine to dfs"},
	{Text: "exit", Description: "exit dfs-prompt"},
	{Text: "head", Description: "show few starting lines of a file"},
	{Text: "help", Description: "show usage"},
	{Text: "ls", Description: "list all the file and directories in the current path"},
	{Text: "mkdir", Description: "make a new directory"},
	{Text: "rmdir", Description: "remove a existing directory"},
	{Text: "pwd", Description: "show the current working directory"},
	{Text: "rm", Description: "remove a file"},
}

func completer(in prompt.Document) []prompt.Suggest {
	w := in.GetWordBeforeCursor()
	if w == "" {
		return []prompt.Suggest{}
	}
	return prompt.FilterHasPrefix(suggestions, w, true)
}

func executor(in string) {
	in = strings.TrimSpace(in)
	blocks := strings.Split(in, " ")
	switch blocks[0] {
	case "help":
		help()
	case "exit":
		os.Exit(0)
	case "init":
		err := dfsAPI.Init(dataDir, "")
		if err != nil {
			fmt.Println("init failed: ", err)
			return
		}
		currentPrompt = DefaultPrompt
	case "pod":
		if len(blocks) < 2 {
			log.Println("invalid command.")
			help()
			return
		}
		switch blocks[1] {
		case "new":
			if len(blocks) < 3 {
				fmt.Println("invalid command. Missing \"name\" argument ")
				return
			}
			podName := blocks[2]
			podInfo, err := dfsAPI.CreatePod(podName, "")
			if err != nil {
				fmt.Println("could not create pod: ", err)
				return
			}
			currentPrompt = getCurrentPrompt(podInfo)
			currentPodInfo = podInfo
		case "del":
			lastPrompt := currentPrompt
			if len(blocks) < 3 {
				fmt.Println("invalid command. Missing \"name\" argument ")
				return
			}
			podName := blocks[2]
			err := dfsAPI.DeletePod(podName)
			if err != nil {
				fmt.Println("could not delete pod: ", err)
				return
			}
			fmt.Println("successfully deleted pod: ", podName)
			if podName == currentPodInfo.GetCurrentPodNameOnly() {
				currentPrompt = DefaultPrompt
			} else {
				currentPrompt = lastPrompt
			}
		case "login":
			if len(blocks) < 3 {
				fmt.Println("invalid command. Missing \"name\" argument ")
				return
			}
			podName := blocks[2]
			podInfo, err := dfsAPI.LoginPod(podName, "")
			if err != nil {
				fmt.Println("Login failed: ", err)
				return
			}
			currentPrompt = getCurrentPrompt(podInfo)
			currentPodInfo = podInfo
		case "logout":
			if !isLoggedIn() {
				return
			}
			err := dfsAPI.LogoutPod(currentPodInfo.GetCurrentPodNameOnly())
			if err != nil {
				fmt.Println("error logging out: ", err)
				return
			}
			currentPrompt = DefaultPrompt
			currentPodInfo = nil
		case "stat":
			if !isLoggedIn() {
				return
			}
			podStat, err := dfsAPI.PodStat(currentPodInfo.GetCurrentPodNameOnly())
			if err != nil {
				fmt.Println("error getting stat: ", err)
				return
			}
			fmt.Println("Version          : ", podStat.Version)
			fmt.Println("pod Name         : ", podStat.PodName)
			fmt.Println("Path             : ", podStat.PodPath)
			fmt.Println("Creation Time    :", podStat.CreationTime)
			fmt.Println("Access Time      :", podStat.AccessTime)
			fmt.Println("Modification Time:", podStat.ModificationTime)
			currentPrompt = getCurrentPrompt(currentPodInfo)
		case "sync":
			if !isLoggedIn() {
				return
			}
			err := dfsAPI.SyncPod(currentPodInfo.GetCurrentPodNameOnly())
			if err != nil {
				fmt.Println("could not sync pod: ", err)
				return
			}
			fmt.Println("pod synced.")
			currentPrompt = getCurrentPrompt(currentPodInfo)
		case "ls":
			err := dfsAPI.ListPods()
			if err != nil {
				fmt.Println("error while listing pods: %w", err)
				return
			}
			currentPrompt = getCurrentPrompt(currentPodInfo)
		default:
			fmt.Println("invalid pod command!!")
			help()
		} // end of pod commands
	case "cd":
		if !isLoggedIn() {
			return
		}
		if len(blocks) < 2 {
			fmt.Println("invalid command. Missing one or more arguments")
			return
		}
		podInfo, err := dfsAPI.ChangeDirectory(currentPodInfo.GetCurrentPodNameOnly(), blocks[1])
		if err != nil {
			fmt.Println("cd failed: ", err)
			return
		}
		currentPodInfo = podInfo
		currentPrompt = getCurrentPrompt(currentPodInfo)
	case "ls":
		if !isLoggedIn() {
			return
		}
		listing, err := dfsAPI.ListDir(currentPodInfo.GetCurrentPodNameOnly())
		if err != nil {
			fmt.Println("ls failed: ", err)
			return
		}
		for _, l := range listing {
			fmt.Println(l)
		}
	case "copyToLocal":
		if !isLoggedIn() {
			return
		}
		if len(blocks) < 3 {
			fmt.Println("invalid command. Missing one or more arguments")
			return
		}
		err := dfsAPI.CopyToLocal(currentPodInfo.GetCurrentPodNameOnly(), blocks[1], blocks[2])
		if err != nil {
			fmt.Println("copyToLocal failed: ", err)
			return
		}
	case "copyFromLocal":
		if !isLoggedIn() {
			return
		}
		if len(blocks) < 4 {
			fmt.Println("invalid command. Missing one or more arguments")
			return
		}
		err := dfsAPI.CopyFromLocal(currentPodInfo.GetCurrentPodNameOnly(), blocks[1], blocks[2], blocks[3])
		if err != nil {
			fmt.Println("copyFromLocal failed: ", err)
			return
		}
	case "mkdir":
		if !isLoggedIn() {
			return
		}
		if len(blocks) < 2 {
			fmt.Println("invalid command. Missing one or more arguments")
			return
		}
		err := dfsAPI.Mkdir(currentPodInfo.GetCurrentPodNameOnly(), blocks[1])
		if err != nil {
			fmt.Println("mkdir failed: ", err)
			return
		}
	case "rmdir":
		if !isLoggedIn() {
			return
		}
		if len(blocks) < 2 {
			fmt.Println("invalid command. Missing one or more arguments")
			return
		}
		err := dfsAPI.RmDir(currentPodInfo.GetCurrentPodNameOnly(), blocks[1])
		if err != nil {
			fmt.Println("rmdir failed: ", err)
			return
		}
	case "cat":
		if !isLoggedIn() {
			return
		}
		if len(blocks) < 2 {
			fmt.Println("invalid command. Missing one or more arguments")
			return
		}
		err := dfsAPI.Cat(currentPodInfo.GetCurrentPodNameOnly(), blocks[1])
		if err != nil {
			fmt.Println("cat failed: ", err)
			return
		}
	case "stat":
		if !isLoggedIn() {
			return
		}
		if len(blocks) < 2 {
			fmt.Println("invalid command. Missing one or more arguments")
			return
		}
		err := dfsAPI.DirectoryOrFileStat(currentPodInfo.GetCurrentPodNameOnly(), blocks[1])
		if err != nil {
			fmt.Println("stat failed: ", err)
			return
		}
	case "pwd":
		if !isLoggedIn() {
			return
		}
		if currentPodInfo.IsCurrentDirRoot() {
			fmt.Println("/")
		} else {
			podDir := currentPodInfo.GetCurrentPodPathAndName()
			curDir := strings.TrimPrefix(currentPodInfo.GetCurrentDirPathAndName(), podDir)
			fmt.Println(curDir)
		}
	case "rm":
		if !isLoggedIn() {
			return
		}
		if len(blocks) < 2 {
			fmt.Println("invalid command. Missing one or more arguments")
			return
		}
		err := dfsAPI.RemoveFile(currentPodInfo.GetCurrentPodNameOnly(), blocks[1])
		if err != nil {
			fmt.Println("rm failed: ", err)
			return
		}
	case "mv":
		fmt.Println("not yet implemented")
	case "head":
		fmt.Println("not yet implemented")
	default:
		fmt.Println("invalid command")
	}
}

func help() {
	fmt.Println("Usage: <command> <sub-command> (args1) (args2) ...")
	fmt.Println("commands:")
	fmt.Println(" - pod <new> (pod-name) - create a new pod and login to that pod")
	fmt.Println(" - pod <del> (pod-name) - Deletes a already created pod")
	fmt.Println(" - pod <login> (pod-name) - login to a already created pod.")
	fmt.Println(" - pod <stat> (pod-name) - display meta information about a pod")
	fmt.Println(" - pod <sync> (pod-name) - sync the contents of a logged in pod from Swarm")
	fmt.Println(" - pod <logout>  - logout of a logged in pod")
	fmt.Println(" - pod <ls> - lists all the pods created for this account")

	fmt.Println(" - cd <directory name>")
	fmt.Println(" - ls ")
	fmt.Println(" - copyToLocal <source file in pod, destination directory in local fs>")
	fmt.Println(" - copyFromLocal <source file in local fs, destination directory in pod, block size in MB>")
	fmt.Println(" - mkdir <directory name>")
	fmt.Println(" - rmdir <directory name>")
	fmt.Println(" - rm <file name>")
	fmt.Println(" - pwd - show present working directory")
	fmt.Println(" - head <no of lines>")
	fmt.Println(" - cat  - stream the file to stdout")
	fmt.Println(" - stat <file name or directory name> - shows the information about a file or directory")
	fmt.Println(" - help - display this help")
	fmt.Println(" - exit - exits from the prompt")

}

func getCurrentPrompt(podInfo *pod.Info) string {
	if podInfo == nil || podInfo.GetCurrentDirInode() == nil {
		return DefaultPrompt
	}
	if podInfo.GetCurrentPodPathOnly() == podInfo.GetCurrentDirPathOnly() {
		return DefaultPrompt + " " + podInfo.GetCurrentDirNameOnly() + " " + PodSeperator
	} else {
		podPathAndName := podInfo.GetCurrentPodPathAndName()
		pathExceptPod := strings.TrimPrefix(podInfo.GetCurrentDirPathOnly(), podPathAndName)
		return DefaultPrompt + " " +
			podInfo.GetCurrentPodNameOnly() + " " + PodSeperator + " " +
			pathExceptPod + utils.PathSeperator + podInfo.GetCurrentDirNameOnly() + PromptSeperator
	}
}

func isLoggedIn() bool {
	if currentPodInfo == nil {
		fmt.Println("login to do the operation")
		return false
	}
	return true
}
