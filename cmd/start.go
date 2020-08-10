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
	"log"
	"net/http"

	"github.com/gorilla/mux"
	//"github.com/gorilla/mux"
	//"github.com/jmozah/intOS-dfs/pkg/api"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "starts a HTTP server for the dfs",
	Long: `Serves all the dfs commands through an HTTP server so that the upper layers
can consume it.`,
	Run: func(cmd *cobra.Command, args []string) {
		startHttpService()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func startHttpService() {
	router := mux.NewRouter().StrictSlash(true)
	//router.HandleFunc("/pod/new/{podName}", api.CreatePod).Methods("POST")
	//router.HandleFunc("/pod/del/{podName}", api.DeletePod).Methods("DELETE")
	//router.HandleFunc("/pod/login/{podName}", api.LoginPod).Methods("POST")
	//router.HandleFunc("/pod/logout/{podName}", api.LogoutPod).Methods("POST")
	//router.HandleFunc("/pod/ls/{dirName}", api.ListPod).Methods("GET")
	//router.HandleFunc("/pod/stat/{fileOrDirName}", api.InfoPod).Methods("GET")
	//router.HandleFunc("/pod/sync/{podName}", api.SyncPod).Methods("POST")
	//
	//router.HandleFunc("/mkdir/{dirName}", api.Mkdir).Methods("POST")
	//router.HandleFunc("/rmdir/{dirName}", api.Rmdir).Methods("DELETE")
	//router.HandleFunc("/copyFromLocal/", api.CopyFromLocal).Methods("POST")
	//router.HandleFunc("/copyToLocal/", api.CopyToLocal).Methods("POST")
	//router.HandleFunc("/rm/{fileName}", api.RemoveFile).Methods("DELETE")

	log.Fatal(http.ListenAndServe("127.0.0.1:"+httpPort, router))
}
