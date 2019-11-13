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
	"rsm/models"
)

type rsmWriterInstance struct {
	sdl common.ISdlInstance
}

/*
RNibReader interface allows retrieving data from redis BD by various keys
*/
type RsmWriter interface {
	SaveRsmRanInfo(rsmRanInfo *models.RsmRanInfo) error
}

/*
GetRNibReader returns reference to RNibReader
*/
func GetRsmWriter(sdl common.ISdlInstance) RsmWriter {
	return &rsmWriterInstance{sdl: sdl}
}

func (r *rsmWriterInstance) SaveRsmRanInfo(rsmRanInfo *models.RsmRanInfo) error {

	nodebNameKey, err := common.ValidateAndBuildNodeBNameKey(rsmRanInfo.RanName)

	if err != nil {
		return err
	}

	data, err := json.Marshal(rsmRanInfo)

	if err != nil {
		return common.NewInternalError(err)
	}

	var pairs []interface{}
	pairs = append(pairs, nodebNameKey, data)

	err = r.sdl.Set(pairs)

	if err != nil {
		return common.NewInternalError(err)
	}

	return nil
}
