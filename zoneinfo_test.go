// Copyright 2019 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build !windows

package procfs

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestZoneinfo(t *testing.T) {
	fs := getProcFixtures(t)
	refs := Zoneinfo{
		Nodes: []*NodeStats{
			{Node: "0", NrInactiveAnon: newPInt64(230981), NrActiveAnon: newPInt64(547580), NrInactiveFile: newPInt64(316904),
				NrActiveFile: newPInt64(346282), NrUnevictable: newPInt64(115467), NrSlabReclaimable: newPInt64(131220),
				NrSlabUnreclaimable: newPInt64(47320), NrIsolatedAnon: newPInt64(0), NrIsolatedFile: newPInt64(0),
				/*TODO: what about workingset_nodes, workingset_restore? */
				NrWorkingsetRefault: newPInt64(466886), NrWorkingsetActivate: newPInt64(276925), NrWorkingsetNodereclaim: newPInt64(487),
				NrAnonPages: newPInt64(795576), NrMapped: newPInt64(215483), NrFilePages: newPInt64(761874), NrDirty: newPInt64(908),
				NrWriteback: newPInt64(0), NrWritebackTemp: newPInt64(0), NrShmem: newPInt64(224925), NrShmemHugepages: newPInt64(0),
				NrShmemPmdMapped: newPInt64(0), NrAnonTransparentHugepages: newPInt64(0), NrUnstable: newPInt64(0),
				NrVmscanWrite: newPInt64(12950), NrVmscanImmediateReclaim: newPInt64(3033), NrDirtied: newPInt64(8007423),
				NrWritten: newPInt64(7752121), /* nr_kernel_misc_reclaimable is not present in current kernels */
			},
		},
		Zones: []*ZoneStats{
			{Node: "0", Zone: "DMA", Free: newPInt64(3952), Min: newPInt64(33), Low: newPInt64(41), High: newPInt64(49),
				Spanned: newPInt64(4095), Present: newPInt64(3975), Managed: newPInt64(3956), NrInactiveAnon: newPInt64(0),
				NrActiveAnon: newPInt64(0), NrInactiveFile: newPInt64(0), NrActiveFile: newPInt64(0), NrUnevictable: newPInt64(0),
				NrWritePending: newPInt64(0), NrMlock: newPInt64(0), NrPageTablePages: newPInt64(0), NrKernelStack: newPInt64(0),
				NrBounce: newPInt64(0), NrZsPages: newPInt64(0), NrFreeCma: newPInt64(0), NumaHit: newPInt64(1), NumaMiss: newPInt64(0),
				NumaForeign: newPInt64(0), NumaInterleave: newPInt64(0), NumaLocal: newPInt64(1), NumaOther: newPInt64(0),
				Protection: []*int64{newPInt64(0), newPInt64(2877), newPInt64(7826), newPInt64(7826), newPInt64(7826)}},
			{Node: "0", Zone: "DMA32", Free: newPInt64(204252), Min: newPInt64(19510), Low: newPInt64(21059), High: newPInt64(22608),
				Spanned: newPInt64(1044480), Present: newPInt64(759231), Managed: newPInt64(742806), NrInactiveAnon: newPInt64(118558),
				NrActiveAnon: newPInt64(106598), NrInactiveFile: newPInt64(75475), NrActiveFile: newPInt64(70293), NrUnevictable: newPInt64(66195),
				NrWritePending: newPInt64(64), NrMlock: newPInt64(4), NrPageTablePages: newPInt64(1756), NrKernelStack: newPInt64(2208),
				NrBounce: newPInt64(0), NrZsPages: newPInt64(0), NrFreeCma: newPInt64(0), NumaHit: newPInt64(113952967), NumaMiss: newPInt64(0),
				NumaForeign: newPInt64(0), NumaInterleave: newPInt64(0), NumaLocal: newPInt64(113952967), NumaOther: newPInt64(0),
				Protection: []*int64{newPInt64(0), newPInt64(0), newPInt64(4949), newPInt64(4949), newPInt64(4949)}},
			{Node: "0", Zone: "Normal", Free: newPInt64(18553), Min: newPInt64(11176), Low: newPInt64(13842), High: newPInt64(16508),
				Spanned: newPInt64(1308160), Present: newPInt64(1308160), Managed: newPInt64(1268711), NrInactiveAnon: newPInt64(112423),
				NrActiveAnon: newPInt64(440982), NrInactiveFile: newPInt64(241429), NrActiveFile: newPInt64(275989), NrUnevictable: newPInt64(49272),
				NrWritePending: newPInt64(844), NrMlock: newPInt64(154), NrPageTablePages: newPInt64(9750), NrKernelStack: newPInt64(15136),
				NrBounce: newPInt64(0), NrZsPages: newPInt64(0), NrFreeCma: newPInt64(0), NumaHit: newPInt64(162718019), NumaMiss: newPInt64(0),
				NumaForeign: newPInt64(0), NumaInterleave: newPInt64(26812), NumaLocal: newPInt64(162718019), NumaOther: newPInt64(0),
				Protection: []*int64{newPInt64(0), newPInt64(0), newPInt64(0), newPInt64(0), newPInt64(0)}},
			{Node: "0", Zone: "Movable", Free: newPInt64(0), Min: newPInt64(0), Low: newPInt64(0), High: newPInt64(0),
				Spanned: newPInt64(0), Present: newPInt64(0), Managed: newPInt64(0),
				Protection: []*int64{newPInt64(0), newPInt64(0), newPInt64(0), newPInt64(0), newPInt64(0)}},
			{Node: "0", Zone: "Device", Free: newPInt64(0), Min: newPInt64(0), Low: newPInt64(0), High: newPInt64(0),
				Spanned: newPInt64(0), Present: newPInt64(0), Managed: newPInt64(0),
				Protection: []*int64{newPInt64(0), newPInt64(0), newPInt64(0), newPInt64(0), newPInt64(0)}},
		},
	}
	data, err := fs.Zoneinfo()
	if err != nil {
		t.Fatalf("failed to parse zoneinfo: %v", err)
	}

	want, got := refs, data
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("unexpected zoneinfo entry (-want +got):\n%s", diff)
	}
}
