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

	//"github.com/gorilla/mux"
	"github.com/gorilla/mux"
	//"github.com/jmozah/intOS-dfs/pkg/api"
	"github.com/rs/cors"
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
	fs := http.FileServer(http.Dir("pkg/web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	//handler := cors.Default().Handler(mux)
	router := mux.NewRouter()

	// Web page handlers
	router.HandleFunc("/", handler.WebHandlers.IndexPageHandler)
	router.HandleFunc("/login_page", handler.WebHandlers.LoginPageHandler).Methods("POST")
	router.HandleFunc("/signup_page", handler.WebHandlers.SignupPageHandler).Methods("POST")

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
	router.HandleFunc("/v0/dir/ls", handler.DirectoryLsHandler).Methods("POST")
	router.HandleFunc("/v0/dir/stat", handler.DirectoryStatHandler).Methods("POST")

	// file related handlers
	router.HandleFunc("/v0/file/download", handler.FileDownloadHandler).Methods("POST")
	router.HandleFunc("/v0/file/upload", handler.FileUploadHandler).Methods("POST")
	router.HandleFunc("/v0/file/stat", handler.FileStatHandler).Methods("POST")
	router.HandleFunc("/v0/file/delete", handler.FileDeleteHandler).Methods("POST")

	http.Handle("/", router)

	fmt.Println("listening on port:", httpPort)

	//handler := cors.Default().Handler(router)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000", "http://localhost:9090"},
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		// Debug: true,
	})
	
	// Insert the middleware
	handler := c.Handler(router)

	err := http.ListenAndServe(":"+httpPort, handler)
	if err != nil {
		fmt.Println("listenAndServe: %w", err)
		return
	}

}
