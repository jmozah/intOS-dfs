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

package web

import (
	"html/template"

	"github.com/jmozah/intOS-dfs/pkg/dfs"
)

type Handler struct {
	dfsAPI         *dfs.DfsAPI
	indexTmpl      *template.Template
	loginPageTmpl  *template.Template
	signupPageTmpl *template.Template
	errorTmpl      *template.Template
}

func NewHandler(dataDir string, beeHost string, beePort string) *Handler {
	return &Handler{
		dfsAPI:         dfs.NewDfsAPI(dataDir, beeHost, beePort),
		indexTmpl:      template.Must(template.ParseFiles("pkg/web/template/index.html")),
		loginPageTmpl:  template.Must(template.ParseFiles("pkg/web/template/login.html")),
		signupPageTmpl: template.Must(template.ParseFiles("pkg/web/template/signup.html")),
		errorTmpl:      template.Must(template.ParseFiles("pkg/web/template/error.html")),
	}
}
