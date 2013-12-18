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
	"encoding/json"
	"math"
	"sort"

	"github.com/couchbaselabs/cbfullofit/analysis"
	"github.com/couchbaselabs/cbfullofit/index"
	"github.com/dustin/go-jsonpointer"
)

type mockFreq struct {
	freq    uint64
	norm    float64
	vectors []*index.TermFieldVector
}

// key doc id
type mockDocFreq map[string]*mockFreq

//key field
type mockFieldDocFreq map[string]mockDocFreq

// 2 dim array
// inner level are always pairs (field name, term)
type mockBackIndexEntry [][]string

type MockIndex struct {

	//this level of the map, the key is the term
	termIndex map[string]mockFieldDocFreq

	// key is docid
	backIndex map[string]mockBackIndexEntry

	docCount uint64
	analyzer map[string]*analysis.Analyzer
	schema   []*index.Field
}

func NewMockIndexWithDocs(schema []*index.Field, docs map[string]interface{}) *MockIndex {
	rv := NewMockIndex(schema)
	for id, doc := range docs {
		bytes, err := json.Marshal(doc)
		if err == nil {
			rv.Update([]byte(id), bytes)
		}
	}
	return rv
}

func NewMockIndex(schema []*index.Field) *MockIndex {
	mi := MockIndex{
		termIndex: make(map[string]mockFieldDocFreq),
		backIndex: make(map[string]mockBackIndexEntry),
		analyzer:  make(map[string]*analysis.Analyzer),
		schema:    schema,
	}

	for _, field := range schema {
		fieldAnalyzer, err := analysis.AnalyzerInstance(field.Analyzer)
		if err != nil {
			panic("error building analyzer")
		}
		mi.analyzer[field.Analyzer] = fieldAnalyzer
	}

	return &mi
}

func (index *MockIndex) Open() error {
	return nil
}

func (index *MockIndex) Close() {}

// for this implementation we dont care about performance
// update is simply delete then add
func (index *MockIndex) Update(id []byte, doc []byte) error {
	index.Delete(id)

	backIndexEntry := make(mockBackIndexEntry, 0)
	for fieldIndex, field := range index.schema {
		fieldValue, err := jsonpointer.Find(doc, field.Path)
		if err != nil {
			return err
		}

		analyzer := index.analyzer[field.Analyzer]
		tokens := analyzer.Analyze(fieldValue)
		fieldLength := len(tokens) // number of tokens in this doc field
		fieldNorm := 1.0 / math.Sqrt(float64(fieldLength))
		tokenFreqs := analysis.TokenFrequency(tokens)
		for _, tf := range tokenFreqs {
			mf := mockFreq{
				freq: uint64(len(tf.Locations)),
				norm: fieldNorm,
			}
			if field.IncludeTermVectors {
				mf.vectors = index.mockVectorsFromTokenFreq(uint16(fieldIndex), tf)
			}
			termString := string(tf.Term)
			fieldMap, ok := index.termIndex[termString]
			if !ok {
				fieldMap = make(map[string]mockDocFreq)
				index.termIndex[termString] = fieldMap
			}
			docMap, ok := fieldMap[field.Name]
			if !ok {
				docMap = make(map[string]*mockFreq)
				fieldMap[field.Name] = docMap
			}
			docMap[string(id)] = &mf
			backIndexInnerEntry := []string{field.Name, termString}
			backIndexEntry = append(backIndexEntry, backIndexInnerEntry)
		}
	}
	index.backIndex[string(id)] = backIndexEntry
	index.docCount += 1
	return nil
}

func (index *MockIndex) Delete(id []byte) error {
	backIndexEntry, existed := index.backIndex[string(id)]
	if existed {
		for _, backIndexPair := range backIndexEntry {
			if len(backIndexPair) == 2 {
				field := backIndexPair[0]
				term := backIndexPair[1]
				delete(index.termIndex[term][field], string(id))
				if len(index.termIndex[term][field]) == 0 {
					delete(index.termIndex[term], field)
					if len(index.termIndex[term]) == 0 {
						delete(index.termIndex, term)
					}
				}
			}
		}
		delete(index.backIndex, string(id))
		index.docCount -= 1
	}

	return nil
}

func (index *MockIndex) TermFieldReader(term []byte, field string) (index.TermFieldReader, error) {

	fdf, ok := index.termIndex[string(term)]
	if !ok {
		fdf = make(mockFieldDocFreq)
	}
	docFreqs, ok := fdf[field]
	if !ok {
		docFreqs = make(mockDocFreq)
	}
	mtfr := mockTermFieldReader{
		index:        docFreqs,
		sortedDocIds: make(sort.StringSlice, len(docFreqs)),
		curr:         -1,
	}
	i := 0
	for k, _ := range docFreqs {
		mtfr.sortedDocIds[i] = k
		i += 1
	}
	sort.Sort(mtfr.sortedDocIds)

	return &mtfr, nil
}

func (index *MockIndex) DocCount() uint64 {
	return index.docCount
}

type mockTermFieldReader struct {
	index        mockDocFreq
	sortedDocIds sort.StringSlice
	curr         int
}

func (reader *mockTermFieldReader) Next() (*index.TermFieldDoc, error) {
	next := reader.curr + 1
	if next < len(reader.sortedDocIds) {
		nextTermKey := reader.sortedDocIds[next]
		nextTerm := reader.index[nextTermKey]
		reader.curr = next
		return &index.TermFieldDoc{ID: nextTermKey, Freq: nextTerm.freq, Norm: nextTerm.norm, Vectors: nextTerm.vectors}, nil
	}
	return nil, nil
}

func (reader *mockTermFieldReader) Advance(ID []byte) (*index.TermFieldDoc, error) {
	if reader.curr >= len(reader.sortedDocIds) {
		return nil, nil
	}
	for currTermID := reader.sortedDocIds[reader.curr]; currTermID < string(ID) && reader.curr < len(reader.sortedDocIds); reader.curr += 1 {

	}

	if reader.curr < len(reader.sortedDocIds) {
		nextTermKey := reader.sortedDocIds[reader.curr]
		nextTerm := reader.index[nextTermKey]
		return &index.TermFieldDoc{ID: nextTermKey, Freq: nextTerm.freq, Norm: nextTerm.norm, Vectors: nextTerm.vectors}, nil
	}
	return nil, nil
}

func (reader *mockTermFieldReader) Count() uint64 {
	return uint64(len(reader.sortedDocIds))
}

func (reader *mockTermFieldReader) Close() {}

func (mi *MockIndex) mockVectorsFromTokenFreq(field uint16, tf *analysis.TokenFreq) []*index.TermFieldVector {
	rv := make([]*index.TermFieldVector, len(tf.Locations))

	for i, l := range tf.Locations {
		mv := index.TermFieldVector{
			Field: mi.schema[field].Name,
			Pos:   uint64(l.Position),
			Start: uint64(l.Start),
			End:   uint64(l.End),
		}
		rv[i] = &mv
	}

	return rv
}
