// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

//IncomingMsg :
type IncomingMsg struct {
	User        string
	Repo        string
	GithubToken string
	Msg         string
}

func bookmarkPost(w http.ResponseWriter, req *http.Request) {
	var in IncomingMsg

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println("Data read error:", err)
		return
	}
	log.Println("Body:", string(body))
	err = json.Unmarshal(body, &in)
	log.Println("Get request:", in)
	if err != nil {
		log.Println("json unmarkshal error:", err)
		return
	}

	//Pass parameter
	bm := NewBookmark(in.User, in.Repo, in.GithubToken)
	err = bm.SaveBookmark(in.Msg)
	if err != nil {
		log.Println("err=", err)
	}
	log.Println("Github issue post success! Msg=", in.Msg)
}

func serveHttpAPI(port string, existC chan bool) {
	go func() {
		if err, ok := <-existC; ok {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/", bookmarkPost)
	http.ListenAndServe(":"+port, mux)
}
