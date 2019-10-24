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

package rmrcgo

// #cgo LDFLAGS: -L/usr/local/lib -lrmr_nng -lnng
// #include <rmr/rmr.h>
// #include <stdlib.h>
import "C"
import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
	"time"
	"unsafe"

	"rsm/logger"
)

func (*Context) Init(readyIntervalSec int, port string, maxMsgSize int, flags int, logger *logger.Logger) RmrMessenger {
	pp := C.CString(port)
	defer C.free(unsafe.Pointer(pp))
	logger.Debugf("#rmr_c_go_api.Init - Going to initiate RMR router")
	ctx := NewContext(maxMsgSize, flags, C.rmr_init(pp, C.int(maxMsgSize), C.int(flags)), logger)
	start := time.Now()
	ticker := time.NewTicker(time.Duration(readyIntervalSec) * time.Second)
	for ; !ctx.IsReady(); <-ticker.C {
		if time.Since(start) >= time.Minute {
			logger.Debugf("#rmr_c_go_api.Init - Routing table is not ready")
			start = time.Now()
		}
	}
	logger.Infof("#rmr_c_go_api.Init - RMR router has been initiated")

	// Configure the rmr to make rounds of attempts to send a message before notifying the application that it should retry.
	// Each round is about 1000 attempts with a short sleep between each round.
	C.rmr_set_stimeout(ctx.RmrCtx, C.int(1000))
	return ctx
}

func (ctx *Context) SendMsg(msg *MBuf) (*MBuf, error) {
	ctx.checkContextInitialized()
	ctx.Logger.Debugf("#rmr_c_go_api.SendMsg - Going to send message. MBuf: %v", *msg)
	allocatedCMBuf := ctx.getAllocatedCRmrMBuf(ctx.Logger, msg, ctx.MaxMsgSize)
	defer C.rmr_free_msg(allocatedCMBuf)
	state := allocatedCMBuf.state
	if state != RMR_OK {
		errorMessage := fmt.Sprintf("#rmr_c_go_api.SendMsg - Failed to get allocated message. state: %v - %s", state, states[int(state)])
		return nil, errors.New(errorMessage)
	}

	transactionId := string(*msg.XAction)
	tmpTid := strings.TrimSpace(transactionId)
	ctx.Logger.Infof("[RSM -> RMR] #rmr_c_go_api.SendMsg - Going to send message %v for transaction id: %s", *msg, tmpTid)

	currCMBuf := C.rmr_send_msg(ctx.RmrCtx, allocatedCMBuf)
	state = currCMBuf.state

	if ctx.Logger.DebugEnabled() {
		ctx.Logger.Debugf("#rmr_c_go_api.SendMsg - The current message  state: %v, message buffer:%v", state, currCMBuf)
	}

	if state != RMR_OK {
		errorMessage := fmt.Sprintf("#rmr_c_go_api.SendMsg - Failed to send message. state: %v - %s", state, states[int(state)])
		return nil, errors.New(errorMessage)
	}

	ctx.Logger.Debugf("#rmr_c_go_api.SendMsg - The message has been sent successfully ")
	return convertToMBuf(currCMBuf), nil
}

func (ctx *Context) RecvMsg() (*MBuf, error) {
	ctx.checkContextInitialized()
	ctx.Logger.Debugf("#rmr_c_go_api.RecvMsg - Going to receive message")
	allocatedCMBuf := C.rmr_alloc_msg(ctx.RmrCtx, C.int(ctx.MaxMsgSize))
	defer C.rmr_free_msg(allocatedCMBuf)

	currCMBuf := C.rmr_rcv_msg(ctx.RmrCtx, allocatedCMBuf)
	state := currCMBuf.state

	if state != RMR_OK {
		errorMessage := fmt.Sprintf("#rmr_c_go_api.RecvMsg - Failed to receive message. state: %v - %s", state, states[int(state)])
		ctx.Logger.Errorf(errorMessage)
		return nil, errors.New(errorMessage)
	}

	mbuf := convertToMBuf(currCMBuf)
	transactionId := string(*mbuf.XAction)
	tmpTid := strings.TrimSpace(transactionId)
	ctx.Logger.Infof("[RMR ->RSM] #rmr_c_go_api.RecvMsg - message %v has been received for transaction id: %s", *mbuf, tmpTid)
	return mbuf, nil
}

func (ctx *Context) IsReady() bool {
	ctx.Logger.Debugf("#rmr_c_go_api.IsReady - Going to check if routing table is initialized")
	return int(C.rmr_ready(ctx.RmrCtx)) != 0
}

func (ctx *Context) Close() {
	ctx.Logger.Debugf("#rmr_c_go_api.Close - Going to close RMR context")
	C.rmr_close(ctx.RmrCtx)
	time.Sleep(100 * time.Millisecond)
}
