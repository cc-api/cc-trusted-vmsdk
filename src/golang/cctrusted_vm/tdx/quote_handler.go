package tdx

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"github.com/cc-api/cc-trusted-api/common/golang/cctrusted_base/tdx"
)

type QuoteHandler interface {
	// Quote gets the quote of the td vm, which is refered as cc report
	Quote([tdx.TD_REPORT_LEN]byte) ([]byte, error)
	// TdReport gets the td report of the td vm, where nonce ad userData
	// are encoded in base64
	TdReport(nonce, userData string) ([tdx.TD_REPORT_LEN]byte, error)
}

var _ QuoteHandler = (*QuoteHandler15)(nil)

type QuoteHandler15 struct {
	devicePath          string
	getTdReportOperator uintptr
	getTdQuoteOperator  uintptr
}

// TdReport implements QuoteHandler.
func (q *QuoteHandler15) TdReport(nonce, userData string) ([tdx.TD_REPORT_LEN]byte, error) {
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

	// encode nonce and userData
	hasher := sha512.New()
	if len(nonce) > 0 {
		nonceDecoded, err := base64.StdEncoding.DecodeString(nonce)
		if err != nil {
			return tdreport, err
		}
		hasher.Write(nonceDecoded)
	}

	if len(userData) > 0 {
		userDataDecoded, err := base64.StdEncoding.DecodeString(userData)
		if err != nil {
			return tdreport, err
		}
		hasher.Write(userDataDecoded)
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
	var file *os.File

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
		return &QuoteHandler15{
			devicePath:          spec.DevicePath,
			getTdReportOperator: spec.AllowedOperation[tdx.GetTdReport],
			getTdQuoteOperator:  spec.AllowedOperation[tdx.GetQuote],
		}, nil
	}
	return nil, errors.New("no supported version of tdx device")
}
