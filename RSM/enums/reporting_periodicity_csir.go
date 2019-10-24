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
