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
	"log"
	"time"

	"github.com/couchbaselabs/go-couchbase"
)

type viewMarker struct {
	Version   int       `json:"version"`
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"`
}

const ddoc = "cbfullofit"
const ddocKey = "/@cbfullofitDdocVersion"
const ddocVersion = 4
const designDoc = `
{
  "views": {
  	"indexesByBucket": {
  		"map": "function (doc, meta) { if (doc.type == 'index') { emit([doc.bucket,doc.name], null);}}"
  	},
	"nodes": {
  		"map": "function (doc, meta) { if (doc.type == 'node') { emit(doc.name, null);}}"
  	},
	"assignmentsByIndex": {
  		"map": "function (doc, meta) { if (doc.type == 'assignment') { emit(doc.index, doc.node);}}"
  	},
	"assignmentsByNode": {
  		"map": "function (doc, meta) { if (doc.type == 'assignment') { emit(doc.node, doc.index);}}"
  	}
  }
}`

func dbConnect(serv, pool, bucket string) (*couchbase.Bucket, error) {

	log.Printf("Connecting to couchbase bucket %v at %v",
		bucket, serv)
	rv, err := couchbase.GetBucket(serv, pool, bucket)
	if err != nil {
		return nil, err
	}

	marker := viewMarker{}
	err = rv.Get(ddocKey, &marker)
	if err != nil {
		log.Printf("Error checking view version: %v", err)
	}
	if marker.Version < ddocVersion {
		log.Printf("Installing new version of views (old version=%v)",
			marker.Version)
		doc := json.RawMessage([]byte(designDoc))
		err = rv.PutDDoc(ddoc, &doc)
		if err != nil {
			return nil, err
		}
		marker.Version = ddocVersion
		marker.Timestamp = time.Now().UTC()
		marker.Type = "ddocmarker"

		rv.Set(ddocKey, 0, &marker)
	}

	return rv, nil
}
