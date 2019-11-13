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

import (
	"strconv"
)

type ReportingPeriodicity int

var ReportingPeriodicityValues = map[int]ReportingPeriodicity{
	1000:  ReportingPeriodicity_one_thousand_ms,
	2000:  ReportingPeriodicity_two_thousand_ms,
	5000:  ReportingPeriodicity_five_thousand_ms,
	10000: ReportingPeriodicity_ten_thousand_ms,
}

var ReportingPeriodicityNames = map[int]string{
	1: "1000",
	2: "2000",
	3: "5000",
	4: "10000",
}

const (
	ReportingPeriodicity_one_thousand_ms ReportingPeriodicity = iota + 1
	ReportingPeriodicity_two_thousand_ms
	ReportingPeriodicity_five_thousand_ms
	ReportingPeriodicity_ten_thousand_ms
)

func (x ReportingPeriodicity) String() string {
	s, ok := ReportingPeriodicityNames[int(x)]

	if ok {
		return s
	}

	return strconv.Itoa(int(x))
}

func GetReportingPeriodicityValuesAsKeys() []int {
	keys := make([]int, len(ReportingPeriodicityValues))

	i := 0
	for k := range ReportingPeriodicityValues {
		keys[i] = k
		i++
	}

	return keys
}
