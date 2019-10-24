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

package configuration

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"testing"
)

func TestParseConfigurationSuccess(t *testing.T) {
	config, err := ParseConfiguration()
	if err != nil {
		t.Errorf("failed to parse configuration: %s", err)
	}
	assert.Equal(t, 4800, config.Http.Port)
	assert.Equal(t, 4801, config.Rmr.Port)
	assert.Equal(t, 65536, config.Rmr.MaxMsgSize)
	assert.Equal(t, "info", config.Logging.LogLevel)

	assert.Equal(t, 3, config.Rnib.MaxRnibConnectionAttempts)
	assert.Equal(t, 10, config.Rnib.RnibRetryIntervalMs)

	assert.Equal(t, true, config.ResourceStatusParams.EnableResourceStatus)
	assert.Equal(t, true, config.ResourceStatusParams.PrbPeriodic)
	assert.Equal(t, true, config.ResourceStatusParams.TnlLoadIndPeriodic)
	assert.Equal(t, true, config.ResourceStatusParams.HwLoadIndPeriodic)
	assert.Equal(t, true, config.ResourceStatusParams.AbsStatusPeriodic)
	assert.Equal(t, true, config.ResourceStatusParams.RsrpMeasurementPeriodic)
	assert.Equal(t, true, config.ResourceStatusParams.CsiPeriodic)

	/*assert.Equal(t, 1, config.ResourceStatusParams.PeriodicityMs)
	assert.Equal(t, 120, config.ResourceStatusParams.PeriodicityRsrpMeasurementMs)
	assert.Equal(t, 5, config.ResourceStatusParams.PeriodicityCsiMs)*/

}

func TestParseConfigurationFileNotFoundFailure(t *testing.T) {
	configPath := "../resources/configuration.yaml"
	configPathTmp := "../resources/configuration.yaml_tmp"
	err := os.Rename(configPath, configPathTmp)
	if err != nil {
		t.Errorf("#http_server_test.TestParseConfigurationFileNotFoundFailure - failed to rename configuration file: %s\n", configPath)
	}
	defer func() {
		err = os.Rename(configPathTmp, configPath)
		if err != nil {
			t.Errorf("#http_server_test.TestParseConfigurationFileNotFoundFailure - failed to rename configuration file: %s\n", configPath)
		}
	}()

	_, cErr := ParseConfiguration()
	assert.Error(t, cErr)
}

func TestRmrConfigNotFoundFailure(t *testing.T) {
	yamlMap := map[string]interface{}{
		"logging":         map[string]interface{}{"logLevel": "info"},
		"http":            map[string]interface{}{"port": 631},
		"resourceStatusParams": map[string]interface{}{"enableResourceStatus": true, "periodicityCsiMs": 5},
		"rnib":            map[string]interface{}{"maxRnibConnectionAttempts": 3, "rnibRetryIntervalMs": 10},
	}
	cleanUp := prepareTempConfigForTest(t, yamlMap)
	defer cleanUp()

	_, cErr := ParseConfiguration()
	assert.EqualError(t, cErr, "#configuration.fillRmrConfig - failed to fill RMR configuration: The entry 'rmr' not found\n")
}

func TestLoggingConfigNotFoundFailure(t *testing.T) {
	yamlMap := map[string]interface{}{
		"rmr":             map[string]interface{}{"port": 6942, "maxMsgSize": 4096},
		"http":            map[string]interface{}{"port": 631},
		"resourceStatusParams": map[string]interface{}{"enableResourceStatus": true, "periodicityCsiMs": 5},
		"rnib":            map[string]interface{}{"maxRnibConnectionAttempts": 3, "rnibRetryIntervalMs": 10},
	}
	cleanUp := prepareTempConfigForTest(t, yamlMap)
	defer cleanUp()

	_, cErr := ParseConfiguration()
	assert.EqualError(t, cErr, "#configuration.fillLoggingConfig - failed to fill logging configuration: The entry 'logging' not found\n")
}

func TestHttpConfigNotFoundFailure(t *testing.T) {
	yamlMap := map[string]interface{}{
		"rmr":             map[string]interface{}{"port": 6942, "maxMsgSize": 4096},
		"logging":         map[string]interface{}{"logLevel": "info"},
		"resourceStatusParams": map[string]interface{}{"enableResourceStatus": true, "periodicityCsiMs": 5},
		"rnib":            map[string]interface{}{"maxRnibConnectionAttempts": 3, "rnibRetryIntervalMs": 10},
	}
	cleanUp := prepareTempConfigForTest(t, yamlMap)
	defer cleanUp()

	_, cErr := ParseConfiguration()
	assert.EqualError(t, cErr, "#configuration.fillHttpConfig - failed to fill HTTP configuration: The entry 'http' not found\n")
}

func TestRnibConfigNotFoundFailure(t *testing.T) {
	yamlMap := map[string]interface{}{
		"rmr":             map[string]interface{}{"port": 6942, "maxMsgSize": 4096},
		"logging":         map[string]interface{}{"logLevel": "info"},
		"http":            map[string]interface{}{"port": 631},
		"resourceStatusParams": map[string]interface{}{"enableResourceStatus": true, "periodicityCsiMs": 5},
	}
	cleanUp := prepareTempConfigForTest(t, yamlMap)
	defer cleanUp()

	_, cErr := ParseConfiguration()
	assert.EqualError(t, cErr, "#configuration.fillRnibConfig - failed to fill RNib configuration: The entry 'rnib' not found\n")
}

func TestResourceStatusParamsConfigNotFoundFailure(t *testing.T) {

	yamlMap := map[string]interface{}{
		"rmr":     map[string]interface{}{"port": 6942, "maxMsgSize": 4096},
		"logging": map[string]interface{}{"logLevel": "info"},
		"http":    map[string]interface{}{"port": 631},
		"rnib":    map[string]interface{}{"maxRnibConnectionAttempts": 3, "rnibRetryIntervalMs": 10},
	}
	cleanUp := prepareTempConfigForTest(t, yamlMap)
	defer cleanUp()

	_, cErr := ParseConfiguration()
	assert.EqualError(t, cErr, "#configuration.fillResourceStatusParamsConfig - failed to fill resourceStatusParams configuration: The entry 'resourceStatusParams' not found\n")
}

func TestCharacteristicsConfigInvalidPeriodicityMs(t *testing.T) {

	yamlMap := map[string]interface{}{
		"rmr":     map[string]interface{}{"port": 6942, "maxMsgSize": 4096},
		"logging": map[string]interface{}{"logLevel": "info"},
		"http":    map[string]interface{}{"port": 631},
		"rnib":    map[string]interface{}{"maxRnibConnectionAttempts": 3, "rnibRetryIntervalMs": 10},
		"resourceStatusParams": map[string]interface{}{"enableResourceStatus": true, "periodicityMs": 50, "periodicityRsrpMeasurementMs": 480, "periodicityCsiMs": 20},
	}
	cleanUp := prepareTempConfigForTest(t, yamlMap)
	defer cleanUp()

	_, cErr := ParseConfiguration()
	assert.Error(t, cErr)
}

func TestResourceStatusParamsConfigInvalidPeriodicityRsrpMeasurementMs(t *testing.T) {

	yamlMap := map[string]interface{}{
		"rmr":     map[string]interface{}{"port": 6942, "maxMsgSize": 4096},
		"logging": map[string]interface{}{"logLevel": "info"},
		"http":    map[string]interface{}{"port": 631},
		"rnib":    map[string]interface{}{"maxRnibConnectionAttempts": 3, "rnibRetryIntervalMs": 10},
		"resourceStatusParams": map[string]interface{}{"enableResourceStatus": true, "periodicityMs": 1000, "periodicityRsrpMeasurementMs": 50, "periodicityCsiMs": 20},
	}
	cleanUp := prepareTempConfigForTest(t, yamlMap)
	defer cleanUp()

	_, cErr := ParseConfiguration()
	assert.Error(t, cErr)
}

func TestResourceStatusParamsConfigInvalidPeriodicityCsiMs(t *testing.T) {

	yamlMap := map[string]interface{}{
		"rmr":     map[string]interface{}{"port": 6942, "maxMsgSize": 4096},
		"logging": map[string]interface{}{"logLevel": "info"},
		"http":    map[string]interface{}{"port": 631},
		"rnib":    map[string]interface{}{"maxRnibConnectionAttempts": 3, "rnibRetryIntervalMs": 10},
		"resourceStatusParams": map[string]interface{}{"enableResourceStatus": true, "periodicityMs": 1000, "periodicityRsrpMeasurementMs": 480, "periodicityCsiMs": 50},
	}
	cleanUp := prepareTempConfigForTest(t, yamlMap)
	defer cleanUp()

	_, cErr := ParseConfiguration()
	assert.Error(t, cErr)
}

func TestConfigurationString(t *testing.T) {
	config, err := ParseConfiguration()
	if err != nil {
		t.Errorf("failed to parse configuration. error: %s", err)
	}
	str := config.String()
	assert.NotEmpty(t, str)
	assert.Contains(t, str, "logging")
	assert.Contains(t, str, "http")
	assert.Contains(t, str, "rmr")
	assert.Contains(t, str, "rnib")
	assert.Contains(t, str, "resourceStatusParams")
}

func prepareTempConfigForTest(t *testing.T, yamlMap map[string]interface{}) func() {
	configPath := "../resources/configuration.yaml"
	configPathTmp := "../resources/configuration.yaml_tmp"
	err := os.Rename(configPath, configPathTmp)
	if err != nil {
		t.Errorf("#http_server_test.TestRnibConfigNotFoundFailure - failed to rename configuration file: %s\n", configPath)
	}
	buf, err := yaml.Marshal(yamlMap)
	if err != nil {
		t.Errorf("#http_server_test.TestRnibConfigNotFoundFailure - failed to marshal configuration map\n")
	}
	err = ioutil.WriteFile("../resources/configuration.yaml", buf, 0644)
	if err != nil {
		t.Errorf("#http_server_test.TestRnibConfigNotFoundFailure - failed to write configuration file: %s\n", configPath)
	}

	return func() {
		err = os.Rename(configPathTmp, configPath)
		if err != nil {
			t.Errorf("#http_server_test.TestRnibConfigNotFoundFailure - failed to rename configuration file: %s\n", configPath)
		}
	}
}
