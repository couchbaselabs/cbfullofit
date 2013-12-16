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
	"log"
	"time"
)

var assignments = make(map[string]*Indexer)

func lookupAssignments(assignments map[string]*Indexer) {

	log.Printf("checking for assignments")

	indexes, err := indexesAssignedToNode(nodeID)
	log.Printf("assignments: %v", indexes)
	if err != nil {
		log.Printf("Unable to lookup assignments: %v", err)
		return
	}

	// look for new assignments
	for _, indexName := range indexes {
		indexer, ok := assignments[indexName]
		if !ok {
			log.Printf("starting new indexer for '%s'", indexName)
			// start up an indexer
			index, err := getIndexDoc(indexName)
			if err != nil {
				log.Printf("cannot find index '%s' in assignment", indexName)
				continue
			}

			indexer = NewIndexer(index.Name, index.Bucket, index.Schema)
			assignments[indexName] = indexer
			go indexer.Run()
		}
	}

	// look for removed assignments
OUTER:
	for assignedIndex, indexer := range assignments {
		for _, index := range indexes {
			if assignedIndex == index {
				continue OUTER
			}
		}
		log.Printf("stopping indexer for '%s'", assignedIndex)
		// stop an indexer
		indexer.Stop()
		delete(assignments, assignedIndex)
	}
}

func pollAssignments() {

	lookupAssignments(assignments)

	period := time.Second * 60
	ticker := time.NewTicker(period)

	for {
		select {
		case <-ticker.C:
			lookupAssignments(assignments)
		}
	}
}
