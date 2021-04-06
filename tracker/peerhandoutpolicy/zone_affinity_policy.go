// Copyright (c) 2016-2019 Uber Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package peerhandoutpolicy

import "github.com/uber/kraken/core"

const _zone_affinity_policy = "zone_affinity"

// zoneAffinityPolicy assigns priorities based on download completeness within a Zone.
// Zone-Peers who've completed downloading are highest, then origins, then other peers.
// This allows for situations where bandwidth between zones is quite limited
type zoneAffinityPolicy struct{}

func newZoneAffinityPolicy() assignmentPolicy {
	return &zoneAffinityPolicy{}
}

func (p *zoneAffinityPolicy) assignPriority(source *core.PeerInfo, peer *core.PeerInfo) (int, string) {
	if source.Zone == peer.Zone {
		if peer.Complete {
			return 0, "peer_seeder"
		} else {
			return 1, "peer_incomplete"
		}
	}

	if peer.Origin {
		return 2, "origin"
	}

	return -1, "restricted"
}
