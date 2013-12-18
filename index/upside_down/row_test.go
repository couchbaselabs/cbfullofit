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
	"reflect"
	"testing"
)

func TestRows(t *testing.T) {
	tests := []struct {
		input  UpsideDownCouchRow
		outKey []byte
		outVal []byte
	}{
		{
			NewVersionRow(1),
			[]byte{'v'},
			[]byte{0x1},
		},
		{
			NewFieldRow(0, "name", "/name", "standard", false),
			[]byte{'f', 0, 0},
			[]byte{'n', 'a', 'm', 'e', BYTE_SEPARATOR, '/', 'n', 'a', 'm', 'e', BYTE_SEPARATOR, 's', 't', 'a', 'n', 'd', 'a', 'r', 'd', BYTE_SEPARATOR, 0},
		},
		{
			NewFieldRow(1, "desc", "/description", "standard", true),
			[]byte{'f', 1, 0},
			[]byte{'d', 'e', 's', 'c', BYTE_SEPARATOR, '/', 'd', 'e', 's', 'c', 'r', 'i', 'p', 't', 'i', 'o', 'n', BYTE_SEPARATOR, 's', 't', 'a', 'n', 'd', 'a', 'r', 'd', BYTE_SEPARATOR, 1},
		},
		{
			NewFieldRow(513, "style", "/style", "keyword", false),
			[]byte{'f', 1, 2},
			[]byte{'s', 't', 'y', 'l', 'e', BYTE_SEPARATOR, '/', 's', 't', 'y', 'l', 'e', BYTE_SEPARATOR, 'k', 'e', 'y', 'w', 'o', 'r', 'd', BYTE_SEPARATOR, 0},
		},
		{
			NewTermFrequencyRow([]byte{'b', 'e', 'e', 'r'}, 0, nil, 3, 3.14),
			[]byte{'t', 'b', 'e', 'e', 'r', BYTE_SEPARATOR, 0, 0},
			[]byte{3, 0, 0, 0, 0, 0, 0, 0, 195, 245, 72, 64},
		},
		{
			NewTermFrequencyRow([]byte{'b', 'e', 'e', 'r'}, 0, []byte{'b', 'u', 'd', 'w', 'e', 'i', 's', 'e', 'r'}, 3, 3.14),
			[]byte{'t', 'b', 'e', 'e', 'r', BYTE_SEPARATOR, 0, 0, 'b', 'u', 'd', 'w', 'e', 'i', 's', 'e', 'r'},
			[]byte{3, 0, 0, 0, 0, 0, 0, 0, 195, 245, 72, 64},
		},
		{
			NewTermFrequencyRowWithTermVectors([]byte{'b', 'e', 'e', 'r'}, 0, []byte{'b', 'u', 'd', 'w', 'e', 'i', 's', 'e', 'r'}, 3, 3.14, []*TermVector{&TermVector{field: 0, pos: 1, start: 3, end: 11}, &TermVector{field: 0, pos: 2, start: 23, end: 31}, &TermVector{field: 0, pos: 3, start: 43, end: 51}}),
			[]byte{'t', 'b', 'e', 'e', 'r', BYTE_SEPARATOR, 0, 0, 'b', 'u', 'd', 'w', 'e', 'i', 's', 'e', 'r'},
			[]byte{3, 0, 0, 0, 0, 0, 0, 0, 195, 245, 72, 64, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 11, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 23, 0, 0, 0, 0, 0, 0, 0, 31, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 43, 0, 0, 0, 0, 0, 0, 0, 51, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			NewBackIndexRow([]byte{'b', 'u', 'd', 'w', 'e', 'i', 's', 'e', 'r'}, []*BackIndexEntry{&BackIndexEntry{[]byte{'b', 'e', 'e', 'r'}, 0}}),
			[]byte{'b', 'b', 'u', 'd', 'w', 'e', 'i', 's', 'e', 'r'},
			[]byte{'b', 'e', 'e', 'r', BYTE_SEPARATOR, 0, 0},
		},
		{
			NewBackIndexRow([]byte{'b', 'u', 'd', 'w', 'e', 'i', 's', 'e', 'r'}, []*BackIndexEntry{&BackIndexEntry{[]byte{'b', 'e', 'e', 'r'}, 0}, &BackIndexEntry{[]byte{'b', 'e', 'a', 't'}, 1}}),
			[]byte{'b', 'b', 'u', 'd', 'w', 'e', 'i', 's', 'e', 'r'},
			[]byte{'b', 'e', 'e', 'r', BYTE_SEPARATOR, 0, 0, 'b', 'e', 'a', 't', BYTE_SEPARATOR, 1, 0},
		},
	}

	// test going from struct to k/v bytes
	for _, test := range tests {
		rk := test.input.Key()
		if !reflect.DeepEqual(rk, test.outKey) {
			t.Errorf("Expected key to be %v got: %v", test.outKey, rk)
		}
		rv := test.input.Value()
		if !reflect.DeepEqual(rv, test.outVal) {
			t.Errorf("Expected value to be %v got: %v", test.outVal, rv)
		}
	}

	// now test going back from k/v bytes to struct
	for _, test := range tests {
		row := ParseFromKeyValue(test.outKey, test.outVal)
		if !reflect.DeepEqual(row, test.input) {
			t.Errorf("Expected: %#v got: %#v", test.input, row)
		}
	}

}
