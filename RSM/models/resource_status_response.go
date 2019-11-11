package models

import (
	"fmt"
	"strings"
)

type MeasurementFailureCause struct {
	MeasurementFailedReportCharacteristics []byte
}

func (r MeasurementFailureCause) String() string {
	return fmt.Sprintf("MeasurementFailedReportCharacteristics: %x", r.MeasurementFailedReportCharacteristics)
}

type MeasurementInitiationResult struct {
	CellId                   string
	MeasurementFailureCauses []*MeasurementFailureCause
}

func (r MeasurementInitiationResult)String() string {
	var strBuilder strings.Builder
	strBuilder.WriteString("[ ")
	for _, cause := range r.MeasurementFailureCauses {
		strBuilder.WriteString(cause.String())
		strBuilder.WriteString(" ")
	}
	strBuilder.WriteString(" ]")
	return fmt.Sprintf("CellId: %s, MeasurementFailureCauses: %s", r.CellId, strBuilder.String())
}

type ResourceStatusResponse struct {
	ENB1_Measurement_ID          int64
	ENB2_Measurement_ID          int64
	MeasurementInitiationResults []*MeasurementInitiationResult
}

func (r ResourceStatusResponse)String() string{
	var strBuilder strings.Builder
	strBuilder.WriteString("[ ")
	for _, result := range r.MeasurementInitiationResults {
		strBuilder.WriteString(result.String())
		strBuilder.WriteString(" ")
	}
	strBuilder.WriteString(" ]")
	return fmt.Sprintf("ENB1_Measurement_ID: %d, ENB2_Measurement_ID: %d, MeasurementInitiationResults:%s", r.ENB1_Measurement_ID, r.ENB2_Measurement_ID, strBuilder.String())
}