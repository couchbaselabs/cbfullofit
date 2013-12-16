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
	"testing"

	"github.com/couchbaselabs/cbfullofit/index"
)

func TestTermSearch(t *testing.T) {

	tests := []struct {
		index   index.Index
		query   *TermQuery
		results []*DocumentMatch
	}{
		{
			index: twoDocIndex,
			query: &TermQuery{
				Term:    "beer",
				Field:   "desc",
				Boost:   1.0,
				Explain: true,
			},
			results: []*DocumentMatch{
				&DocumentMatch{
					ID:    "1",
					Score: 1.0,
				},
				&DocumentMatch{
					ID:    "2",
					Score: 0.5,
				},
				&DocumentMatch{
					ID:    "3",
					Score: 0.5,
				},
				&DocumentMatch{
					ID:    "4",
					Score: 1.0,
				},
			},
		},
		{
			index: twoDocIndex,
			query: &TermQuery{
				Term:    "marty",
				Field:   "name",
				Boost:   1.0,
				Explain: true,
			},
			results: []*DocumentMatch{
				&DocumentMatch{
					ID:    "1",
					Score: 1.916290731874155,
				},
			},
		},
	}

	for testIndex, test := range tests {
		searcher, err := NewTermSearcher(test.index, test.query)
		defer searcher.Close()

		next, err := searcher.Next()
		i := 0
		for err == nil && next != nil {
			if i < len(test.results) {
				if next.ID != test.results[i].ID {
					t.Errorf("expected result %d to have id %s got %s for test %d", i, test.results[i].ID, next.ID, testIndex)
				}
				if next.Score != test.results[i].Score {
					t.Errorf("expected result %d to have score %v got  %v for test %d", i, test.results[i].Score, next.Score, testIndex)
					t.Logf("explanation: %v", next.Expl)
				}
			}
			next, err = searcher.Next()
			i++
		}
		if err != nil {
			t.Fatalf("error iterating searcher: %v for test %d", err, testIndex)
		}
		if len(test.results) != i {
			t.Errorf("expected %d results got %d for test %d", len(test.results), i, testIndex)
		}
	}
}
