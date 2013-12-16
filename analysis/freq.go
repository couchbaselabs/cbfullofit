//  Copyright (c) 2013 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.
package analysis

type TokenFreq struct {
	Term []byte
	Freq uint64
}

func TokenFrequency(tokens TokenStream) []*TokenFreq {
	index := make(map[string]*TokenFreq)

	for _, token := range tokens {
		curr, ok := index[string(token.Term)]
		if ok {
			curr.Freq = curr.Freq + 1
		} else {
			index[string(token.Term)] = &TokenFreq{token.Term, 1}
		}
	}

	rv := make([]*TokenFreq, len(index))
	i := 0
	for _, tf := range index {
		rv[i] = tf
		i += 1
	}

	return rv
}
