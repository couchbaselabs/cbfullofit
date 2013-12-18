//  Copyright (c) 2013 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.
package standard_analyzer

import (
	"github.com/couchbaselabs/cbfullofit/analysis"
	"github.com/couchbaselabs/cbfullofit/analysis/filters/lower_case_filter"
	"github.com/couchbaselabs/cbfullofit/analysis/filters/stop_words_filter"
	"github.com/couchbaselabs/cbfullofit/analysis/sanitizer/json_string_sanitizer"
	"github.com/couchbaselabs/cbfullofit/analysis/tokenizers/unicode_word_boundary"
)

func NewStandardAnalyzer() (*analysis.Analyzer, error) {
	lower_case_filter, err := lower_case_filter.NewLowerCaseFilter()
	if err != nil {
		return nil, err
	}

	stop_words_filter, err := stop_words_filter.NewStopWordsFilter()
	if err != nil {
		return nil, err
	}

	standard := analysis.Analyzer{
		Sanitizer: json_string_sanitizer.NewJsonStringSanitizer(),
		Tokenizer: unicode_word_boundary.NewUnicodeWordBoundaryTokenizer(),
		Filters: []analysis.TokenFilter{
			lower_case_filter,
			stop_words_filter,
		},
	}

	return &standard, nil
}

func init() {
	analysis.RegisterAnalyzer("standard", NewStandardAnalyzer)
}
