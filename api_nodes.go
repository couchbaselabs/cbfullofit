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
	"net/http"

	"github.com/couchbaselabs/go-couchbase"
)

func serveNodesList(w http.ResponseWriter, r *http.Request) {

	nodes, err := nodeList(*cbServ)

	if err != nil {
		showError(w, r, err.Error(), 500)
	}

	mustEncode(w, nodes)
}

type NodesViewResult struct {
	TotalRows int `json:"total_rows"`
	Rows      []struct {
		ID    string
		Key   string
		Value interface{}
		Doc   *interface{}
	}
	Errors []couchbase.ViewError
}

func nodeList(serv string) ([]string, error) {

	viewResult := NodesViewResult{}

	err := db.ViewCustom(ddoc, "nodes", map[string]interface{}{}, &viewResult)
	if err != nil {
		return nil, err
	}
	rv := make([]string, len(viewResult.Rows))
	for i, row := range viewResult.Rows {
		rv[i] = row.Key
	}
	return rv, nil
}
