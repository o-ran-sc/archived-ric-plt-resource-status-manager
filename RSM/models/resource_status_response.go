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

func (r MeasurementInitiationResult) String() string {
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

func (r ResourceStatusResponse) String() string {
	var strBuilder strings.Builder
	strBuilder.WriteString("[ ")
	for _, result := range r.MeasurementInitiationResults {
		strBuilder.WriteString(result.String())
		strBuilder.WriteString(" ")
	}
	strBuilder.WriteString(" ]")
	return fmt.Sprintf("ENB1_Measurement_ID: %d, ENB2_Measurement_ID: %d, MeasurementInitiationResults:%s", r.ENB1_Measurement_ID, r.ENB2_Measurement_ID, strBuilder.String())
}
