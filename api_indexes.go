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
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/couchbaselabs/cbfullofit/search"
	"github.com/couchbaselabs/go-couchbase"
	"github.com/gorilla/mux"
)

func serveIndexesList(w http.ResponseWriter, r *http.Request) {

	indexes, err := indexList(*cbServ)

	if err != nil {
		showError(w, r, err.Error(), 500)
	}

	mustEncode(w, indexes)
}

type IndexByBucketViewResult struct {
	TotalRows int `json:"total_rows"`
	Rows      []struct {
		ID    string
		Key   []string
		Value interface{}
		Doc   *interface{}
	}
	Errors []couchbase.ViewError
}

func indexList(serv string) ([]string, error) {

	viewResult := IndexByBucketViewResult{}

	err := db.ViewCustom(ddoc, "indexesByBucket", map[string]interface{}{}, &viewResult)
	if err != nil {
		return nil, err
	}
	rv := make([]string, len(viewResult.Rows))
	for i, row := range viewResult.Rows {
		rv[i] = row.Key[1]
	}
	return rv, nil
}

type Index struct {
	Name   string           `json:"name"`
	Type   string           `json:"type"`
	Bucket string           `json:"bucket"`
	Schema map[string]Field `json:"schema"`
}

func createIndex(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	indexName := vars["index"]

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		showError(w, r, err.Error(), 500)
		return
	}

	index := Index{}
	err = json.Unmarshal(body, &index)
	if err != nil {
		showError(w, r, err.Error(), 400)
		return
	}

	// assert that type is index
	if index.Type != "index" {
		showError(w, r, "type must be 'index'", 400)
		return
	}

	// assert that name field in JSON matches URL
	if indexName != index.Name {
		showError(w, r, fmt.Sprintf("index name '%s' in URL does not match '%s' in request body", indexName, index.Name), 400)
		return
	}

	// assert that bucket exists
	_, bucketExists := db.GetPool().BucketMap[index.Bucket]
	if !bucketExists {
		showError(w, r, fmt.Sprintf("bucket '%s' does not exist", index.Bucket), 400)
		return
	}

	added, err := db.Add("index_"+indexName, 0, index)
	if err != nil {
		showError(w, r, err.Error(), 500)
		return
	}
	if !added {
		showError(w, r, fmt.Sprintf("index named '%s' already exists", indexName), 400)
		return
	}

	mustEncode(w, index)
}

func deleteIndex(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	indexName := vars["index"]

	err := db.Delete("index_" + indexName)
	if err != nil {
		showError(w, r, err.Error(), 500)
		return
	}

	mustEncode(w, "ok")
}

func getIndex(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	indexName := vars["index"]

	index, err := getIndexDoc(indexName)
	if err != nil {
		showError(w, r, err.Error(), 500)
		return
	}

	mustEncode(w, index)
}

func getIndexDoc(indexName string) (*Index, error) {
	index := Index{}
	err := db.Get("index_"+indexName, &index)
	if err != nil {
		return nil, err
	}
	return &index, nil
}

func searchIndexTerm(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	indexName := vars["index"]

	indexer, ok := assignments[indexName]
	if !ok {
		// FIXME, redirect to a node that can?
		showError(w, r, "sorry this node cannot search this index", 500)
		return
	}

	queryTerm := r.FormValue("q")
	queryField := r.FormValue("f")

	tq := search.TermQuery{
		Term:    queryTerm,
		Field:   queryField,
		Boost:   1.0,
		Explain: true,
	}
	log.Printf("query: %#v", tq)

	collector := search.NewTopScorerCollector(10)
	searcher, err := tq.Searcher(indexer.index)
	if err != nil {
		showError(w, r, fmt.Sprintf("searcher error: %v", err), 500)
		return
	}
	err = collector.Collect(searcher)
	if err != nil {
		showError(w, r, fmt.Sprintf("search error: %v", err), 500)
		return
	}
	results := collector.Results()

	fres := struct {
		MaxScore  float64                        `json:"max_score"`
		TotalHits uint64                         `json:"total_hits"`
		Took      float64                        `json:"took"`
		Hits      search.DocumentMatchCollection `json:"hits"`
	}{
		Hits:      results,
		MaxScore:  collector.MaxScore(),
		TotalHits: collector.Total(),
		Took:      collector.Took().Seconds(),
	}

	mustEncode(w, fres)
}

type SearchRequest struct {
	Q       search.Query `json:"query"`
	Size    float64      `json:"size"`
	Explain bool         `json:"explain"`
}

func (r *SearchRequest) UnmarshalJSON(input []byte) error {
	var temp struct {
		Q       json.RawMessage `json:"query"`
		Size    float64         `json:"size"`
		Explain bool            `json:"explain"`
	}

	err := json.Unmarshal(input, &temp)
	if err != nil {
		return err
	}

	r.Size = temp.Size
	r.Explain = temp.Explain
	r.Q, err = search.ParseQuery(temp.Q)
	if err != nil {
		return err
	}

	if r.Size <= 0 {
		r.Size = 10
	}

	return nil
}

func searchIndex(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	indexName := vars["index"]

	indexer, ok := assignments[indexName]
	if !ok {
		// FIXME, redirect to a node that can?
		showError(w, r, "sorry this node cannot search this index", 500)
		return
	}

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		showError(w, r, fmt.Sprintf("error reading request body: %v", err), 500)
		return
	}

	log.Printf("query: %s", string(requestBody))

	var sr SearchRequest
	err = json.Unmarshal(requestBody, &sr)
	if err != nil {
		showError(w, r, fmt.Sprintf("error parsing query: %v", err), 500)
		return
	}

	err = sr.Q.Validate()
	if err != nil {
		showError(w, r, fmt.Sprintf("error validating query: %v", err), 500)
		return
	}

	collector := search.NewTopScorerCollector(int(sr.Size))
	searcher, err := sr.Q.Searcher(indexer.index)
	if err != nil {
		showError(w, r, fmt.Sprintf("searcher error: %v", err), 500)
		return
	}
	err = collector.Collect(searcher)
	if err != nil {
		showError(w, r, fmt.Sprintf("search error: %v", err), 500)
		return
	}
	results := collector.Results()

	fres := struct {
		MaxScore  float64                        `json:"max_score"`
		TotalHits uint64                         `json:"total_hits"`
		Took      float64                        `json:"took"`
		Hits      search.DocumentMatchCollection `json:"hits"`
	}{
		Hits:      results,
		MaxScore:  collector.MaxScore(),
		TotalHits: collector.Total(),
		Took:      collector.Took().Seconds(),
	}

	mustEncode(w, fres)
}
