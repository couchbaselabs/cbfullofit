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
	"reflect"
	"testing"
)

func TestParseQuery(t *testing.T) {
	tests := []struct {
		input []byte
		query Query
	}{
		{
			input: []byte(`{"term":"test","field":"desc","boost":1.0,"explain":true}`),
			query: &TermQuery{
				Term:    "test",
				Field:   "desc",
				Boost:   1.0,
				Explain: true,
			},
		},
		{
			input: []byte(`{"must":{"terms":[{"term":"test_must","field":"desc","boost":1.0,"explain":true}],"boost":1.0,"explain":true},"must_not":{"terms":[{"term":"test_must_not","field":"desc","boost":1.0,"explain":true}],"boost":1.0,"explain":true},"should":{"terms":[{"term":"test_should","field":"desc","boost":1.0,"explain":true}],"boost":1.0,"explain":true,"min":1.0},"boost":1.0,"explain":true}`),
			query: &TermBooleanQuery{
				Must: &TermConjunctionQuery{
					Terms: []Query{
						&TermQuery{
							Term:    "test_must",
							Field:   "desc",
							Boost:   1.0,
							Explain: true,
						},
					},
					Boost:   1.0,
					Explain: true,
				},
				Should: &TermDisjunctionQuery{
					Terms: []Query{
						&TermQuery{
							Term:    "test_should",
							Field:   "desc",
							Boost:   1.0,
							Explain: true,
						},
					},
					Boost:   1.0,
					Explain: true,
					Min:     1.0,
				},
				MustNot: &TermDisjunctionQuery{
					Terms: []Query{
						&TermQuery{
							Term:    "test_must_not",
							Field:   "desc",
							Boost:   1.0,
							Explain: true,
						},
					},
					Boost:   1.0,
					Explain: true,
				},
				Boost:   1.0,
				Explain: true,
			},
		},
	}

	for _, test := range tests {
		q, err := ParseQuery(test.input)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(q, test.query) {
			t.Errorf("expected %v got %v for %s", test.query, q, string(test.input))
			qbytes, err := json.MarshalIndent(&q, "", "  ")
			if err == nil {
				t.Logf("q in json is: %s", string(qbytes))
			}

			q2bytes, err := json.MarshalIndent(&test.query, "", "  ")
			if err == nil {
				t.Logf("q2 in json is: %s", string(q2bytes))
			}
		}
	}
}
