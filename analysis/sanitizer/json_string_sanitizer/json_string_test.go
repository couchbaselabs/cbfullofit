//  Copyright (c) 2013 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.
package json_string_sanitizer

import (
	"reflect"
	"testing"
)

func TestJsonStringSanitizer(t *testing.T) {
	tests := []struct {
		input  []byte
		output []byte
	}{
		// normal JSON string as seen from JSON pointer (contains whitespace and surrounding quotes)
		{
			input:  []byte(` "a json string"`),
			output: []byte(`a json string`),
		},
		// something that isn't a JSON string
		{
			input:  []byte(` false`),
			output: []byte(` false`),
		},
	}

	for _, test := range tests {
		sanitizer := NewJsonStringSanitizer()
		output := sanitizer.Sanitize(test.input)
		if !reflect.DeepEqual(output, test.output) {
			t.Errorf("Expected: `%s` got: `%s` for `%s`", string(test.output), string(output), string(test.input))
		}
	}
}
