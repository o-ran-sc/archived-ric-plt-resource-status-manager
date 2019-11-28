//
// Copyright 2019 AT&T Intellectual Property
// Copyright 2019 Nokia
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//  This source code is part of the near-RT RIC (RAN Intelligent Controller)
//  platform project (RICP).


//  This source code is part of the near-RT RIC (RAN Intelligent Controller)
//  platform project (RICP).

package enums

import (
	"strconv"
)

var ReportingPeriodicityRsrPmrValues = map[int]ReportingPeriodicityRSRPMR{
	120: ReportingPeriodicityRSRPMR_one_hundred_20_ms,
	240: ReportingPeriodicityRSRPMR_two_hundred_40_ms,
	480: ReportingPeriodicityRSRPMR_four_hundred_80_ms,
	640: ReportingPeriodicityRSRPMR_six_hundred_40_ms,
}

var ReportingPeriodicityRsrPmrNames = map[int]string{
	1: "120",
	2: "240",
	3: "480",
	4: "640",
}

type ReportingPeriodicityRSRPMR int

const (
	ReportingPeriodicityRSRPMR_one_hundred_20_ms ReportingPeriodicityRSRPMR = iota + 1
	ReportingPeriodicityRSRPMR_two_hundred_40_ms
	ReportingPeriodicityRSRPMR_four_hundred_80_ms
	ReportingPeriodicityRSRPMR_six_hundred_40_ms
)

func (x ReportingPeriodicityRSRPMR) String() string {
	s, ok := ReportingPeriodicityRsrPmrNames[int(x)]

	if ok {
		return s
	}

	return strconv.Itoa(int(x))
}

func GetReportingPeriodicityRsrPmrValuesAsKeys() []int {
	keys := make([]int, len(ReportingPeriodicityRsrPmrValues))

	i := 0
	for k := range ReportingPeriodicityRsrPmrValues {
		keys[i] = k
		i++
	}

	return keys
}
