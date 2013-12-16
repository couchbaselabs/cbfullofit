//  Copyright (c) 2013 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.
package search

import (
	"fmt"
)

type TermDisjunctionQueryScorer struct {
	explain bool
}

func NewTermDisjunctionQueryScorer(explain bool) *TermDisjunctionQueryScorer {
	return &TermDisjunctionQueryScorer{
		explain: explain,
	}
}

func (s *TermDisjunctionQueryScorer) Score(constituents []*DocumentMatch, countMatch, countTotal int) *DocumentMatch {
	rv := DocumentMatch{
		ID: constituents[0].ID,
	}

	var sum float64
	var childrenExplanations []*Explanation
	if s.explain {
		childrenExplanations = make([]*Explanation, len(constituents))
	}

	for i, docMatch := range constituents {
		sum += docMatch.Score
		if s.explain {
			childrenExplanations[i] = docMatch.Expl
		}
	}

	var rawExpl *Explanation
	if s.explain {
		rawExpl = &Explanation{Value: sum, Message: "sum of:", Children: childrenExplanations}
	}

	coord := float64(countMatch) / float64(countTotal)
	rv.Score = sum * coord
	if s.explain {
		ce := make([]*Explanation, 2)
		ce[0] = rawExpl
		ce[1] = &Explanation{Value: coord, Message: fmt.Sprintf("coord(%d/%d)", countMatch, countTotal)}
		rv.Expl = &Explanation{Value: rv.Score, Message: "product of:", Children: ce}
	}

	return &rv
}
