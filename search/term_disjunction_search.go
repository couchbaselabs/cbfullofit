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
	"math"
	"sort"

	"github.com/couchbaselabs/cbfullofit/index"
)

type TermDisjunctionSearcher struct {
	index     index.Index
	searchers OrderedSearcherList
	queryNorm float64
	currs     []*DocumentMatch
	currentId string
	scorer    *TermDisjunctionQueryScorer
	min       float64
}

func NewTermDisjunctionSearcher(index index.Index, query *TermDisjunctionQuery) (*TermDisjunctionSearcher, error) {
	// build the downstream searchres
	searchers := make(OrderedSearcherList, len(query.Terms))
	for i, termQuery := range query.Terms {
		searcher, err := termQuery.Searcher(index)
		if err != nil {
			return nil, err
		}
		searchers[i] = searcher
	}
	// sort the searchers
	sort.Sort(sort.Reverse(searchers))
	// build our searcher
	rv := TermDisjunctionSearcher{
		index:     index,
		searchers: searchers,
		currs:     make([]*DocumentMatch, len(searchers)),
		scorer:    NewTermDisjunctionQueryScorer(query.Explain),
		min:       query.Min,
	}
	rv.computeQueryNorm()
	err := rv.initSearchers()
	if err != nil {
		return nil, err
	}

	return &rv, nil
}

func (s *TermDisjunctionSearcher) computeQueryNorm() {
	// first calculate sum of squared weights
	sumOfSquaredWeights := 0.0
	for _, termSearcher := range s.searchers {
		sumOfSquaredWeights += termSearcher.Weight()
	}
	// now compute query norm from this
	s.queryNorm = 1.0 / math.Sqrt(sumOfSquaredWeights)
	// finally tell all the downsteam searchers the norm
	for _, termSearcher := range s.searchers {
		termSearcher.SetQueryNorm(s.queryNorm)
	}
}

func (s *TermDisjunctionSearcher) initSearchers() error {
	var err error
	// get all searchers pointing at their first match
	for i, termSearcher := range s.searchers {
		s.currs[i], err = termSearcher.Next()
		if err != nil {
			return err
		}
	}

	s.currentId = s.nextSmallestId()
	return nil
}

func (s *TermDisjunctionSearcher) nextSmallestId() string {
	rv := ""
	for _, curr := range s.currs {
		if curr != nil && (curr.ID < rv || rv == "") {
			rv = curr.ID
		}
	}
	return rv
}

func (s *TermDisjunctionSearcher) Weight() float64 {
	var rv float64
	for _, searcher := range s.searchers {
		rv += searcher.Weight()
	}
	return rv
}

func (s *TermDisjunctionSearcher) SetQueryNorm(qnorm float64) {
	for _, searcher := range s.searchers {
		searcher.SetQueryNorm(qnorm)
	}
}

func (s *TermDisjunctionSearcher) Next() (*DocumentMatch, error) {
	var err error
	var rv *DocumentMatch
	matching := make([]*DocumentMatch, 0)

	found := false
	for !found && s.currentId != "" {
		for _, curr := range s.currs {
			if curr != nil && curr.ID == s.currentId {
				matching = append(matching, curr)
			}
		}

		if len(matching) > int(s.min) {
			found = true
			// score this match
			rv = s.scorer.Score(matching, len(matching), len(s.searchers))
		}

		// invoke next on all the matching searchers
		for i, curr := range s.currs {
			if curr != nil && curr.ID == s.currentId {
				searcher := s.searchers[i]
				s.currs[i], err = searcher.Next()
				if err != nil {
					return nil, err
				}
			}
		}
		s.currentId = s.nextSmallestId()
	}
	return rv, nil
}

func (s *TermDisjunctionSearcher) Advance(ID string) (*DocumentMatch, error) {

	// get all searchers pointing at their first match
	var err error
	for i, termSearcher := range s.searchers {
		s.currs[i], err = termSearcher.Advance(ID)
		if err != nil {
			return nil, err
		}
	}

	s.currentId = s.nextSmallestId()

	return s.Next()
}

func (s *TermDisjunctionSearcher) Count() uint64 {
	// for now return a worst case
	var sum uint64 = 0
	for _, searcher := range s.searchers {
		sum += searcher.Count()
	}
	return sum
}

func (s *TermDisjunctionSearcher) Close() {
	for _, searcher := range s.searchers {
		searcher.Close()
	}
}
