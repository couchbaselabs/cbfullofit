//  Copyright (c) 2013 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.
package keyword_analyzer

import (
	"github.com/couchbaselabs/cbfullofit/analysis"
	"github.com/couchbaselabs/cbfullofit/analysis/sanitizer/json_string_sanitizer"
	"github.com/couchbaselabs/cbfullofit/analysis/tokenizers/single_token"
)

func NewKeywordAnalyzer() (*analysis.Analyzer, error) {
	keyword := analysis.Analyzer{
		Sanitizer: json_string_sanitizer.NewJsonStringSanitizer(),
		Tokenizer: single_token.NewSingleTokenTokenizer(),
		Filters:   []analysis.TokenFilter{},
	}

	return &keyword, nil
}

func init() {
	analysis.RegisterAnalyzer("keyword", NewKeywordAnalyzer)
}
