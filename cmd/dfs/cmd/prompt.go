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
	DefaultPrompt   = "dfs"
	UserSeperator   = ">>>"
	PodSeperator    = ">>"
	PromptSeperator = "> "
)

var (
	currentUser    string
	currentPodInfo *pod.Info
	currentPrompt  string
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
		initPrompt()
	},
}

func init() {
	rootCmd.AddCommand(promptCmd)
}

func initPrompt() {
	currentPrompt = DefaultPrompt + " " + UserSeperator
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
	{Text: "user new", Description: "create a new user"},
	{Text: "user del", Description: "delete a existing user"},
	{Text: "user login", Description: "login to a existing user"},
	{Text: "user logout", Description: "logout from a logged in user"},
	{Text: "user present", Description: "is user present"},
	{Text: "user ls", Description: "list all users"},
	{Text: "pod new", Description: "create a new pod for a user"},
	{Text: "pod del", Description: "delete a existing pod of a user"},
	{Text: "pod open", Description: "open to a existing pod of a user"},
	{Text: "pod close", Description: "close a already opened pod of a user"},
	{Text: "pod ls", Description: "list all the existing pods of  auser"},
	{Text: "pod stat", Description: "show the metadata of a pod of a user"},
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
	case "user":
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
			userName := blocks[2]
			ref, mnemonic, err := dfsAPI.CreateUser(userName, "")
			if err != nil {
				fmt.Println("create user: ", err)
				return
			}
			fmt.Println("user created with address ", ref)
			fmt.Println("Please store the following 24 words safely")
			fmt.Println("if you loose this, you cannot recover the data in-case of an emergency.")
			fmt.Println("you can also use this mnemonic to access the data from another device")
			fmt.Println("=============== Mnemonic ==========================")
			fmt.Println(mnemonic)
			fmt.Println("=============== Mnemonic ==========================")
			currentUser = userName
			currentPodInfo = nil
			currentPrompt = getCurrentPrompt()
		case "del":
			if len(blocks) < 3 {
				fmt.Println("invalid command. Missing \"name\" argument ")
				return
			}
			userName := blocks[2]
			err := dfsAPI.DeleteUser(userName, "")
			if err != nil {
				fmt.Println("delete user: ", err)
				return
			}
			currentUser = ""
			currentPodInfo = nil
			currentPrompt = getCurrentPrompt()
		case "login":
			if len(blocks) < 3 {
				fmt.Println("invalid command. Missing \"name\" argument ")
				return
			}
			userName := blocks[2]
			err := dfsAPI.LoginUser(userName, "")
			if err != nil {
				fmt.Println("login user: ", err)
				return
			}
			currentUser = userName
			currentPodInfo = nil
			currentPrompt = getCurrentPrompt()
		case "logout":
			err := dfsAPI.LogoutUser(currentUser)
			if err != nil {
				fmt.Println("logout user: ", err)
				return
			}
			currentUser = ""
			currentPodInfo = nil
			currentPrompt = getCurrentPrompt()
		case "present":
			if len(blocks) < 3 {
				fmt.Println("invalid command. Missing \"name\" argument ")
				return
			}
			userName := blocks[2]
			yes := dfsAPI.IsUserNameAvailable(userName)
			if yes {
				fmt.Println("user name: available")
			} else {
				fmt.Println("user name: not available")
			}
			currentPrompt = getCurrentPrompt()
		case "ls":
			users := dfsAPI.ListAllUsers()
			for _, user := range users {
				fmt.Println(user)
			}
			currentPrompt = getCurrentPrompt()
		default:
			fmt.Println("invalid user command")
		}
	case "pod":
		if currentUser == "" {
			fmt.Println("login as a user to execute these commands")
			return
		}
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
			podInfo, err := dfsAPI.CreatePod(currentUser, podName, "")
			if err != nil {
				fmt.Println("could not create pod: ", err)
				return
			}
			currentPodInfo = podInfo
			currentPrompt = getCurrentPrompt()
		case "del":
			lastPrompt := currentPrompt
			if len(blocks) < 3 {
				fmt.Println("invalid command. Missing \"name\" argument ")
				return
			}
			podName := blocks[2]
			err := dfsAPI.DeletePod(currentUser, podName)
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
		case "open":
			if len(blocks) < 3 {
				fmt.Println("invalid command. Missing \"name\" argument ")
				return
			}
			podName := blocks[2]
			podInfo, err := dfsAPI.OpenPod(currentUser, podName, "")
			if err != nil {
				fmt.Println("Login failed: ", err)
				return
			}
			currentPodInfo = podInfo
			currentPrompt = getCurrentPrompt()
		case "close":
			if !isPodOpened() {
				return
			}
			err := dfsAPI.ClosePod(currentUser, currentPodInfo.GetCurrentPodNameOnly())
			if err != nil {
				fmt.Println("error logging out: ", err)
				return
			}
			currentPrompt = DefaultPrompt + " " + UserSeperator
			currentPodInfo = nil
		case "stat":
			if !isPodOpened() {
				return
			}
			podStat, err := dfsAPI.PodStat(currentUser, currentPodInfo.GetCurrentPodNameOnly())
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
			currentPrompt = getCurrentPrompt()
		case "sync":
			if !isPodOpened() {
				return
			}
			err := dfsAPI.SyncPod(currentUser, currentPodInfo.GetCurrentPodNameOnly())
			if err != nil {
				fmt.Println("could not sync pod: ", err)
				return
			}
			fmt.Println("pod synced.")
			currentPrompt = getCurrentPrompt()
		case "ls":
			pods, err := dfsAPI.ListPods(currentUser)
			if err != nil {
				fmt.Println("error while listing pods: %w", err)
				return
			}
			for _, pod := range pods {
				fmt.Println(pod)
			}
			fmt.Println("")
			currentPrompt = getCurrentPrompt()
		default:
			fmt.Println("invalid pod command!!")
			help()
		} // end of pod commands
	case "cd":
		if !isPodOpened() {
			return
		}
		if len(blocks) < 2 {
			fmt.Println("invalid command. Missing one or more arguments")
			return
		}
		podInfo, err := dfsAPI.ChangeDirectory(currentUser, currentPodInfo.GetCurrentPodNameOnly(), blocks[1])
		if err != nil {
			fmt.Println("cd failed: ", err)
			return
		}
		currentPodInfo = podInfo
		currentPrompt = getCurrentPrompt()
	case "ls":
		if !isPodOpened() {
			return
		}
		fl, dl, err := dfsAPI.ListDir(currentUser, currentPodInfo.GetCurrentPodNameOnly(), "")
		if err != nil {
			fmt.Println("ls failed: ", err)
			return
		}
		for _, l := range fl {
			fmt.Println(l)
		}
		for _, l := range dl {
			fmt.Println(l)
		}
	case "copyToLocal":
		if !isPodOpened() {
			return
		}
		if len(blocks) < 3 {
			fmt.Println("invalid command. Missing one or more arguments")
			return
		}
		err := dfsAPI.CopyToLocal(currentUser, currentPodInfo.GetCurrentPodNameOnly(), blocks[1], blocks[2])
		if err != nil {
			fmt.Println("copyToLocal failed: ", err)
			return
		}
	case "copyFromLocal":
		if !isPodOpened() {
			return
		}
		if len(blocks) < 4 {
			fmt.Println("invalid command. Missing one or more arguments")
			return
		}
		err := dfsAPI.CopyFromLocal(currentUser, currentPodInfo.GetCurrentPodNameOnly(), blocks[1], blocks[2], blocks[3])
		if err != nil {
			fmt.Println("copyFromLocal failed: ", err)
			return
		}
	case "mkdir":
		if !isPodOpened() {
			return
		}
		if len(blocks) < 2 {
			fmt.Println("invalid command. Missing one or more arguments")
			return
		}
		err := dfsAPI.Mkdir(currentUser, currentPodInfo.GetCurrentPodNameOnly(), blocks[1], "")
		if err != nil {
			fmt.Println("mkdir failed: ", err)
			return
		}
	case "rmdir":
		if !isPodOpened() {
			return
		}
		if len(blocks) < 2 {
			fmt.Println("invalid command. Missing one or more arguments")
			return
		}
		err := dfsAPI.RmDir(currentUser, currentPodInfo.GetCurrentPodNameOnly(), blocks[1])
		if err != nil {
			fmt.Println("rmdir failed: ", err)
			return
		}
	case "cat":
		if !isPodOpened() {
			return
		}
		if len(blocks) < 2 {
			fmt.Println("invalid command. Missing one or more arguments")
			return
		}
		err := dfsAPI.Cat(currentUser, currentPodInfo.GetCurrentPodNameOnly(), blocks[1])
		if err != nil {
			fmt.Println("cat failed: ", err)
			return
		}
	case "stat":
		if !isPodOpened() {
			return
		}
		if len(blocks) < 2 {
			fmt.Println("invalid command. Missing one or more arguments")
			return
		}
		fs, err := dfsAPI.FileStat(currentUser, currentPodInfo.GetCurrentPodNameOnly(), blocks[1])
		if err != nil {
			fmt.Println("stat failed: ", err)
			return
		}
		fmt.Println("Account 	: ", fs.Account)
		fmt.Println("PodName 	: ", fs.PodName)
		fmt.Println("File Path	: ", fs.FilePath)
		fmt.Println("File Name	: ", fs.FileName)
		fmt.Println("File Size	: ", fs.FileSize)
		fmt.Println("Cr. Time	: ", fs.CreationTime)
		fmt.Println("Mo. Time	: ", fs.ModificationTime)
		fmt.Println("Ac. Time	: ", fs.AccessTime)
		for _, b := range fs.Blocks {
			blkStr := fmt.Sprintf("%s, 0x%s, %d bytes", b.Name, b.Reference, b.Size)
			fmt.Println(blkStr)
		}
	case "pwd":
		if !isPodOpened() {
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
		if !isPodOpened() {
			return
		}
		if len(blocks) < 2 {
			fmt.Println("invalid command. Missing one or more arguments")
			return
		}
		err := dfsAPI.DeleteFile(currentUser, currentPodInfo.GetCurrentPodNameOnly(), blocks[1])
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
	fmt.Println(" - user <new> (user-name) - create a new user and login as that user")
	fmt.Println(" - user <del> (user-name) - deletes a already created user")
	fmt.Println(" - user <login> (user-name) - login as a given user")
	fmt.Println(" - user <logout> (user-name) - logout as user")
	fmt.Println(" - user <ls> - lists all the user present in this instance")

	fmt.Println(" - pod <new> (pod-name) - create a new pod for the logged in user and opens the pod")
	fmt.Println(" - pod <del> (pod-name) - deletes a already created pod of the user")
	fmt.Println(" - pod <open> (pod-name) - open a already created pod")
	fmt.Println(" - pod <stat> (pod-name) - display meta information about a pod")
	fmt.Println(" - pod <sync> (pod-name) - sync the contents of a logged in pod from Swarm")
	fmt.Println(" - pod <close>  - close a opened pod")
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

func getCurrentPrompt() string {
	currPrompt := getUserPrompt()
	podPrompt := getPodPrompt()
	if podPrompt != "" {
		currPrompt = currPrompt + " " + podPrompt + " " + PodSeperator
	}
	dirPrompt := getCurrentDirPrompt()
	if dirPrompt != "" {
		currPrompt = currPrompt + " " + dirPrompt + " " + PromptSeperator
	}
	return currPrompt
}

func isPodOpened() bool {
	if currentPodInfo == nil {
		fmt.Println("open the pod to do the operation")
		return false
	}
	return true
}

func getUserPrompt() string {
	if currentUser == "" {
		return DefaultPrompt + " " + UserSeperator
	} else {
		return DefaultPrompt + "@" + currentUser + " " + UserSeperator
	}
}

func getPodPrompt() string {
	if currentPodInfo != nil {
		return currentPodInfo.GetCurrentPodNameOnly()
	} else {
		return ""
	}
}

func getCurrentDirPrompt() string {
	currentDir := ""
	if currentPodInfo != nil {
		if currentPodInfo.IsCurrentDirRoot() {
			return utils.PathSeperator
		}
		podPathAndName := currentPodInfo.GetCurrentPodPathAndName()
		pathExceptPod := strings.TrimPrefix(currentPodInfo.GetCurrentDirPathOnly(), podPathAndName)
		currentDir = pathExceptPod + utils.PathSeperator + currentPodInfo.GetCurrentDirNameOnly()
	}
	return currentDir
}
