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
	"fmt"
	"net/http"

	"github.com/couchbaselabs/go-couchbase"
	"github.com/gorilla/mux"
)

func indexAssignmentList(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	indexName := vars["index"]

	buckets, err := nodesAssignedToIndex(indexName)

	if err != nil {
		showError(w, r, err.Error(), 500)
	}

	mustEncode(w, buckets)
}

func nodesAssignedToIndex(index string) ([]string, error) {
	viewResult := couchbase.ViewResult{}

	err := db.ViewCustom(ddoc, "assignmentsByIndex", map[string]interface{}{"key": index, "stale": false}, &viewResult)
	if err != nil {
		return nil, err
	}
	rv := make([]string, len(viewResult.Rows))
	for i, row := range viewResult.Rows {
		rv[i] = row.Value.(string)
	}
	return rv, nil
}

type Assignment struct {
	Index string `json:"index"`
	Node  string `json:"node"`
	Type  string `json:"type"`
}

func assignNodeToIndex(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	indexName := vars["index"]
	nodeName := vars["node"]

	// FIXME add validation here

	assignment := Assignment{
		Index: indexName,
		Node:  nodeName,
		Type:  "assignment",
	}
	added, err := db.Add("assignment_"+nodeName+indexName, 0, assignment)
	if err != nil {
		showError(w, r, err.Error(), 500)
		return
	}
	if !added {
		showError(w, r, fmt.Sprintf("assigment of '%s' to '%s' already exists", indexName, nodeName), 400)
		return
	}

	mustEncode(w, assignment)
}

func unAssignNodeToIndex(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	indexName := vars["index"]
	nodeName := vars["node"]

	// FIXME add validation here

	err := db.Delete("assignment_" + nodeName + indexName)
	if err != nil {
		showError(w, r, err.Error(), 500)
		return
	}

	mustEncode(w, "ok")
}

func nodeAssignmentList(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	nodeName := vars["node"]

	buckets, err := indexesAssignedToNode(nodeName)

	if err != nil {
		showError(w, r, err.Error(), 500)
	}

	mustEncode(w, buckets)
}

func indexesAssignedToNode(node string) ([]string, error) {
	viewResult := couchbase.ViewResult{}

	err := db.ViewCustom(ddoc, "assignmentsByNode", map[string]interface{}{"key": node, "stale": false}, &viewResult)
	if err != nil {
		return nil, err
	}
	rv := make([]string, len(viewResult.Rows))
	for i, row := range viewResult.Rows {
		rv[i] = row.Value.(string)
	}
	return rv, nil
}
