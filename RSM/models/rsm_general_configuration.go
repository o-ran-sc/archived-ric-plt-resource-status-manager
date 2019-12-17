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

import "rsm/enums"

type RsmGeneralConfiguration struct {
	EnableResourceStatus         bool                             `json:"enableResourceStatus"`
	PartialSuccessAllowed        bool                             `json:"partialSuccessAllowed"`
	PrbPeriodic                  bool                             `json:"prbPeriodic"`
	TnlLoadIndPeriodic           bool                             `json:"tnlLoadIndPeriodic"`
	HwLoadIndPeriodic            bool                             `json:"wwLoadIndPeriodic"`
	AbsStatusPeriodic            bool                             `json:"absStatusPeriodic"`
	RsrpMeasurementPeriodic      bool                             `json:"rsrpMeasurementPeriodic"`
	CsiPeriodic                  bool                             `json:"csiPeriodic"`
	PeriodicityMs                enums.ReportingPeriodicity       `json:"periodicityMs"`
	PeriodicityRsrpMeasurementMs enums.ReportingPeriodicityRSRPMR `json:"periodicityRsrpMeasurementMs"`
	PeriodicityCsiMs             enums.ReportingPeriodicityCSIR   `json:"periodicityCsiMs"`
}
