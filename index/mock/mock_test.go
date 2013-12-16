//  Copyright (c) 2013 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.
package mock

import (
	"reflect"
	"testing"

	_ "github.com/couchbaselabs/cbfullofit/analysis/analyzers/standard_analyzer"
	"github.com/couchbaselabs/cbfullofit/index"
)

func TestCRUD(t *testing.T) {
	schema := []*index.Field{
		&index.Field{
			Name:     "name",
			Path:     "/name",
			Analyzer: "standard",
		},
		&index.Field{
			Name:     "desc",
			Path:     "/description",
			Analyzer: "standard",
		},
	}
	i := NewMockIndex(schema)

	// create doc, assert doc count goes up
	i.Update([]byte("1"), []byte(`{"name": "marty"}`))
	count := i.DocCount()
	if count != 1 {
		t.Errorf("expected document count to be 1, was: %d", count)
	}

	// add another doc, assert doc count goes up again
	i.Update([]byte("2"), []byte(`{"name": "bob"}`))
	count = i.DocCount()
	if count != 2 {
		t.Errorf("expected document count to be 2, was: %d", count)
	}

	// search for doc with term that should exist
	expectedMatch := &index.TermFieldDoc{
		ID:   "1",
		Freq: 1,
		Norm: 1,
	}
	tfr, err := i.TermFieldReader([]byte("marty"), "name")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	match, err := tfr.Next()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(expectedMatch, match) {
		t.Errorf("got %v, expected %v", match, expectedMatch)
	}
	nomatch, err := tfr.Next()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if nomatch != nil {
		t.Errorf("expected nil after last match")
	}

	// update doc, assert doc count doesn't go up
	i.Update([]byte("1"), []byte(`{"name": "salad"}`))
	count = i.DocCount()
	if count != 2 {
		t.Errorf("expected document count to be 2, was: %d", count)
	}

	// perform the original search again, should NOT find anything this time
	tfr, err = i.TermFieldReader([]byte("marty"), "name")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	nomatch, err = tfr.Next()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if nomatch != nil {
		t.Errorf("expected no matches, found one")
		t.Logf("%v", i)
	}

	// delete a doc, ensure the count is 1
	err = i.Delete([]byte("2"))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	count = i.DocCount()
	if count != 1 {
		t.Errorf("expected document count to be 1, was: %d", count)
	}
}
