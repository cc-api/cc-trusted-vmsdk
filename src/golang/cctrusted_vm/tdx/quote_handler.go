package tdx

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"strconv"
	"syscall"
	"time"
	"unsafe"

	"github.com/cc-api/evidence-api/common/golang/evidence_api/tdx"
	"github.com/mdlayher/vsock"
)

type QuoteHandler interface {
	// Quote gets the quote of the td vm, which is refered as cc report
	Quote([tdx.TD_REPORT_LEN]byte) ([]byte, error)
	// TdReport gets the td report of the td vm, where nonce ad userData
	// are encoded in base64
	TdReport(nonce, userData []byte) ([tdx.TD_REPORT_LEN]byte, error)
}

var _ QuoteHandler = (*QuoteHandler15)(nil)

type QuoteHandler15 struct {
	devicePath          string
	tdxAttestConfig     map[string]string
	getTdReportOperator uintptr
	getTdQuoteOperator  uintptr
}

// TdReport implements QuoteHandler.
func (q *QuoteHandler15) TdReport(nonce, userData []byte) ([tdx.TD_REPORT_LEN]byte, error) {
	tdreport := [tdx.TD_REPORT_LEN]uint8{}
	var err error
	var file *os.File

	defer func() {
		if file != nil {
			file.Close()
		}
	}()

	// open tdx device
	file, err = os.Open(q.devicePath)
	if err != nil {
		return tdreport, err
	}

	// concatenate nonce and userData
	// check if the data is base64 encoded, if yes, decode before doing hash
	hasher := sha512.New()
	if len(nonce) > 0 {
		nonceDecoded, err := base64.StdEncoding.DecodeString(string(nonce))
		if err != nil {
			hasher.Write(nonce)
		} else {
			hasher.Write(nonceDecoded)
		}
	}

	if len(userData) > 0 {
		userDataDecoded, err := base64.StdEncoding.DecodeString(string(userData))
		if err != nil {
			hasher.Write(userData)
		} else {
			hasher.Write(userDataDecoded)
		}
	}

	reportData := [64]byte(hasher.Sum(nil))
	req := tdx.TdxReportReq15{
		ReportData: reportData,
		Tdreport:   tdreport,
	}

	argPtr := uintptr(unsafe.Pointer(&req))

	// get td report via ioctl
	if _, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(file.Fd()),
		q.getTdReportOperator,
		argPtr,
	); errno != 0 {
		return tdreport, errno
	}

	return req.Tdreport, nil
}

// Quote implements QuoteHandler.
func (q *QuoteHandler15) Quote(tdreport [tdx.TD_REPORT_LEN]byte) ([]byte, error) {
	var err error
	var quote []byte

	if len(q.tdxAttestConfig) != 0 {
		if val, ok := q.tdxAttestConfig["port"]; ok {
			val_int, _ := strconv.Atoi(val)
			quote, _ = q.FetchQuoteByVsock(val_int, tdreport)
		}
	}

	if quote == nil {
		quote, err = q.FetchQuoteByTdvmcall(tdreport)
	}

	return quote, err
}

func (q *QuoteHandler15) FetchQuoteByVsock(vsockPort int, tdreport [tdx.TD_REPORT_LEN]byte) ([]byte, error) {
	// host context id for vsock connection
	const hostCid = 2

	// connect to QGS socket
	conn, err := vsock.Dial(hostCid, uint32(vsockPort), nil)
	if err != nil {
		return nil, err
	}

	// set deadline for connection
	err = conn.SetDeadline(time.Now().Add(30 * time.Second))
	if err != nil {
		return nil, err
	}

	// create tdx quote request
	headerSize := 4
	qgsMsgGetQuoteReq := tdx.NewQgsMsgGetQuoteReqVer15(tdreport).Bytes()

	msgSize := make([]byte, headerSize)
	binary.BigEndian.PutUint32(msgSize, uint32(len(qgsMsgGetQuoteReq)))

	pBlobPayload := make([]byte, len(qgsMsgGetQuoteReq)+headerSize)
	copy(pBlobPayload[:headerSize], msgSize)
	copy(pBlobPayload[headerSize:], qgsMsgGetQuoteReq[:len(qgsMsgGetQuoteReq)])

	_, err = conn.Write(pBlobPayload)
	if err != nil {
		return nil, err
	}

	// fetch QGS response size
	header := make([]byte, headerSize)
	nRead, err := conn.Read(header)
	if err != nil {
		return nil, err
	}
	var size uint32
	size = 0
	for _, item := range header[:nRead] {
		size = size*256 + uint32(item&0xFF)
	}

	// fetch QGS response
	qgsResp := make([]byte, size)
	nRead, err = conn.Read(qgsResp)
	if err != nil {
		return nil, err
	}

	// close connection
	if err = conn.Close(); err != nil {
		return nil, err
	}

	resp := tdx.NewQgsMsgGetQuoteRespFromBytes(qgsResp[:nRead])
	return resp.IdQuote[:resp.QuoteSize], nil
}

func (q *QuoteHandler15) FetchQuoteByTdvmcall(tdreport [tdx.TD_REPORT_LEN]byte) ([]byte, error) {
	var file *os.File
	var err error

	defer func() {
		if file != nil {
			file.Close()
		}
	}()

	// open tdx device
	file, err = os.Open(q.devicePath)
	if err != nil {
		return nil, err
	}

	// create tdx quote request
	qsgMsgGetQuoteReq := tdx.NewQgsMsgGetQuoteReqVer15(tdreport)
	tdxQuoteHdr := tdx.NewTdxQuoteHdrVer15(qsgMsgGetQuoteReq)
	tdxQuoteReq := tdx.NewTdxQuoteReqVer15(tdxQuoteHdr)

	argPtr := uintptr(unsafe.Pointer(tdxQuoteReq))

	// get tdx quote via ioctl
	if _, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(file.Fd()),
		q.getTdQuoteOperator,
		argPtr,
	); errno != 0 {
		return nil, errno
	}

	if tdxQuoteHdr.Status != 0 {
		return nil, fmt.Errorf("get quote failed! status code 0x%x", tdxQuoteHdr.Status)
	}

	dataLen := binary.BigEndian.Uint32(tdxQuoteHdr.DataLenBeBytes[:])
	if uint64(tdxQuoteHdr.OutLen)-4 != uint64(dataLen) {
		return nil, errors.New("td quote data length sanity check failed")
	}
	resp := tdx.NewQgsMsgGetQuoteRespFromBytes(tdxQuoteHdr.Data[:])
	return resp.IdQuote[:resp.QuoteSize], nil
}

func GetQuoteHandler(spec tdx.TDXDeviceSpec) (QuoteHandler, error) {
	switch spec.Version {
	case tdx.TDX_VERSION_1_0:
		// TODO: support tdx 1.0
		return nil, errors.New("tdx 1.0 version not supported now temporarily")
	case tdx.TDX_VERSION_1_5:
		attestConfig := spec.ProbeAttestConfig()
		return &QuoteHandler15{
			devicePath:          spec.DevicePath,
			tdxAttestConfig:     attestConfig,
			getTdReportOperator: spec.AllowedOperation[tdx.GetTdReport],
			getTdQuoteOperator:  spec.AllowedOperation[tdx.GetQuote],
		}, nil
	}
	return nil, errors.New("no supported version of tdx device")
}
