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
	"log"

	"github.com/couchbaselabs/cbfullofit/index"
	"github.com/couchbaselabs/cbfullofit/index/upside_down"
	"github.com/dustin/gomemcached/client"
)

type Indexer struct {
	name   string
	bucket string
	index  index.Index
	stop   StopChannel
}

func NewIndexer(indexName string, bucket string, schema map[string]Field) *Indexer {
	usdschema := make([]*index.Field, 0)
	for fn, f := range schema {
		usdschema = append(usdschema,
			&index.Field{
				Name:     fn,
				Path:     f.Path,
				Analyzer: f.Analyzer,
			},
		)
	}
	return &Indexer{
		name:   indexName,
		bucket: bucket,
		stop:   make(StopChannel),
		index:  upside_down.NewUpsideDownCouch(*dataDir+"/"+indexName, usdschema),
	}
}

func (i *Indexer) Run() {
	i.index.Open()
	defer i.index.Close()

	args := memcached.DefaultTapArguments()
	args.Backfill = 0
	args.Checkpoint = true

	bucketDb, err := dbConnect(*cbServ, *cbPool, i.bucket)
	if err != nil {
		log.Printf("unable to find bucket to index: %v", err)
		return
	}

	feed, err := bucketDb.StartTapFeed(&args)
	if err != nil {
		log.Fatalf("Error starting tap feed: %v", err)
	}

OUTER:
	for {
		select {
		case cbEvent, ok := <-feed.C:
			if ok {
				switch cbEvent.Opcode {
				case memcached.TapMutation:
					i.index.Update(cbEvent.Key, cbEvent.Value)
				case memcached.TapDeletion:
					i.index.Delete(cbEvent.Key)
				}
			}
		case <-i.stop:
			log.Printf("Indexer '%s' asked to stop", i.name)
			break OUTER
		}
	}
	log.Printf("Indexer '%s' stoped", i.name)
}

func (i *Indexer) Stop() {
	log.Printf("Asking indexer '%s' to stop", i.name)
	close(i.stop)
}
