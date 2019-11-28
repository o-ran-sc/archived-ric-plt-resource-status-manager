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


package configuration

import (
	"fmt"
	"github.com/spf13/viper"
)

type Configuration struct {
	Logging struct {
		LogLevel string
	}
	Http struct {
		Port int
	}
	Rmr struct {
		Port             int
		MaxMsgSize       int
		ReadyIntervalSec int
	}
	Rnib struct {
		MaxRnibConnectionAttempts int
		RnibRetryIntervalMs       int
	}
}

func ParseConfiguration() (*Configuration, error) {
	viper.SetConfigType("yaml")
	viper.SetConfigName("configuration")
	viper.AddConfigPath("RSM/resources/")
	viper.AddConfigPath("./resources/")     //For production
	viper.AddConfigPath("../resources/")    //For test under Docker
	viper.AddConfigPath("../../resources/") //For test under Docker
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("#configuration.parseConfiguration - failed to read configuration file: %s\n", err)
	}

	config := Configuration{}
	if err := config.fillRmrConfig(viper.Sub("rmr")); err != nil {
		return nil, err
	}
	if err := config.fillHttpConfig(viper.Sub("http")); err != nil {
		return nil, err
	}
	if err := config.fillLoggingConfig(viper.Sub("logging")); err != nil {
		return nil, err
	}
	if err := config.fillRnibConfig(viper.Sub("rnib")); err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Configuration) fillLoggingConfig(logConfig *viper.Viper) error {
	if logConfig == nil {
		return fmt.Errorf("#configuration.fillLoggingConfig - failed to fill logging configuration: The entry 'logging' not found\n")
	}
	c.Logging.LogLevel = logConfig.GetString("logLevel")
	return nil
}

func (c *Configuration) fillHttpConfig(httpConfig *viper.Viper) error {
	if httpConfig == nil {
		return fmt.Errorf("#configuration.fillHttpConfig - failed to fill HTTP configuration: The entry 'http' not found\n")
	}
	c.Http.Port = httpConfig.GetInt("port")
	return nil
}

func (c *Configuration) fillRmrConfig(rmrConfig *viper.Viper) error {
	if rmrConfig == nil {
		return fmt.Errorf("#configuration.fillRmrConfig - failed to fill RMR configuration: The entry 'rmr' not found\n")
	}
	c.Rmr.Port = rmrConfig.GetInt("port")
	c.Rmr.MaxMsgSize = rmrConfig.GetInt("maxMsgSize")
	c.Rmr.ReadyIntervalSec = rmrConfig.GetInt("readyIntervalSec")
	return nil
}

func (c *Configuration) fillRnibConfig(rnibConfig *viper.Viper) error {
	if rnibConfig == nil {
		return fmt.Errorf("#configuration.fillRnibConfig - failed to fill RNib configuration: The entry 'rnib' not found\n")
	}
	c.Rnib.MaxRnibConnectionAttempts = rnibConfig.GetInt("maxRnibConnectionAttempts")
	c.Rnib.RnibRetryIntervalMs = rnibConfig.GetInt("rnibRetryIntervalMs")
	return nil
}


func (c *Configuration) String() string {
	return fmt.Sprintf("{logging.logLevel: %s, http.port: %d, rmr.port: %d, rmr.maxMsgSize: %d, rmr.readyIntervalSec: %d, rnib.maxRnibConnectionAttempts: %d, rnib.rnibRetryIntervalMs: %d, ",
		c.Logging.LogLevel,
		c.Http.Port,
		c.Rmr.Port,
		c.Rmr.MaxMsgSize,
		c.Rmr.ReadyIntervalSec,
		c.Rnib.MaxRnibConnectionAttempts,
		c.Rnib.RnibRetryIntervalMs,
	)
}
