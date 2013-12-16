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
	"flag"
	"log"

	"github.com/couchbaselabs/cbfullofit/index/upside_down"

	_ "github.com/couchbaselabs/cbfullofit/analysis/analyzers/standard_analyzer"
)

func main() {

	flag.Parse()

	idx := upside_down.NewUpsideDownCouch(flag.Arg(0), nil)
	err := idx.Open()
	if err != nil {
		log.Printf("error opening index: %v", err)
		return
	}
	defer idx.Close()
	idx.Dump()
}
