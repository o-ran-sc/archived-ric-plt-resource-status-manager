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
	"github.com/stretchr/testify/assert"
	"rsm/enums"
	"rsm/mocks"
	"rsm/models"
	"testing"
)

func TestSaveRsmRanInfo(t *testing.T) {
	sdl := &mocks.MockSdlInstance{}
	writer := GetRsmWriter(sdl)
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
	sdl.On("Set",[]interface{}{[]interface{}{key, []byte(infoAsDbType)}}).Return(nil)
	err := writer.SaveRsmRanInfo(&infoAsGoType)
	if err != nil {
		t.Errorf("want: success, got: error: %v\n", err)
	}
	sdl.AssertNumberOfCalls(t, "Set",1)
}


func TestSaveRsmRanInfoValidationError(t *testing.T) {
	sdl := &mocks.MockSdlInstance{}
	writer := GetRsmWriter(sdl)
	ranName := ""
	key, _ := common.ValidateAndBuildNodeBNameKey(ranName)
	infoAsGoType := models.RsmRanInfo{
		RanName:           ranName,
		Enb1MeasurementId: 1,
		Enb2MeasurementId: 2,
		Action:            enums.Start,
		ActionStatus:      false,
	}
	infoAsDbType:= "{\"ranName\":\"test1\",\"enb1MeasurementId\":1,\"enb2MeasurementId\":2,\"action\":\"start\",\"actionStatus\":false}"
	sdl.On("Set",[]interface{}{[]interface{}{key, []byte(infoAsDbType)}}).Return(nil)
	err := writer.SaveRsmRanInfo(&infoAsGoType)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "#utils.ValidateAndBuildNodeBNameKey - an empty inventory name received")
}


func TestSaveGeneralConfiguration(t *testing.T) {
	sdl := &mocks.MockSdlInstance{}
	writer := GetRsmWriter(sdl)
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
			key:= buildRsmGeneralConfigurationKey()
			sdl.On("Set",[]interface{}{[]interface{}{key, []byte(tc.cfgAsDbType)}}).Return(nil)
			err := writer.SaveRsmGeneralConfiguration(&tc.cfgAsGoType)
			if err != nil {
				t.Errorf("want: success, got: error: %v\n", err)
			}

			sdl.AssertNumberOfCalls(t, "Set",1)
		})
	}
}


func TestSaveGeneralConfigurationDbError(t *testing.T) {
	sdl := &mocks.MockSdlInstance{}
	writer := GetRsmWriter(sdl)

	cfgAsGoType:= models.RsmGeneralConfiguration{
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
	cfgAsDbTYpe:= "{\"enableResourceStatus\":true,\"partialSuccessAllowed\":true,\"prbPeriodic\":true,\"tnlLoadIndPeriodic\":true,\"wwLoadIndPeriodic\":true,\"absStatusPeriodic\":true,\"rsrpMeasurementPeriodic\":true,\"csiPeriodic\":true,\"periodicityMs\":1,\"periodicityRsrpMeasurementMs\":3,\"periodicityCsiMs\":3}"
	key:= buildRsmGeneralConfigurationKey()
	sdl.On("Set",[]interface{}{[]interface{}{key, []byte(cfgAsDbTYpe)}}).Return(fmt.Errorf("db error"))
	err := writer.SaveRsmGeneralConfiguration(&cfgAsGoType)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "db error")
}
