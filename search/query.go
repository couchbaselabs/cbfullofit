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
	"encoding/json"
	"fmt"

	"github.com/couchbaselabs/cbfullofit/index"
)

type Query interface {
	GetBoost() float64
	Searcher(index index.Index) (Searcher, error)
}

func ParseQuery(input []byte) (Query, error) {
	var tmp map[string]interface{}
	err := json.Unmarshal(input, &tmp)
	if err != nil {
		return nil, err
	}
	_, isTermQuery := tmp["term"]
	if isTermQuery {
		var rv *TermQuery
		err := json.Unmarshal(input, &rv)
		if err != nil {
			return nil, err
		}
		return rv, nil
	}
	_, hasMust := tmp["must"]
	_, hasShould := tmp["should"]
	if hasMust || hasShould {
		var rv *TermBooleanQuery
		err := json.Unmarshal(input, &rv)
		if err != nil {
			return nil, err
		}
		return rv, nil
	}
	return nil, fmt.Errorf("Unrecognized query")
}

type TermQuery struct {
	Term    string  `json:"term"`
	Field   string  `json:"field,omitempty"`
	Boost   float64 `json:"boost,omitempty"`
	Explain bool    `json:"explain,omitempty"`
}

func (q *TermQuery) GetBoost() float64 {
	return q.Boost
}

func (q *TermQuery) Searcher(index index.Index) (Searcher, error) {
	return NewTermSearcher(index, q)
}

type TermConjunctionQuery struct {
	Terms   []Query `json:"terms"`
	Boost   float64 `json:"boost"`
	Explain bool    `json:"explain"`
}

func (q *TermConjunctionQuery) UnmarshalJSON(input []byte) error {
	var temp struct {
		Terms   []json.RawMessage `json:"terms"`
		Boost   float64           `json:"boost"`
		Explain bool              `json:"explain"`
	}

	err := json.Unmarshal(input, &temp)
	if err != nil {
		return err
	}

	q.Boost = temp.Boost
	q.Explain = temp.Explain
	q.Terms = make([]Query, len(temp.Terms))
	for i, term := range temp.Terms {
		tq, err := ParseQuery(term)
		if err != nil {
			return err
		}
		q.Terms[i] = tq
	}
	return nil
}

func (q *TermConjunctionQuery) GetBoost() float64 {
	return q.Boost
}

func (q *TermConjunctionQuery) Searcher(index index.Index) (Searcher, error) {
	return NewTermConjunctionSearcher(index, q)
}

type TermDisjunctionQuery struct {
	Terms   []Query `json:"terms"`
	Boost   float64 `json:"boost"`
	Explain bool    `json:"explain"`
	Min     float64 `json:"min"`
}

func (q *TermDisjunctionQuery) UnmarshalJSON(input []byte) error {
	var temp struct {
		Terms   []json.RawMessage `json:"terms"`
		Boost   float64           `json:"boost"`
		Explain bool              `json:"explain"`
		Min     float64           `json:"min"`
	}

	err := json.Unmarshal(input, &temp)
	if err != nil {
		return err
	}

	q.Boost = temp.Boost
	q.Explain = temp.Explain
	q.Min = temp.Min
	q.Terms = make([]Query, len(temp.Terms))
	for i, term := range temp.Terms {
		tq, err := ParseQuery(term)
		if err != nil {
			return err
		}
		q.Terms[i] = tq
	}
	return nil
}

func (q *TermDisjunctionQuery) GetBoost() float64 {
	return q.Boost
}

func (q *TermDisjunctionQuery) Searcher(index index.Index) (Searcher, error) {
	return NewTermDisjunctionSearcher(index, q)
}

type TermBooleanQuery struct {
	Must    *TermConjunctionQuery `json:"must,omitempty"`
	MustNot *TermDisjunctionQuery `json:"must_not,omitempty"`
	Should  *TermDisjunctionQuery `json:"should,omitempty"`
	Boost   float64               `json:"boost,omitempty"`
	Explain bool                  `json:"explain,omitempty"`
}

func (q *TermBooleanQuery) GetBoost() float64 {
	return q.Boost
}

func (q *TermBooleanQuery) Searcher(index index.Index) (Searcher, error) {
	return NewTermBooleanSearcher(index, q)
}
