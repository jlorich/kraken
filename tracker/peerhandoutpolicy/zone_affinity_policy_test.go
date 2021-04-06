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

import (
	"math/rand"
	"testing"

	"github.com/uber/kraken/core"

	"github.com/stretchr/testify/require"
	"github.com/uber-go/tally"
)

func TestZoneAffinityPriorityPolicy(t *testing.T) {
	require := require.New(t)

	policy, err := NewPriorityPolicy(tally.NoopScope, _zone_affinity_policy)
	require.NoError(err)

	inZoneSeeders := 10
	inZonePeers := 3
	origins := 3
	outOfZonePeers := 20

	source := core.PeerInfoFixture()
	source.Zone = "Test"

	peers := make([]*core.PeerInfo, inZoneSeeders+inZonePeers+origins+outOfZonePeers)
	for k := 0; k < len(peers); k++ {
		p := core.PeerInfoFixture()

		if k < inZoneSeeders {
			p.Zone = "Test"
			p.Complete = true
		} else if k < inZoneSeeders+inZonePeers {
			p.Zone = "Test"
		} else if k < inZoneSeeders+inZonePeers+origins {
			p.Complete = true
			p.Origin = true
			p.Zone = "Origin"
		} else {
			p.Zone = "NotTest"
		}
		peers[k] = p
	}

	// shuffle
	for i := 0; i < len(peers); i++ {
		j := rand.Intn(i + 1)
		peers[i], peers[j] = peers[j], peers[i]
	}

	peers = policy.SortPeers(source, peers)
	require.Len(peers, inZoneSeeders+inZonePeers+origins)
	for k := 0; k < len(peers); k++ {
		p := peers[k]
		if k < inZoneSeeders {
			require.EqualValues(p.Zone, "Test")
			require.True(p.Complete)
			require.False(p.Origin)
		} else if k < inZonePeers {
			require.EqualValues(p.Zone, "Test")
			require.False(p.Complete)
			require.False(p.Origin)
		} else if k < origins {
			require.EqualValues(p.Zone, "Origin")
			require.True(p.Complete)
			require.True(p.Origin)
		}
	}
}
