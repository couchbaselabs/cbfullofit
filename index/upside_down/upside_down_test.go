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

	_ "github.com/couchbaselabs/cbfullofit/analysis/analyzers/standard_analyzer"
	"github.com/couchbaselabs/cbfullofit/index"
)

func TestIndexOpenReopen(t *testing.T) {
	defer os.RemoveAll("test")

	schema := []*index.Field{
		&index.Field{
			Name:     "name",
			Path:     "/name",
			Analyzer: "standard",
		},
	}
	idx := NewUpsideDownCouch("test", schema)

	err := idx.Open()
	if err != nil {
		t.Errorf("error opening index: %v", err)
	}

	var expectedCount uint64 = 0
	docCount := idx.DocCount()
	if docCount != expectedCount {
		t.Errorf("Expected document count to be %d got %d", expectedCount, docCount)
	}

	// opening database should have inserted version/schema
	expectedLength := uint64(1 + len(schema))
	rowCount := idx.rowCount()
	if rowCount != expectedLength {
		t.Errorf("expected %d rows, got: %d", expectedLength, rowCount)
	}

	// now close it
	idx.Close()

	// create a new instance with a different schema
	newSchema := []*index.Field{
		&index.Field{
			Name:     "desc",
			Path:     "/desc",
			Analyzer: "keyword",
		},
	}
	idx = NewUpsideDownCouch("test", newSchema)
	err = idx.Open()
	if err != nil {
		t.Errorf("error opening index: %v", err)
	}

	// schema SHOULD be the original schema, and NOT the new one
	if !reflect.DeepEqual(idx.schema, schema) {
		t.Errorf("wrong schema, expected: %v got: %v", schema, idx.schema)
	}

	// now close it
	idx.Close()
}

func TestIndexInsert(t *testing.T) {
	defer os.RemoveAll("test")

	schema := []*index.Field{
		&index.Field{
			Name:     "name",
			Path:     "/name",
			Analyzer: "standard",
		},
	}
	idx := NewUpsideDownCouch("test", schema)

	err := idx.Open()
	if err != nil {
		t.Errorf("error opening index: %v", err)
	}
	defer idx.Close()

	var expectedCount uint64 = 0
	docCount := idx.DocCount()
	if docCount != expectedCount {
		t.Errorf("Expected document count to be %d got %d", expectedCount, docCount)
	}

	doc := []byte(`{"name": "test"}`)
	err = idx.Update([]byte{'1'}, doc)
	if err != nil {
		t.Errorf("Error updating index: %v", err)
	}
	expectedCount += 1

	docCount = idx.DocCount()
	if docCount != expectedCount {
		t.Errorf("Expected document count to be %d got %d", expectedCount, docCount)
	}

	// should have 4 rows (1 for version, 1 for schema field, and 1 for single term, and 1 for the term count,  and 1 for the back index entry)
	expectedLength := uint64(1 + len(schema) + 1 + 1 + 1)
	rowCount := idx.rowCount()
	if rowCount != expectedLength {
		t.Errorf("expected %d rows, got: %d", expectedLength, rowCount)
	}
}

func TestIndexInsertThenDelete(t *testing.T) {
	defer os.RemoveAll("test")

	schema := []*index.Field{
		&index.Field{
			Name:     "name",
			Path:     "/name",
			Analyzer: "standard",
		},
	}
	idx := NewUpsideDownCouch("test", schema)

	err := idx.Open()
	if err != nil {
		t.Errorf("error opening index: %v", err)
	}
	defer idx.Close()

	var expectedCount uint64 = 0
	docCount := idx.DocCount()
	if docCount != expectedCount {
		t.Errorf("Expected document count to be %d got %d", expectedCount, docCount)
	}

	doc := []byte(`{"name": "test"}`)
	err = idx.Update([]byte{'1'}, doc)
	if err != nil {
		t.Errorf("Error updating index: %v", err)
	}
	expectedCount += 1

	docCount = idx.DocCount()
	if docCount != expectedCount {
		t.Errorf("Expected document count to be %d got %d", expectedCount, docCount)
	}

	err = idx.Delete([]byte{'1'})
	if err != nil {
		t.Errorf("Error deleting entry from index: %v", err)
	}
	expectedCount -= 1

	docCount = idx.DocCount()
	if docCount != expectedCount {
		t.Errorf("Expected document count to be %d got %d", expectedCount, docCount)
	}

	// should have 2 row (1 for version, 1 for schema field)
	expectedLength := uint64(1 + len(schema))
	rowCount := idx.rowCount()
	if rowCount != expectedLength {
		t.Errorf("expected %d rows, got: %d", expectedLength, rowCount)
	}
}

func TestIndexInsertThenUpdate(t *testing.T) {
	defer os.RemoveAll("test")

	schema := []*index.Field{
		&index.Field{
			Name:     "name",
			Path:     "/name",
			Analyzer: "standard",
		},
	}
	idx := NewUpsideDownCouch("test", schema)

	err := idx.Open()
	if err != nil {
		t.Errorf("error opening index: %v", err)
	}
	defer idx.Close()

	doc := []byte(`{"name": "test"}`)
	err = idx.Update([]byte{'1'}, doc)
	if err != nil {
		t.Errorf("Error updating index: %v", err)
	}

	// this update should overwrite one term, and introduce one new one
	doc = []byte(`{"name": "test fail"}`)
	err = idx.Update([]byte{'1'}, doc)
	if err != nil {
		t.Errorf("Error deleting entry from index: %v", err)
	}

	// should have 2 row (1 for version, 1 for schema field, and 2 for the two term, and 2 for the term counts, and 1 for the back index entry)
	expectedLength := uint64(1 + len(schema) + 2 + 2 + 1)
	rowCount := idx.rowCount()
	if rowCount != expectedLength {
		t.Errorf("expected %d rows, got: %d", expectedLength, rowCount)
	}

	// now do another update that should remove one of term
	doc = []byte(`{"name": "fail"}`)
	err = idx.Update([]byte{'1'}, doc)
	if err != nil {
		t.Errorf("Error deleting entry from index: %v", err)
	}

	// should have 2 row (1 for version, 1 for schema field, and 1 for the remaining term, and 1 for the term count, and 1 for the back index entry)
	expectedLength = uint64(1 + len(schema) + 1 + 1 + 1)
	rowCount = idx.rowCount()
	if rowCount != expectedLength {
		t.Errorf("expected %d rows, got: %d", expectedLength, rowCount)
	}
}

func TestIndexInsertMultiple(t *testing.T) {
	defer os.RemoveAll("test")

	schema := []*index.Field{
		&index.Field{
			Name:     "name",
			Path:     "/name",
			Analyzer: "standard",
		},
	}
	idx := NewUpsideDownCouch("test", schema)

	err := idx.Open()
	if err != nil {
		t.Errorf("error opening index: %v", err)
	}
	defer idx.Close()

	doc := []byte(`{"name": "test"}`)
	err = idx.Update([]byte{'1'}, doc)
	if err != nil {
		t.Errorf("Error updating index: %v", err)
	}

	err = idx.Update([]byte{'2'}, doc)
	if err != nil {
		t.Errorf("Error updating index: %v", err)
	}

	// should have 4 rows (1 for version, 1 for schema field, and 2 for single term, and 1 for the term count,  and 2 for the back index entries)
	expectedLength := uint64(1 + len(schema) + 2 + 1 + 2)
	rowCount := idx.rowCount()
	if rowCount != expectedLength {
		t.Errorf("expected %d rows, got: %d", expectedLength, rowCount)
	}
}
