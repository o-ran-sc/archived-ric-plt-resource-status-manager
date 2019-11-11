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

package converters

// #cgo CFLAGS: -I../asn1codec/inc/ -I../asn1codec/e2ap_engine/
// #cgo LDFLAGS: -L ../asn1codec/lib/ -L../asn1codec/e2ap_engine/ -le2ap_codec -lasncodec
// #include <asn1codec_utils.h>
import "C"
import (
	"fmt"
	"github.com/pkg/errors"
	"rsm/logger"
	"unsafe"
)

type X2apPduUnpacker struct {
	logger *logger.Logger

}

func NewX2apPduUnpacker(logger *logger.Logger) X2apPduUnpacker {
	return X2apPduUnpacker{logger :logger}
}

func (r X2apPduUnpacker) UnpackX2apPdu(packedBufferSize int, packedBuf []byte, maxMessageBufferSize int) (*C.E2AP_PDU_t, error) {
	pdu := C.new_pdu()

	if pdu == nil {
		return nil, errors.New("allocation failure (pdu)")
	}

	r.logger.Infof("#x2apPdu_asn1_unpacker.UnpackX2apPdu - Packed pdu(%d):%x", packedBufferSize, packedBuf)

	errBuf := make([]C.char, maxMessageBufferSize)
	if !C.per_unpack_pdu(pdu, C.ulong(packedBufferSize), (*C.uchar)(unsafe.Pointer(&packedBuf[0])), C.ulong(len(errBuf)), &errBuf[0]) {
		return nil, errors.New(fmt.Sprintf("unpacking error: %s", C.GoString(&errBuf[0])))
	}

	if r.logger.DebugEnabled() {
		C.asn1_pdu_printer(pdu, C.size_t(len(errBuf)), &errBuf[0])
		r.logger.Debugf("#x2apPdu_asn1_unpacker.UnpackX2apPdu - PDU: %v  packed size:%d", C.GoString(&errBuf[0]), packedBufferSize)
	}

	return pdu, nil
}

func (r X2apPduUnpacker)UnpackX2apPduAsString(packedBufferSize int, packedBuf []byte, maxMessageBufferSize int) (string, error) {
	pdu, err := r.UnpackX2apPdu(packedBufferSize, packedBuf, maxMessageBufferSize)
	if err != nil {
		return "", err
	}

	defer C.delete_pdu(pdu)

	buf := make([]C.char, 16*maxMessageBufferSize)
	C.asn1_pdu_printer(pdu, C.size_t(len(buf)), &buf[0])
	return C.GoString(&buf[0]), nil
}



