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
	_ "github.com/couchbaselabs/cbfullofit/analysis/analyzers/standard_analyzer"
	"github.com/couchbaselabs/cbfullofit/index"
	"github.com/couchbaselabs/cbfullofit/index/mock"
)

// sets up some mock data used in many tests in this package

var twoDocIndexSchema = []*index.Field{
	&index.Field{
		Name:     "name",
		Path:     "/name",
		Analyzer: "standard",
	},
	&index.Field{
		Name:     "desc",
		Path:     "/description",
		Analyzer: "standard",
	},
}

var twoDocIndexDocs = map[string]interface{}{
	// must have 4/4 beer
	"1": map[string]interface{}{
		"name":        "marty",
		"description": "beer beer beer beer",
	},
	// must have 1/4 beer
	"2": map[string]interface{}{
		"name":        "steve",
		"description": "angst beer couch database",
	},
	// must have 1/4 beer
	"3": map[string]interface{}{
		"name":        "dustin",
		"description": "apple beer column dank",
	},
	// must have 65/65 beer
	"4": map[string]interface{}{
		"name":        "ravi",
		"description": "beer beer beer beer beer beer beer beer beer beer beer beer beer beer beer beer beer beer beer beer beer beer beer beer beer beer beer beer beer beer beer beer beer beer beer beer beer beer beer beer beer beer beer beer beer beer beer beer beer beer beer beer beer beer beer beer beer beer beer beer beer beer beer beer beer",
	},
	// must have 0/x beer
	"5": map[string]interface{}{
		"name":        "bobert",
		"description": "water",
	},
}

var twoDocIndex *mock.MockIndex = mock.NewMockIndexWithDocs(twoDocIndexSchema, twoDocIndexDocs)
