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
	"rsm/models"
)

type rsmReaderInstance struct {
	sdl common.ISdlInstance
}

// RsmReader interface allows retrieving data from redis BD by various keys
type RsmReader interface {
	GetRsmGeneralConfiguration() (*models.RsmGeneralConfiguration, error)
	GetRsmRanInfo(ranName string) (*models.RsmRanInfo, error)
}

// GetRsmReader returns reference to RsmReader
func GetRsmReader(sdl common.ISdlInstance) RsmReader {
	return &rsmReaderInstance{sdl: sdl}
}

// GetRsmRanInfo returns the rsm data associated with ran 'ranName'
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

// GetRsmGeneralConfiguration returns resource status request related configuration
func (r *rsmReaderInstance) GetRsmGeneralConfiguration() (*models.RsmGeneralConfiguration, error) {
	cfg := &models.RsmGeneralConfiguration{}
	err := r.getByKeyAndUnmarshal(buildRsmGeneralConfigurationKey(), cfg)
	return cfg, err
}

// getByKeyAndUnmarshal returns the value that is associated with key 'key' as a Go structure
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

func buildRsmGeneralConfigurationKey() string {
	return "CFG:GENERAL:v1.0.0"
}
