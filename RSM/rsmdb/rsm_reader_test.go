/*******************************************************************************
 *
 *   Copyright (c) 2019 AT&T Intellectual Property.
 *
 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 *   Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and
 *   limitations under the License.
 *
 *******************************************************************************/
package rsmdb

import (
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
//	"gerrit.o-ran-sc.org/r/ric-plt/sdlgo"
	"github.com/stretchr/testify/assert"
//	"os"
	"rsm/enums"
	"rsm/mocks"
	"rsm/models"
	"testing"
)

func TestGetRsmRanInfo(t *testing.T) {
	sdl := &mocks.MockSdlInstance{}
	reader := GetRsmReader(sdl)
	ranName := "test1"
	key, _ := common.ValidateAndBuildNodeBNameKey(ranName)
	infoAsGoType := models.RsmRanInfo{
		RanName:           ranName,
		Enb1MeasurementId: 1,
		Enb2MeasurementId: 2,
		Action:            enums.Start,
		ActionStatus:      false,
	}
	infoAsDbType:= "{\"ranName\":\"test1\",\"enb1MeasurementId\":1,\"enb2MeasurementId\":2,\"action\":\"start\",\"actionStatus\":false}"
	sdl.On("Get", []string{key}).Return(map[string]interface{}{key: infoAsDbType}, nil)
	info, err := reader.GetRsmRanInfo(ranName)
	if err != nil {
		t.Errorf("want: success, got: error: %v\n", err)
	}
	assert.Equal(t, info, &infoAsGoType)
}

func TestGetRsmRanInfoValidationError(t *testing.T) {
	sdl := &mocks.MockSdlInstance{}
	reader := GetRsmReader(sdl)
	ranName := ""
	key, _ := common.ValidateAndBuildNodeBNameKey(ranName)
	sdl.On("Get", []string{key}).Return(map[string]interface{}{key: ""}, nil)
	_, err := reader.GetRsmRanInfo(ranName)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "#utils.ValidateAndBuildNodeBNameKey - an empty inventory name received")
}

func TestGetRsmRanInfoDbError(t *testing.T) {
	sdl := &mocks.MockSdlInstance{}
	reader := GetRsmReader(sdl)
	ranName := "test1"
	key, _ := common.ValidateAndBuildNodeBNameKey(ranName)
	sdl.On("Get", []string{key}).Return((map[string]interface{})(nil), fmt.Errorf("db error"))
	_, err := reader.GetRsmRanInfo(ranName)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "db error")
}

func TestGetGeneralConfiguration(t *testing.T) {
	sdl := &mocks.MockSdlInstance{}
	reader := GetRsmReader(sdl)
	var testCases = []struct {
		cfgAsGoType models.RsmGeneralConfiguration
		cfgAsDbType string
	}{
		{
			cfgAsGoType: models.RsmGeneralConfiguration{
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
			},

			cfgAsDbType: "{\"enableResourceStatus\":true,\"partialSuccessAllowed\":true,\"prbPeriodic\":true,\"tnlLoadIndPeriodic\":true,\"wwLoadIndPeriodic\":true,\"absStatusPeriodic\":true,\"rsrpMeasurementPeriodic\":true,\"csiPeriodic\":true,\"periodicityMs\":1,\"periodicityRsrpMeasurementMs\":3,\"periodicityCsiMs\":3}",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.cfgAsDbType, func(t *testing.T) {
			key := buildRsmGeneralConfigurationKey()
			sdl.On("Get", []string{key}).Return(map[string]interface{}{key: tc.cfgAsDbType}, nil)
			cfg, err := reader.GetRsmGeneralConfiguration()
			if err != nil {
				t.Errorf("want: success, got: error: %v\n", err)
			}
			assert.Equal(t, cfg, &tc.cfgAsGoType)
		})
	}
}

func TestGetGeneralConfigurationNotFound(t *testing.T) {
	sdl := &mocks.MockSdlInstance{}
	reader := GetRsmReader(sdl)

	key := buildRsmGeneralConfigurationKey()
	sdl.On("Get", []string{key}).Return((map[string]interface{})(nil), nil)
	_, err := reader.GetRsmGeneralConfiguration()
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "#rsmReader.getByKeyAndUnmarshal - entity of type *models.RsmGeneralConfiguration not found. Key: CFG:GENERAL:v1.0.0")
}

func TestGetGeneralConfigurationDbError(t *testing.T) {
	sdl := &mocks.MockSdlInstance{}
	reader := GetRsmReader(sdl)

	key := buildRsmGeneralConfigurationKey()
	sdl.On("Get", []string{key}).Return((map[string]interface{})(nil), fmt.Errorf("db error"))
	_, err := reader.GetRsmGeneralConfiguration()
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "db error")
}

func TestGetGeneralConfigurationUnmarshalError(t *testing.T) {
	sdl := &mocks.MockSdlInstance{}
	reader := GetRsmReader(sdl)
	cfgAsDbTYpe := "{\"enableResourceStatus\":true, partialSuccessAllowed\":true,\"prbPeriodic\":true,\"tnlLoadIndPeriodic\":true,\"wwLoadIndPeriodic\":true,\"absStatusPeriodic\":true,\"rsrpMeasurementPeriodic\":true,\"csiPeriodic\":true,\"periodicityMs\":1,\"periodicityRsrpMeasurementMs\":3,\"periodicityCsiMs\":3}"
	key := buildRsmGeneralConfigurationKey()
	sdl.On("Get", []string{key}).Return(map[string]interface{}{key: cfgAsDbTYpe}, nil)
	_, err := reader.GetRsmGeneralConfiguration()
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "invalid character 'p' looking for beginning of object key string")
}

/*
Test against redis.
Test execution depends on the existence of the environment variable DBAAS_SERVICE_HOST.
*

func TestGetGeneralConfigurationIntegration(t *testing.T) {
	if len(os.Getenv("DBAAS_SERVICE_HOST")) == 0 {
		return
	}
	db := sdlgo.NewDatabase()
	sdl := sdlgo.NewSdlInstance("rsm", db)
	reader := GetRsmReader(sdl)
	cfgAsGoType := models.RsmGeneralConfiguration{
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
	}
	cfg, err := reader.GetRsmGeneralConfiguration()
	if err != nil {
		t.Errorf("want: success, got: error: %v\n", err)
	}

	assert.Equal(t, &cfgAsGoType, cfg)
}
*/