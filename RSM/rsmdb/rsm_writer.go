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

package rsmdb

import (
	"encoding/json"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"rsm/models"
)

type rsmWriterInstance struct {
	sdl common.ISdlInstance
}

// RsmWriter interface allows inserting/updating data to/in redis BD with various keys
type RsmWriter interface {
	SaveRsmRanInfo(rsmRanInfo *models.RsmRanInfo) error
	SaveRsmGeneralConfiguration(cfg *models.RsmGeneralConfiguration) error
}

// GetRsmWriter returns reference to RsmWriter
func GetRsmWriter(sdl common.ISdlInstance) RsmWriter {
	return &rsmWriterInstance{sdl: sdl}
}

// SaveRsmRanInfo saves the ran related rsm data with key RanName
func (r *rsmWriterInstance) SaveRsmRanInfo(rsmRanInfo *models.RsmRanInfo) error {

	nodebNameKey, err := common.ValidateAndBuildNodeBNameKey(rsmRanInfo.RanName)

	if err != nil {
		return err
	}

	return r.SaveWithKeyAndMarshal(nodebNameKey, rsmRanInfo)
}

// SaveRsmGeneralConfiguration saves the resource status request related configuration
func (r *rsmWriterInstance) SaveRsmGeneralConfiguration(cfg *models.RsmGeneralConfiguration) error {

	return r.SaveWithKeyAndMarshal(buildRsmGeneralConfigurationKey(), cfg)
}


// SaveWithKeyAndMarshal marshals the Go structure to json and saves it to the DB with key 'key'
func (r *rsmWriterInstance) SaveWithKeyAndMarshal(key string, entity interface{}) error {

	data, err := json.Marshal(entity)

	if err != nil {
		return common.NewInternalError(err)
	}

	var pairs []interface{}
	pairs = append(pairs, key, data)

	err = r.sdl.Set(pairs)

	if err != nil {
		return common.NewInternalError(err)
	}

	return nil
}
