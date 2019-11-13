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

package rsmdb

import (
	"encoding/json"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"reflect"
	"rsm/enums"
	"rsm/models"
)

type rsmReaderInstance struct {
	sdl common.ISdlInstance
}

/*
RNibReader interface allows retrieving data from redis BD by various keys
*/
type RsmReader interface {
	GetRsmGeneralConfiguration() (*models.RsmGeneralConfiguration, error)
	GetRsmRanInfo(ranName string) (*models.RsmRanInfo, error)
}

/*
GetRNibReader returns reference to RNibReader
*/
func GetRsmReader(sdl common.ISdlInstance) RsmReader {
	return &rsmReaderInstance{sdl: sdl}
}

func (r *rsmReaderInstance) GetRsmRanInfo(ranName string) (*models.RsmRanInfo, error) {

	key, err := common.ValidateAndBuildNodeBNameKey(ranName)
	if err != nil {
		return nil, err
	}

	rsmRanInfo := &models.RsmRanInfo{}

	err = r.getByKeyAndUnmarshal(key, rsmRanInfo)

	if err != nil {
		return nil, err
	}

	return rsmRanInfo, nil
}

// TODO: implement
func (r *rsmReaderInstance) GetRsmGeneralConfiguration() (*models.RsmGeneralConfiguration, error) {
	return &models.RsmGeneralConfiguration{
		EnableResourceStatus:         true,
		PartialSuccessAllowed:        true,
		PrbPeriodic:                  true,
		TnlLoadIndPeriodic:           true,
		HwLoadIndPeriodic:            true,
		AbsStatusPeriodic:            true,
		RsrpMeasurementPeriodic:      true,
		CsiPeriodic:                  true,
		PeriodicityMs:                enums.ReportingPeriodicity_one_thousand_ms,
		PeriodicityRsrpMeasurementMs: enums.ReportingPeriodicityRSRPMR_four_hundred_80_ms,
		PeriodicityCsiMs:             enums.ReportingPeriodicityCSIR_ms20,
	}, nil
}

func (r *rsmReaderInstance) getByKeyAndUnmarshal(key string, entity interface{}) error {
	data, err := r.sdl.Get([]string{key})
	if err != nil {
		return common.NewInternalError(err)
	}
	if data != nil && data[key] != nil {
		err = json.Unmarshal([]byte(data[key].(string)), entity)
		if err != nil {
			return common.NewInternalError(err)
		}
		return nil
	}
	return common.NewResourceNotFoundErrorf("#rsmReader.getByKeyAndUnmarshal - entity of type %s not found. Key: %s", reflect.TypeOf(entity).String(), key)
}
