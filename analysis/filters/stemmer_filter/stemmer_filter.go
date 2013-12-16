//  Copyright (c) 2013 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.
package stemmer_filter

import (
	"bitbucket.org/tebeka/snowball"
	"github.com/couchbaselabs/cbfullofit/analysis"
)

type StemmerFilter struct {
	lang    string
	stemmer *snowball.Stemmer
}

func NewStemmerFilter(lang string) (*StemmerFilter, error) {
	stemmer, err := snowball.New(lang)
	if err != nil {
		return nil, err
	}
	return &StemmerFilter{
		lang:    lang,
		stemmer: stemmer,
	}, nil
}

func (s *StemmerFilter) Filter(input analysis.TokenStream) analysis.TokenStream {
	rv := make(analysis.TokenStream, 0)

	for _, token := range input {
		stemmed := s.stemmer.Stem(string(token.Term))
		token.Term = []byte(stemmed)
		rv = append(rv, token)
	}

	return rv
}
