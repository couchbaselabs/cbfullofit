//  Copyright (c) 2013 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.
package analysis

import (
	"fmt"
)

type Token struct {
	Start    int
	End      int
	Term     []byte
	Position int
}

func (t *Token) String() string {
	return fmt.Sprintf("Start: %d  End: %d  Position: %d  Token: %s", t.Start, t.End, t.Position, string(t.Term))
}

type TokenStream []*Token

type Tokenizer interface {
	Tokenize([]byte) TokenStream
}

type TokenFilter interface {
	Filter(TokenStream) TokenStream
}

type Analyzer struct {
	Tokenizer Tokenizer
	Filters   []TokenFilter
}

func (a *Analyzer) Analyze(input []byte) TokenStream {
	tokens := a.Tokenizer.Tokenize(input)
	for _, filter := range a.Filters {
		tokens = filter.Filter(tokens)
	}
	return tokens
}

type AnalyzerConstructor func() (*Analyzer, error)

var analyzerRegistry map[string]AnalyzerConstructor = make(map[string]AnalyzerConstructor)

func RegisterAnalyzer(name string, cons AnalyzerConstructor) {
	analyzerRegistry[name] = cons
}

func AnalyzerInstance(name string) (*Analyzer, error) {
	cons, ok := analyzerRegistry[name]
	if !ok {
		return nil, fmt.Errorf("No analyzer registered with the name '%s'", name)
	}

	analyzer, err := cons()
	return analyzer, err
}
