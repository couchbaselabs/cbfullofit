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
	"testing"

	"github.com/couchbaselabs/cbfullofit/index"
)

func TestIndexReader(t *testing.T) {
	defer os.RemoveAll("test")

	schema := []*index.Field{&index.Field{"name", "/name", "standard"}}
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

	doc = []byte(`{"name": "test test test"}`)
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

	reader, err = idx.TermFieldReader([]byte("test"), "name")
	if err != nil {
		t.Errorf("Error accessing term field reader: %v", err)
	}

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
}
