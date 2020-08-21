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

	"github.com/jmozah/intOS-dfs/pkg/api"
)

var handler *api.Handler

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "starts a HTTP server for the dfs",
	Long: `Serves all the dfs commands through an HTTP server so that the upper layers
can consume it.`,
	Run: func(cmd *cobra.Command, args []string) {
		handler = api.NewHandler(dataDir, beeHost, beePort)
		startHttpService()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func startHttpService() {
	router := mux.NewRouter()

	// User account related handlers
	router.HandleFunc("/v0/user/signup", handler.UserSignupHandler).Methods("POST")
	router.HandleFunc("/v0/user/delete", handler.UserDeleteHandler).Methods("POST")
	router.HandleFunc("/v0/user/login", handler.UserLoginHandler).Methods("POST")
	router.HandleFunc("/v0/user/logout", handler.UserLogoutHandler).Methods("POST")
	router.HandleFunc("/v0/user/present", handler.UserPresentHandler).Methods("POST")

	// pod related handlers
	router.HandleFunc("/v0/pod/new", handler.PodCreateHandler).Methods("POST")
	router.HandleFunc("/v0/pod/delete", handler.PodDeleteHandler).Methods("POST")
	router.HandleFunc("/v0/pod/open", handler.PodOpenHandler).Methods("POST")
	router.HandleFunc("/v0/pod/close", handler.PodCloseHandler).Methods("POST")
	router.HandleFunc("/v0/pod/ls", handler.PodListHandler).Methods("POST")
	router.HandleFunc("/v0/pod/stat", handler.PodStatHandler).Methods("POST")
	router.HandleFunc("/v0/pod/sync", handler.PodSyncHandler).Methods("POST")

	// directory related handlers
	router.HandleFunc("/v0/dir/mkdir", handler.DirectoryMkdirHandler).Methods("POST")
	router.HandleFunc("/v0/dir/rmdir", handler.DirectoryRmdirHandler).Methods("POST")
	router.HandleFunc("/v0/dir/cd", handler.DirectoryCdHandler).Methods("POST")
	router.HandleFunc("/v0/dir/ls", handler.DirectoryLsHandler).Methods("POST")
	router.HandleFunc("/v0/dir/stat", handler.DirectoryStatHandler).Methods("POST")
	router.HandleFunc("/v0/dir/pwd", handler.DirectoryPwdHandler).Methods("POST")

	// file related handlers
	router.HandleFunc("/v0/file/download", handler.FileDownloadHandler).Methods("POST")
	router.HandleFunc("/v0/file/upload", handler.FileUploadHandler).Methods("POST")
	router.HandleFunc("/v0/file/stat", handler.FileStatHandler).Methods("POST")

	http.Handle("/", router)

	fmt.Println("listening on port:", httpPort)
	err := http.ListenAndServe(":"+httpPort, nil)
	if err != nil {
		fmt.Println("listenAndServe: %w", err)
		return
	}

}
