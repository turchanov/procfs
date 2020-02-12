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
	"bytes"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/prometheus/procfs/internal/util"
)

// Tag values describe a source field from /proc/zoneinfo

type ZoneStats struct {
	Node             string
	Zone             string
	Free             *int64   `json:"free"`
	Min              *int64   `json:"min"`
	Low              *int64   `json:"low"`
	High             *int64   `json:"high"`
	Scanned          *int64   `json:"scanned"`
	Spanned          *int64   `json:"spanned"`
	Present          *int64   `json:"present"`
	Managed          *int64   `json:"managed"`
	NrInactiveAnon   *int64   `json:"nr_zone_inactive_anon"`
	NrActiveAnon     *int64   `json:"nr_zone_active_anon"`
	NrInactiveFile   *int64   `json:"nr_zone_inactive_file"`
	NrActiveFile     *int64   `json:"nr_zone_active_file"`
	NrUnevictable    *int64   `json:"nr_zone_unevictable"`
	NrWritePending   *int64   `json:"nr_zone_write_pending"`
	NrMlock          *int64   `json:"nr_mlock"`
	NrPageTablePages *int64   `json:"nr_page_table_pages"`
	NrKernelStack    *int64   `json:"nr_kernel_stack"`
	NrBounce         *int64   `json:"nr_bounce"`
	NrZsPages        *int64   `json:"nr_zspages"`
	NrFreeCma        *int64   `json:"nr_free_cma"`
	NumaHit          *int64   `json:"numa_hit"`
	NumaMiss         *int64   `json:"numa_miss"`
	NumaForeign      *int64   `json:"numa_foreign"`
	NumaInterleave   *int64   `json:"numa_interleave"`
	NumaLocal        *int64   `json:"numa_local"`
	NumaOther        *int64   `json:"numa_other"`
	Protection       []*int64 `json:"protection"`
}

type NodeStats struct {
	Node                       string
	NrInactiveAnon             *int64 `json:"nr_inactive_anon"`
	NrActiveAnon               *int64 `json:"nr_active_anon"`
	NrInactiveFile             *int64 `json:"nr_inactive_file"`
	NrActiveFile               *int64 `json:"nr_active_file"`
	NrUnevictable              *int64 `json:"nr_unevictable"`
	NrSlabReclaimable          *int64 `json:"nr_slab_reclaimable"`
	NrSlabUnreclaimable        *int64 `json:"nr_slab_unreclaimable"`
	NrIsolatedAnon             *int64 `json:"nr_isolated_anon"`
	NrIsolatedFile             *int64 `json:"nr_isolated_file"`
	NrWorkingsetRefault        *int64 `json:"workingset_refault"`
	NrWorkingsetActivate       *int64 `json:"workingset_activate"`
	NrWorkingsetNodereclaim    *int64 `json:"workingset_nodereclaim"`
	NrAnonPages                *int64 `json:"nr_anon_pages"`
	NrMapped                   *int64 `json:"nr_mapped"`
	NrFilePages                *int64 `json:"nr_file_pages"`
	NrDirty                    *int64 `json:"nr_dirty"`
	NrWriteback                *int64 `json:"nr_writeback"`
	NrWritebackTemp            *int64 `json:"nr_writeback_temp"`
	NrShmem                    *int64 `json:"nr_shmem"`
	NrShmemHugepages           *int64 `json:"nr_shmem_hugepages"`
	NrShmemPmdMapped           *int64 `json:"nr_shmem_pmdmapped"`
	NrAnonTransparentHugepages *int64 `json:"nr_anon_transparent_hugepages"`
	NrUnstable                 *int64 `json:"nr_unstable"`
	NrVmscanWrite              *int64 `json:"nr_vmscan_write"`
	NrVmscanImmediateReclaim   *int64 `json:"nr_vmscan_immediate_reclaim"`
	NrDirtied                  *int64 `json:"nr_dirtied"`
	NrWritten                  *int64 `json:"nr_written"`
}

// Zoneinfo holds info parsed from /proc/zoneinfo.
type Zoneinfo struct {
	Nodes []*NodeStats
	Zones []*ZoneStats
}

const (
	ProcessingNodeStats = 1
	ProcessingZoneStats = 2
)

var nodeZoneRE = regexp.MustCompile(`(\d+), zone\s+(\w+)`)

// Zoneinfo parses an zoneinfo-file (/proc/zoneinfo) and returns a slice of
// structs containing the relevant info.  More information available here:
// https://www.kernel.org/doc/Documentation/sysctl/vm.txt
func (fs FS) Zoneinfo() (Zoneinfo, error) {
	data, err := ioutil.ReadFile(fs.proc.Path("zoneinfo"))
	if err != nil {
		return Zoneinfo{}, fmt.Errorf("error reading zoneinfo %s: %s", fs.proc.Path("zoneinfo"), err)
	}
	zoneinfo, err := parseZoneinfo(data)
	if err != nil {
		return Zoneinfo{}, fmt.Errorf("error parsing zoneinfo %s: %s", fs.proc.Path("zoneinfo"), err)
	}
	return zoneinfo, nil
}

func (stats *NodeStats) parse(line string) {
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return
	}
	vp := util.NewValueParser(parts[1])

	switch parts[0] {
	case "nr_inactive_anon":
		stats.NrInactiveAnon = vp.PInt64()
	case "nr_active_anon":
		stats.NrActiveAnon = vp.PInt64()
	case "nr_inactive_file":
		stats.NrInactiveFile = vp.PInt64()
	case "nr_active_file":
		stats.NrActiveFile = vp.PInt64()
	case "nr_unevictable":
		stats.NrUnevictable = vp.PInt64()
	case "nr_slab_reclaimable":
		stats.NrSlabReclaimable = vp.PInt64()
	case "nr_slab_unreclaimable":
		stats.NrSlabUnreclaimable = vp.PInt64()
	case "nr_isolated_anon":
		stats.NrIsolatedAnon = vp.PInt64()
	case "nr_isolated_file":
		stats.NrIsolatedFile = vp.PInt64()
	case "workingset_refault":
		stats.NrWorkingsetRefault = vp.PInt64()
	case "workingset_activate":
		stats.NrWorkingsetActivate = vp.PInt64()
	case "workingset_nodereclaim":
		stats.NrWorkingsetNodereclaim = vp.PInt64()
	case "nr_anon_pages":
		stats.NrAnonPages = vp.PInt64()
	case "nr_mapped":
		stats.NrMapped = vp.PInt64()
	case "nr_file_pages":
		stats.NrFilePages = vp.PInt64()
	case "nr_dirty":
		stats.NrDirty = vp.PInt64()
	case "nr_writeback":
		stats.NrWriteback = vp.PInt64()
	case "nr_writeback_temp":
		stats.NrWritebackTemp = vp.PInt64()
	case "nr_shmem":
		stats.NrShmem = vp.PInt64()
	case "nr_shmem_hugepages":
		stats.NrShmemHugepages = vp.PInt64()
	case "nr_shmem_pmdmapped":
		stats.NrShmemPmdMapped = vp.PInt64()
	case "nr_anon_transparent_hugepages":
		stats.NrAnonTransparentHugepages = vp.PInt64()
	case "nr_unstable":
		stats.NrUnstable = vp.PInt64()
	case "nr_vmscan_write":
		stats.NrVmscanWrite = vp.PInt64()
	case "nr_vmscan_immediate_reclaim":
		stats.NrVmscanImmediateReclaim = vp.PInt64()
	case "nr_dirtied":
		stats.NrDirtied = vp.PInt64()
	case "nr_written":
		stats.NrWritten = vp.PInt64()
	}
}

func (stats *ZoneStats) parse(line string) {
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return
	}

	if (parts[0] == "pages") && (parts[1] == "free") {
		vp := util.NewValueParser(parts[2])
		stats.Free = vp.PInt64()
		return
	}

	vp := util.NewValueParser(parts[1])

	switch parts[0] {
	case "min":
		stats.Min = vp.PInt64()
	case "low":
		stats.Low = vp.PInt64()
	case "high":
		stats.High = vp.PInt64()
	case "scanned":
		stats.Scanned = vp.PInt64()
	case "spanned":
		stats.Spanned = vp.PInt64()
	case "present":
		stats.Present = vp.PInt64()
	case "managed":
		stats.Managed = vp.PInt64()
	case "nr_zone_inactive_anon":
		stats.NrInactiveAnon = vp.PInt64()
	case "nr_zone_active_anon":
		stats.NrActiveAnon = vp.PInt64()
	case "nr_zone_inactive_file":
		stats.NrInactiveFile = vp.PInt64()
	case "nr_zone_active_file":
		stats.NrActiveFile = vp.PInt64()
	case "nr_zone_unevictable":
		stats.NrUnevictable = vp.PInt64()
	case "nr_zone_write_pending":
		stats.NrWritePending = vp.PInt64()
	case "nr_mlock":
		stats.NrMlock = vp.PInt64()
	case "nr_page_table_pages":
		stats.NrPageTablePages = vp.PInt64()
	case "nr_kernel_stack":
		stats.NrKernelStack = vp.PInt64()
	case "nr_bounce":
		stats.NrBounce = vp.PInt64()
	case "nr_zspages":
		stats.NrZsPages = vp.PInt64()
	case "nr_free_cma":
		stats.NrFreeCma = vp.PInt64()
	case "numa_hit":
		stats.NumaHit = vp.PInt64()
	case "numa_miss":
		stats.NumaMiss = vp.PInt64()
	case "numa_foreign":
		stats.NumaForeign = vp.PInt64()
	case "numa_interleave":
		stats.NumaInterleave = vp.PInt64()
	case "numa_local":
		stats.NumaLocal = vp.PInt64()
	case "numa_other":
		stats.NumaOther = vp.PInt64()
	case "protection:":
		protectionParts := strings.Split(line, ":")
		protectionValues := strings.Replace(protectionParts[1], "(", "", 1)
		protectionValues = strings.Replace(protectionValues, ")", "", 1)
		protectionValues = strings.TrimSpace(protectionValues)
		protectionStringMap := strings.Split(protectionValues, ", ")
		val, err := util.ParsePInt64s(protectionStringMap)
		if err == nil {
			stats.Protection = val
		}
	}
}

func parseZoneinfo(zoneinfoData []byte) (Zoneinfo, error) {
	zoneinfo := Zoneinfo{}

	zoneinfoBlocks := bytes.Split(zoneinfoData, []byte("\nNode"))
	for _, block := range zoneinfoBlocks {
		var currentNode, currentZone string
		var nodestats *NodeStats
		var zonestats *ZoneStats

		data := strings.Split(string(block), "\n")

		// This must not happen but still we have to check the size of "data" before slicing it
		if len(data) < 2 {
			continue
		}
		header, lines := data[0], data[1:]

		// First line must be "(Node )?\d+, zone\s+\w+" since we split zoneinfoData by "\nNode"
		nodeZone := nodeZoneRE.FindStringSubmatch(header)
		if nodeZone == nil {
			continue
		}
		currentNode = nodeZone[1]
		currentZone = nodeZone[2]
		zonestats = &ZoneStats{Node: currentNode, Zone: currentZone}

		state := ProcessingZoneStats
		for _, line := range lines {
			line = strings.TrimSpace(line)

			if strings.HasPrefix(line, "per-node stats") {
				state = ProcessingNodeStats
				nodestats = &NodeStats{Node: currentNode}
				continue
			} else if strings.HasPrefix(line, "pages free") {
				state = ProcessingZoneStats
			}

			switch state {
			case ProcessingNodeStats:
				nodestats.parse(line)
			case ProcessingZoneStats:
				zonestats.parse(line)
			}
		}

		if nodestats != nil {
			zoneinfo.Nodes = append(zoneinfo.Nodes, nodestats)
		}
		zoneinfo.Zones = append(zoneinfo.Zones, zonestats)
	}
	return zoneinfo, nil
}
