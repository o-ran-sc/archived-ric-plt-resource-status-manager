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

package services

import (
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader"
	"net"
	"rsm/configuration"
	"rsm/logger"
	"rsm/models"
	"rsm/rsmdb"
	"time"
)

type RNibDataService interface {
	GetRsmGeneralConfiguration() (*models.RsmGeneralConfiguration, error)
	GetRsmRanInfo(ranName string) (*models.RsmRanInfo, error)
	SaveRsmRanInfo(rsmData *models.RsmRanInfo) error
	GetNodeb(ranName string) (*entities.NodebInfo, error)
	GetListNodebIds() ([]*entities.NbIdentity, error)
	PingRnib() bool
}

type rNibDataService struct {
	logger        *logger.Logger
	rnibReader    reader.RNibReader
	rsmReader     rsmdb.RsmReader
	rsmWriter     rsmdb.RsmWriter
	maxAttempts   int
	retryInterval time.Duration
}

func NewRnibDataService(logger *logger.Logger, config *configuration.Configuration, rnibReader reader.RNibReader, rsmReader rsmdb.RsmReader, rsmWriter rsmdb.RsmWriter) *rNibDataService {
	return &rNibDataService{
		logger:        logger,
		rnibReader:    rnibReader,
		rsmReader:     rsmReader,
		rsmWriter:     rsmWriter,
		maxAttempts:   config.Rnib.MaxRnibConnectionAttempts,
		retryInterval: time.Duration(config.Rnib.RnibRetryIntervalMs) * time.Millisecond,
	}
}

func (w *rNibDataService) GetRsmGeneralConfiguration() (*models.RsmGeneralConfiguration, error) {
	w.logger.Infof("#RnibDataService.GetRsmGeneralConfiguration")

	var rsmGeneralConfiguration *models.RsmGeneralConfiguration = nil

	err := w.retry("GetRsmGeneralConfiguration", func() (err error) {
		rsmGeneralConfiguration, err = w.rsmReader.GetRsmGeneralConfiguration()
		return
	})

	return rsmGeneralConfiguration, err
}

func (w *rNibDataService) GetRsmRanInfo(ranName string) (*models.RsmRanInfo, error) {
	w.logger.Infof("#RnibDataService.GetRsmRanInfo - RAN name: %s", ranName)

	var rsmData *models.RsmRanInfo = nil

	err := w.retry("GetRsmRanInfo", func() (err error) {
		rsmData, err = w.rsmReader.GetRsmRanInfo(ranName)
		return
	})

	return rsmData, err
}

func (w *rNibDataService) SaveRsmRanInfo(rsmRanInfo *models.RsmRanInfo) error {
	w.logger.Infof("#RnibDataService.SaveRsmRanInfo - RAN name: %s", rsmRanInfo.RanName)

	err := w.retry("SaveRsmRanInfo", func() (err error) {
		err = w.rsmWriter.SaveRsmRanInfo(rsmRanInfo)
		return
	})

	return err
}

func (w *rNibDataService) GetNodeb(ranName string) (*entities.NodebInfo, error) {
	w.logger.Infof("#RnibDataService.GetNodeb - ranName: %s", ranName)

	var nodeb *entities.NodebInfo = nil

	err := w.retry("GetNodeb", func() (err error) {
		nodeb, err = w.rnibReader.GetNodeb(ranName)
		return
	})

	return nodeb, err
}

func (w *rNibDataService) GetListNodebIds() ([]*entities.NbIdentity, error) {
	w.logger.Infof("#RnibDataService.GetListNodebIds")

	var nodeIds []*entities.NbIdentity = nil

	err := w.retry("GetListNodebIds", func() (err error) {
		nodeIds, err = w.rnibReader.GetListNodebIds()
		return
	})

	return nodeIds, err
}

func (w *rNibDataService) PingRnib() bool {
	err := w.retry("GetListNodebIds", func() (err error) {
		_, err = w.rnibReader.GetListNodebIds()
		return
	})

	return !isRnibConnectionError(err)
}

func (w *rNibDataService) retry(rnibFunc string, f func() error) (err error) {
	attempts := w.maxAttempts

	for i := 1; ; i++ {
		err = f()
		if err == nil {
			return
		}
		if !isRnibConnectionError(err) {
			return err
		}
		if i >= attempts {
			w.logger.Errorf("#RnibDataService.retry - after %d attempts of %s, last error: %s", attempts, rnibFunc, err)
			return err
		}
		time.Sleep(w.retryInterval)

		w.logger.Infof("#RnibDataService.retry - retrying %d %s after error: %s", i, rnibFunc, err)
	}
}

func isRnibConnectionError(err error) bool {
	internalErr, ok := err.(*common.InternalError)
	if !ok {
		return false
	}
	_, ok = internalErr.Err.(*net.OpError)

	return ok
}
