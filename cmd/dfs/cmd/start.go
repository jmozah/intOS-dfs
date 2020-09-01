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
	"github.com/jmozah/intOS-dfs/pkg/api"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
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
	router := mux.NewRouter()

	// Web page handlers
	router.HandleFunc("/", handler.WebHandlers.IndexPageHandler)
	router.HandleFunc("/login_page", handler.WebHandlers.LoginPageHandler).Methods("POST")
	router.HandleFunc("/signup_page", handler.WebHandlers.SignupPageHandler).Methods("POST")

	// User account related handlers which does not login need middleware
	router.HandleFunc("/v0/user/signup", handler.UserSignupHandler).Methods("POST")
	router.HandleFunc("/v0/user/present", handler.UserPresentHandler).Methods("POST")
	router.HandleFunc("/v0/user/login", handler.UserLoginHandler).Methods("POST")

	// user account related handlers which require login middleware
	userRouter := router.PathPrefix("/v0/user/").Subrouter()
	userRouter.Use(handler.LoginMiddleware)
	userRouter.HandleFunc("/delete", handler.UserDeleteHandler).Methods("POST")
	userRouter.HandleFunc("/logout", handler.UserLogoutHandler).Methods("POST")

	// pod related handlers
	podRouter := router.PathPrefix("/v0/pod/").Subrouter()
	podRouter.Use(handler.LoginMiddleware)
	podRouter.HandleFunc("/new", handler.PodCreateHandler).Methods("POST")
	podRouter.HandleFunc("/delete", handler.PodDeleteHandler).Methods("POST")
	podRouter.HandleFunc("/open", handler.PodOpenHandler).Methods("POST")
	podRouter.HandleFunc("/close", handler.PodCloseHandler).Methods("POST")
	podRouter.HandleFunc("/ls", handler.PodListHandler).Methods("POST")
	podRouter.HandleFunc("/stat", handler.PodStatHandler).Methods("POST")
	podRouter.HandleFunc("/sync", handler.PodSyncHandler).Methods("POST")

	// directory related handlers
	dirRouter := router.PathPrefix("/v0/dir/").Subrouter()
	dirRouter.Use(handler.LoginMiddleware)
	dirRouter.HandleFunc("/mkdir", handler.DirectoryMkdirHandler).Methods("POST")
	dirRouter.HandleFunc("/rmdir", handler.DirectoryRmdirHandler).Methods("POST")
	dirRouter.HandleFunc("/ls", handler.DirectoryLsHandler).Methods("POST")
	dirRouter.HandleFunc("/stat", handler.DirectoryStatHandler).Methods("POST")

	// file related handlers
	fileRouter := router.PathPrefix("/v0/file/").Subrouter()
	fileRouter.Use(handler.LoginMiddleware)
	fileRouter.HandleFunc("/download", handler.FileDownloadHandler).Methods("POST")
	fileRouter.HandleFunc("/upload", handler.FileUploadHandler).Methods("POST")
	fileRouter.HandleFunc("/stat", handler.FileStatHandler).Methods("POST")
	fileRouter.HandleFunc("/delete", handler.FileDeleteHandler).Methods("POST")

	http.Handle("/", router)
	handler := cors.Default().Handler(router)

	fmt.Println("listening on port:", httpPort)
	err := http.ListenAndServe(":"+httpPort, handler)
	if err != nil {
		fmt.Println("listenAndServe: %w", err)
		return
	}

}
