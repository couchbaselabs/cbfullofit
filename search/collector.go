//  Copyright (c) 2013 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.
package search

import (
	"container/list"
	"time"
)

type Collector interface {
	Collect(searcher Searcher) error
	Results() DocumentMatchCollection
	Total() uint64
	MaxScore() float64
	Took() time.Duration
}

type TopScoreCollector struct {
	k        int
	results  *list.List
	took     time.Duration
	maxScore float64
	total    uint64
}

func NewTopScorerCollector(k int) *TopScoreCollector {
	return &TopScoreCollector{
		k:       k,
		results: list.New(),
	}
}

func (tksc *TopScoreCollector) Total() uint64 {
	return tksc.total
}

func (tksc *TopScoreCollector) MaxScore() float64 {
	return tksc.maxScore
}

func (tksc *TopScoreCollector) Took() time.Duration {
	return tksc.took
}

func (tksc *TopScoreCollector) Collect(searcher Searcher) error {
	startTime := time.Now()
	next, err := searcher.Next()
	for err == nil && next != nil {
		tksc.collectSingle(next)
		next, err = searcher.Next()
	}
	// compute search duration
	tksc.took = time.Since(startTime)
	if err != nil {
		return err
	}
	return nil
}

func (tksc *TopScoreCollector) collectSingle(dm *DocumentMatch) {
	// increment total hits
	tksc.total += 1

	// update max score
	if dm.Score > tksc.maxScore {
		tksc.maxScore = dm.Score
	}

	for e := tksc.results.Front(); e != nil; e = e.Next() {
		curr := e.Value.(*DocumentMatch)
		if dm.Score < curr.Score {

			tksc.results.InsertBefore(dm, e)
			// if we just made the list too long
			if tksc.results.Len() > tksc.k {
				// remove the head
				tksc.results.Remove(tksc.results.Front())
			}
			return
		}
	}
	// if we got to the end, we still have to add it
	tksc.results.PushBack(dm)
	if tksc.results.Len() > tksc.k {
		// remove the head
		tksc.results.Remove(tksc.results.Front())
	}
}

func (tksc *TopScoreCollector) Results() DocumentMatchCollection {
	rv := make(DocumentMatchCollection, tksc.results.Len())
	i := 0
	for e := tksc.results.Back(); e != nil; e = e.Prev() {
		rv[i] = e.Value.(*DocumentMatch)
		i++
	}
	return rv
}
