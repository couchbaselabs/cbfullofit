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
	"bytes"
)

var quoteBytes = []byte{'"'}

type JsonStringSanitizer struct {
}

func NewJsonStringSanitizer() *JsonStringSanitizer {
	return &JsonStringSanitizer{}
}

func (s *JsonStringSanitizer) Sanitize(input []byte) []byte {
	firstQuote := bytes.Index(input, quoteBytes)
	if firstQuote < 0 {
		// no open quote, not a JSON string, do nothing
		return input
	}
	lastQuote := bytes.LastIndex(input, quoteBytes)
	if lastQuote == firstQuote {
		// no closing quote, not a JSON string, do nothing
	}
	return input[firstQuote+1 : lastQuote]
}
