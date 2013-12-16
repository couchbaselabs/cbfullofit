//  Copyright (c) 2013 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.
package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"

	"github.com/couchbaselabs/cbfullofit/analysis/filters/lower_case_filter"
	"github.com/couchbaselabs/cbfullofit/analysis/filters/stop_words_filter"
	"github.com/couchbaselabs/cbfullofit/analysis/tokenizers/unicode_word_boundary"
)

func main() {

	in := bufio.NewReader(os.Stdin)
	input, err := ioutil.ReadAll(in)
	if err != nil {
		log.Fatalf("Error reading input: %v", err)
	}

	tokenizer := unicode_word_boundary.NewUnicodeWordBoundaryTokenizer()
	lowerCaseFilter, err := lower_case_filter.NewLowerCaseFilter()
	if err != nil {
		log.Fatalf("Error creating lower-case filter: %v", err)
	}
	stopWordsFilter, err := stop_words_filter.NewStopWordsFilter()
	if err != nil {
		log.Fatalf("Error creating stop-words filter: %v", err)
	}

	for _, token := range stopWordsFilter.Filter(lowerCaseFilter.Filter(tokenizer.Tokenize(input))) {
		log.Printf("%v", token)
	}

}
