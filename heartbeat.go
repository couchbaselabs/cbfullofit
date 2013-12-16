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
	"net"
	"net/url"
	"strings"
	"time"
)

const nodeTTL = 60

func oneHeartbeat(startTime time.Time) {
	u, err := url.Parse(*cbServ)
	c, err := net.Dial("tcp", u.Host)
	localAddr := ""
	if err == nil {
		localAddr = strings.Split(c.LocalAddr().String(), ":")[0]
		c.Close()
	}

	aboutMe := IndexerNode{
		Addr:     localAddr,
		Type:     "node",
		Started:  startTime,
		Time:     time.Now().UTC(),
		BindAddr: *bindAddr,
		Version:  VERSION,
		Name:     nodeID,
	}

	err = db.Set("node_"+nodeID, nodeTTL, aboutMe)
	if err != nil {
		log.Printf("Failed to record a heartbeat: %v", err)
	}
}

func heartbeat() {

	startTime := time.Now().UTC()
	period := time.Second * 5
	ticker := time.NewTicker(period)

	for {
		select {
		case <-ticker.C:
			oneHeartbeat(startTime)
		}
	}
}
