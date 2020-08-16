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
	"net/http"

	"github.com/gorilla/mux"
	//"github.com/gorilla/mux"
	//"github.com/jmozah/intOS-dfs/pkg/api"
	"github.com/spf13/cobra"

	"github.com/jmozah/intOS-dfs/pkg/web"
)

var handler *web.Handler

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "starts a HTTP server for the dfs",
	Long: `Serves all the dfs commands through an HTTP server so that the upper layers
can consume it.`,
	Run: func(cmd *cobra.Command, args []string) {
		handler = web.NewHandler(dataDir, beeHost, beePort)
		startHttpService()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func startHttpService() {
	fs := http.FileServer(http.Dir("pkg/web/public/assets/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	router := mux.NewRouter()
	router.HandleFunc("/", handler.IndexPageHandler)
	router.HandleFunc("/login_page", handler.LoginPageHandler).Methods("POST")
	router.HandleFunc("/signup_page", handler.SignupPageHandler).Methods("POST")

	router.HandleFunc("/user/signup", handler.UserSignupHandler).Methods("POST")
	//router.HandleFunc("/user/delete", handler.UserDeleteHandler).Methods("POST")
	//router.HandleFunc("/user/login", handler.UserLoginHandler).Methods("POST")
	//router.HandleFunc("/user/logout", handler.LogoutHandler).Methods("POST")
	//router.HandleFunc("/user/present", handler.UserPresentHandler).Methods("POST")
	//router.HandleFunc("/user/ls", handler.UserListHandler).Methods("POST")
	//
	//router.HandleFunc("/pod/new", handler.PodCreateHandler).Methods("POST")
	//router.HandleFunc("/pod/delete", handler.PodDeleteHandler).Methods("POST")
	//router.HandleFunc("/pod/open", handler.PodOpenHandler).Methods("POST")
	//router.HandleFunc("/pod/close", handler.PodCloseHandler).Methods("POST")
	//router.HandleFunc("/pod/ls", handler.PodListHandler).Methods("POST")
	//router.HandleFunc("/pod/stat", handler.PodStatHandler).Methods("POST")
	//router.HandleFunc("/pod/sync", handler.PodSyncHandler).Methods("POST")
	//
	//router.HandleFunc("/dir/mkdir", handler.MakeDirectoryHandler).Methods("POST")
	//router.HandleFunc("/dir/rmdir", handler.RemoveDirectoryHandler).Methods("POST")
	//router.HandleFunc("/dir/cd", handler.ChangeDirectoryHandler).Methods("POST")
	//router.HandleFunc("/dir/ls", handler.ListDirectoryHandler).Methods("POST")
	//router.HandleFunc("/dir/stat", handler.StatDirectoryHandler).Methods("POST")
	//router.HandleFunc("/dir/pwd", handler.CurrentDirectoryHandler).Methods("POST")
	//
	//router.HandleFunc("/file/copyToLocal", handler.FileCopyToLocalHandler).Methods("POST")
	//router.HandleFunc("/file/copyFromLocal", handler.FileCopyFromLocalHandler).Methods("POST")
	//router.HandleFunc("/file/stat", handler.FileStatHandler).Methods("POST")

	http.Handle("/", router)

	fmt.Println("listening on port:", httpPort)
	err := http.ListenAndServe(":"+httpPort, nil)
	if err != nil {
		fmt.Println("listenAndServe: %w", err)
		return
	}

}
