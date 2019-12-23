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

//  This source code is part of the near-RT RIC (RAN Intelligent Controller)
//  platform project (RICP).

package models

import (
	"rsm/enums"
)

type RsmRanInfo struct {
	RanName           string          `json:"ranName"`
	Enb1MeasurementId int64           `json:"enb1MeasurementId"`
	Enb2MeasurementId int64           `json:"enb2MeasurementId"`
	Action            enums.RsmAction `json:"action"`
	ActionStatus      bool            `json:"actionStatus"`
}

func NewRsmRanInfo(ranName string, enb1MeasurementId int64, enb2MeasurementId int64, action enums.RsmAction, actionStatus bool) *RsmRanInfo {
	return &RsmRanInfo{
		RanName:           ranName,
		Enb1MeasurementId: enb1MeasurementId,
		Enb2MeasurementId: enb2MeasurementId,
		Action:            action,
		ActionStatus:      actionStatus,
	}
}
