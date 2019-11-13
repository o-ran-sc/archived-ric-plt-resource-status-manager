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
//


package enums

import "strconv"

var ReportingPeriodicityCsirValues = map[int]ReportingPeriodicityCSIR{
	5:  ReportingPeriodicityCSIR_ms5,
	10: ReportingPeriodicityCSIR_ms10,
	20: ReportingPeriodicityCSIR_ms20,
	40: ReportingPeriodicityCSIR_ms40,
	80: ReportingPeriodicityCSIR_ms80,
}

var ReportingPeriodicityCsirNames = map[int]string{
	1: "5",
	2: "10",
	3: "20",
	4: "40",
	5: "80",
}

type ReportingPeriodicityCSIR int

const (
	ReportingPeriodicityCSIR_ms5 ReportingPeriodicityCSIR = iota + 1
	ReportingPeriodicityCSIR_ms10
	ReportingPeriodicityCSIR_ms20
	ReportingPeriodicityCSIR_ms40
	ReportingPeriodicityCSIR_ms80
)

func (x ReportingPeriodicityCSIR) String() string {
	s, ok := ReportingPeriodicityCsirNames[int(x)]

	if ok {
		return s
	}

	return strconv.Itoa(int(x))
}

func GetReportingPeriodicityCsirValuesAsKeys() []int {
	keys := make([]int, len(ReportingPeriodicityCsirValues))

	i := 0
	for k := range ReportingPeriodicityCsirValues {
		keys[i] = k
		i++
	}

	return keys
}
