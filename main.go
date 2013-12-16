//  Copyright (c) 2013 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.
package main

import (
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/couchbaselabs/go-couchbase"
	"github.com/gorilla/mux"

	"github.com/nu7hatch/gouuid"

	_ "github.com/couchbaselabs/cbfullofit/analysis/analyzers/standard_analyzer"
)

var VERSION = "0.0.0"
var db *couchbase.Bucket
var nodeID string

var staticEtag = flag.String("staticEtag", "", "A static etag value.")
var staticPath = flag.String("static", "static", "Path to the static content")

var bindAddr = flag.String("addr", ":8094", "http listen address")

var cbServ = flag.String("couchbase", "http://localhost:8091/", "URL to couchbase")
var cbPool = flag.String("pool", "default", "couchbase pool")
var cbBucket = flag.String("bucket", "cbfullofit", "couchbase bucket")

var dataDir = flag.String("datadir", "data", "data storage directory")

var dump = flag.Bool("dump", false, "dump index contents")

type myFileHandler struct {
	h http.Handler
}

func (mfh myFileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if *staticEtag != "" {
		w.Header().Set("Etag", *staticEtag)
	}
	mfh.h.ServeHTTP(w, r)
}

func RewriteURL(to string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = to
		h.ServeHTTP(w, r)
	})
}

func getOrGenerateNodeID() string {
	nidBytes, err := ioutil.ReadFile(*dataDir + "/NODEID")
	if err != nil {
		log.Printf("NodeID not found, generating new NodeID")
		uuid, err := uuid.NewV4()
		if err != nil {
			log.Fatalf("Unable to generate nodeID: %v", err)
		}

		err = ioutil.WriteFile(*dataDir+"/NODEID", []byte(uuid.String()), os.ModePerm)
		if err != nil {
			log.Fatalf("Unable to record nodeID to disk: %v", err)
		}
		nidBytes = []byte(uuid.String())
	}
	return string(nidBytes)
}

func main() {

	flag.Parse()

	// ensure the data directory exists
	err := os.MkdirAll(*dataDir, os.ModePerm)
	if err != nil {
		log.Fatalf("error making data directory")
	}

	// find my nodeID or generate a new one if it doesn't exist
	nodeID = getOrGenerateNodeID()
	log.Printf("NodeID: %s", nodeID)

	// connect to the database
	//	var err error
	db, err = dbConnect(*cbServ, *cbPool, *cbBucket)
	if err != nil {
		log.Fatalf("Error connecting to couchbase: %v", err)
	}

	// start the heartbeat process
	go heartbeat()

	// start polling for assignments
	go pollAssignments()

	// start server
	r := mux.NewRouter()

	// API
	r.HandleFunc("/api/bucket/", serveBucketsList).Methods("GET")
	r.HandleFunc("/api/index/", serveIndexesList).Methods("GET")
	r.HandleFunc("/api/index/{index}", getIndex).Methods("GET")
	r.HandleFunc("/api/index/{index}", createIndex).Methods("PUT")
	r.HandleFunc("/api/index/{index}", deleteIndex).Methods("DELETE")
	r.HandleFunc("/api/index/{index}/_searchTerm", searchIndexTerm).Methods("GET")
	r.HandleFunc("/api/index/{index}/_search", searchIndex).Methods("POST")
	//r.HandleFunc("/api/index/{index}/_searchAllTerms", searchIndexAllTerms).Methods("GET")
	r.HandleFunc("/api/node/", serveNodesList).Methods("GET")

	// node/index assignment
	r.HandleFunc("/api/assignment/index/{index}", indexAssignmentList).Methods("GET")
	r.HandleFunc("/api/assignment/index/{index}/{node}", assignNodeToIndex).Methods("PUT")
	r.HandleFunc("/api/assignment/index/{index}/{node}", unAssignNodeToIndex).Methods("DELETE")
	r.HandleFunc("/api/assignment/node/{node}", nodeAssignmentList).Methods("GET")
	r.HandleFunc("/api/assignment/node/{node}/{index}", assignNodeToIndex).Methods("PUT")
	r.HandleFunc("/api/assignment/node/{node}/{index}", unAssignNodeToIndex).Methods("DELETE")

	// static
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		myFileHandler{http.FileServer(http.Dir(*staticPath))}))

	// application pages
	appPages := []string{
		"/home",
		"/index",
		"/node",
	}

	for _, p := range appPages {
		r.PathPrefix(p).Handler(RewriteURL("app.html",
			http.FileServer(http.Dir(*staticPath))))
	}

	r.Handle("/", http.RedirectHandler("/static/app.html", 302))

	http.Handle("/", r)

	log.Printf("Listening on %v", *bindAddr)
	log.Fatal(http.ListenAndServe(*bindAddr, nil))
}

func showError(w http.ResponseWriter, r *http.Request,
	msg string, code int) {
	log.Printf("Reporting error %v/%v", code, msg)
	http.Error(w, msg, code)
}

func mustEncode(w io.Writer, i interface{}) {
	if headered, ok := w.(http.ResponseWriter); ok {
		headered.Header().Set("Cache-Control", "no-cache")
		headered.Header().Set("Content-type", "application/json")
	}

	e := json.NewEncoder(w)
	if err := e.Encode(i); err != nil {
		panic(err)
	}
}
