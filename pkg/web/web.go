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
)

type Web struct {
	indexTmpl      *template.Template
	signupPageTmpl *template.Template
	loginPageTmpl  *template.Template
}

func NewWeb() *Web {
	return &Web{
		indexTmpl:      template.Must(template.ParseFiles("pkg/web/templates/index.html")),
		signupPageTmpl: template.Must(template.ParseFiles("pkg/web/templates/signup.html")),
		loginPageTmpl:  template.Must(template.ParseFiles("pkg/web/templates/login.html")),
	}
}