//  Copyright (c) 2013 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.
package upside_down

import (
	"os"
	"reflect"
	"testing"

	"github.com/couchbaselabs/cbfullofit/index"
)

func TestIndexReader(t *testing.T) {
	defer os.RemoveAll("test")

	schema := []*index.Field{
		&index.Field{
			Name:     "name",
			Path:     "/name",
			Analyzer: "standard",
		},
		&index.Field{
			Name:               "desc",
			Path:               "/description",
			Analyzer:           "standard",
			IncludeTermVectors: true,
		},
	}
	idx := NewUpsideDownCouch("test", schema)

	err := idx.Open()
	if err != nil {
		t.Errorf("error opening index: %v", err)
	}
	defer idx.Close()

	var expectedCount uint64 = 0
	doc := []byte(`{"name": "test"}`)
	err = idx.Update([]byte{'1'}, doc)
	if err != nil {
		t.Errorf("Error updating index: %v", err)
	}
	expectedCount += 1

	doc = []byte(`{"name": "test test test", "description": "eat more rice"}`)
	err = idx.Update([]byte{'2'}, doc)
	if err != nil {
		t.Errorf("Error updating index: %v", err)
	}
	expectedCount += 1

	// first look for a term that doesnt exist
	reader, err := idx.TermFieldReader([]byte("nope"), "name")
	if err != nil {
		t.Errorf("Error accessing term field reader: %v", err)
	}
	count := reader.Count()
	if count != 0 {
		t.Errorf("Expected doc count to be: %d got: %d", 0, count)
	}
	reader.Close()

	reader, err = idx.TermFieldReader([]byte("test"), "name")
	if err != nil {
		t.Errorf("Error accessing term field reader: %v", err)
	}
	defer reader.Close()

	expectedCount = 2
	count = reader.Count()
	if count != expectedCount {
		t.Errorf("Exptected doc count to be: %d got: %d", expectedCount, count)
	}

	var match *index.TermFieldDoc
	var actualCount uint64
	match, err = reader.Next()
	for err == nil && match != nil {
		match, err = reader.Next()
		if err != nil {
			t.Errorf("unexpected error reading next")
		}
		actualCount += 1
	}
	if actualCount != count {
		t.Errorf("count was 2, but only saw %d", actualCount)
	}

	expectedMatch := &index.TermFieldDoc{
		ID:   "2",
		Freq: 1,
		Norm: 0.5773502588272095,
		Vectors: []*index.TermFieldVector{
			&index.TermFieldVector{
				Field: "desc",
				Pos:   3,
				Start: 9,
				End:   13,
			},
		},
	}
	tfr, err := idx.TermFieldReader([]byte("rice"), "desc")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	match, err = tfr.Next()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(expectedMatch, match) {
		t.Errorf("got %#v, expected %#v", match, expectedMatch)
	}
}
