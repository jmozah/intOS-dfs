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

package cookie

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/securecookie"
)

const (
	cookieName           = "intOS-dfs"
	cookieUserName       = "user"
	cookieSessionId      = "cookie-id"
	cookiePodName        = "pod"
	cookieExpirationTime = 15 * time.Minute
)

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func GetUniqueSessionId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func SetSession(userName, sessionId string, response http.ResponseWriter) error {
	value := map[string]string{
		cookieUserName:  userName,
		cookieSessionId: sessionId,
	}
	encoded, err := cookieHandler.Encode(cookieName, value)
	if err != nil {
		return err
	}
	expire := time.Now().Add(cookieExpirationTime)
	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    encoded,
		Path:     "/",
		Expires:  expire,
		HttpOnly: true,
		MaxAge:   0, // to make sure that the browser does not persist it in disk
	}
	http.SetCookie(response, cookie)
	return nil
}

func ResetSessionExpiry(request *http.Request, response http.ResponseWriter) error {
	rcvdCookie, err := request.Cookie(cookieName)
	if err != nil {
		return err
	}
	expire := time.Now().Add(cookieExpirationTime)
	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    rcvdCookie.Value,
		Path:     "/",
		Expires:  expire,
		HttpOnly: true,
		MaxAge:   0, // to make sure that the browser does not persist it in disk
	}
	http.SetCookie(response, cookie)
	return nil
}

func GetUserNameAndSessionId(request *http.Request) (userName, sessionId string, err error) {
	cookie, err := request.Cookie(cookieName)
	if err != nil {
		return "", "", err
	}
	cookieValue := make(map[string]string)
	err = cookieHandler.Decode(cookieName, cookie.Value, &cookieValue)
	if err != nil {
		return "", "", err
	}
	userName = cookieValue[cookieUserName]
	sessionId = cookieValue[cookieSessionId]
	return userName, sessionId, nil
}

func GetUserNameSessionIdAndPodName(request *http.Request) (userName, sessionId, podName string, err error) {
	cookie, err := request.Cookie(cookieName)
	if err != nil {
		return "", "", "", err
	}
	cookieValue := make(map[string]string)
	err = cookieHandler.Decode(cookieName, cookie.Value, &cookieValue)
	if err != nil {
		return "", "", "", err
	}
	userName = cookieValue[cookieUserName]
	sessionId = cookieValue[cookieSessionId]
	podName = cookieValue[cookiePodName]
	return userName, sessionId, podName, nil
}

func ClearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Expires:  time.Now().Add(-time.Duration(1) * time.Second), // set the expiry to 1 second
	}
	http.SetCookie(response, cookie)
}

func SetPodNameInSession(podName string, request *http.Request, response http.ResponseWriter) error {
	rcvdCookie, err := request.Cookie(cookieName)
	if err != nil {
		return err
	}
	cookieValue := make(map[string]string)
	err = cookieHandler.Decode(cookieName, rcvdCookie.Value, &cookieValue)
	if err != nil {
		return err
	}
	cookieValue[cookiePodName] = podName
	encoded, err := cookieHandler.Encode(cookieName, cookieValue)
	if err != nil {
		return err
	}
	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    encoded,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   0, // to make sure that the browser does not persist it in disk
	}
	http.SetCookie(response, cookie)
	return nil
}

func RemovePodNameFromSession(request *http.Request, response http.ResponseWriter) error {
	rcvdCookie, err := request.Cookie(cookieName)
	if err != nil {
		return err
	}
	cookieValue := make(map[string]string)
	err = cookieHandler.Decode(cookieName, rcvdCookie.Value, &cookieValue)
	if err != nil {
		return err
	}
	delete(cookieValue, cookiePodName)
	encoded, err := cookieHandler.Encode(cookieName, cookieValue)
	if err != nil {
		return err
	}
	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    encoded,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   0, // to make sure that the browser does not persist it in disk
	}
	http.SetCookie(response, cookie)
	return nil
}
