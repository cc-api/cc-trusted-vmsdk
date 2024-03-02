package vmsdk

import (
	"errors"
	"fmt"
	"sync"

	cctrusted_vm "github.com/cc-api/cc-trusted-vmsdk/src/golang/cctrusted_vm"
	_ "github.com/cc-api/cc-trusted-vmsdk/src/golang/cctrusted_vm/tdx"

	"github.com/cc-api/cc-trusted-api/common/golang/cctrusted_base"
	"github.com/cc-api/cc-trusted-api/common/golang/cctrusted_base/tdx"
)

var _ cctrusted_base.CCTrustedAPI = (*SDK)(nil)

type SDK struct {
	cvm cctrusted_vm.ConfidentialVM
}

// DumpCCReport implements cctrusted_base.CCTrustedAPI.
func (s *SDK) DumpCCReport(reportBytes []byte) error {
	vmCtx := s.cvm.CVMContext()
	switch vmCtx.VMType {
	case cctrusted_base.TYPE_CC_TDX:
		report, err := tdx.NewTdxReportFromBytes(reportBytes)
		if err != nil {
			return err
		}
		report.Dump(cctrusted_base.QuoteDumpFormatHuman)
	default:
	}
	return nil
}

// GetCCMeasurement implements cctrusted_base.CCTrustedAPI.
func (s *SDK) GetCCMeasurement(index int, alg cctrusted_base.TCG_ALG) (cctrusted_base.TcgDigest, error) {
	emptyRet := cctrusted_base.TcgDigest{}
	report, err := s.GetCCReport("", "", nil)
	if err != nil {
		return emptyRet, err
	}
	group := report.IMRGroup()
	if index > group.MaxIndex {
		return emptyRet, fmt.Errorf("index %d larger than max index %d", index, group.MaxIndex)
	}
	entry := group.Group[index]
	if entry.AlgID != alg {
		return emptyRet, fmt.Errorf("required alg %v, but got alg %v", alg, entry.AlgID)
	}
	return entry, nil
}

// GetMeasurementCount implements cctrusted_base.CCTrustedAPI.
func (s *SDK) GetMeasurementCount() (int, error) {
	return s.cvm.MaxImrIndex() + 1, nil
}

// ReplayCCEventLog implements cctrusted_base.CCTrustedAPI.
func (s *SDK) ReplayCCEventLog(formatedEventLogs []cctrusted_base.FormatedTcgEvent) map[int]map[cctrusted_base.TCG_ALG][]byte {
	return cctrusted_base.ReplayFormatedEventLog(formatedEventLogs)
}

// GetDefaultAlgorithm implements cctrusted_base.CCTrustedAPI.
func (s *SDK) GetDefaultAlgorithm() cctrusted_base.TCG_ALG {
	return s.cvm.DefaultAlgorithm()
}

// SelectEventlog implements CCTrustedAPI.
func (s *SDK) GetCCEventLog(start int32, count int32) (*cctrusted_base.EventLogger, error) {
	el, err := s.internelEventlog()
	if err != nil {
		return nil, err
	}
	el.Parse()

	if start != 0 || count != 0 {
		el, err = el.Select(int(start), int(count))
		if err != nil {
			return nil, err
		}
	}

	return el, nil
}

func (s *SDK) internelEventlog() (*cctrusted_base.EventLogger, error) {
	if s.cvm == nil {
		return nil, errors.New("no available cvm in sdk")
	}

	eventLogBytes, err := s.cvm.FullEventLog()
	if err != nil {
		return nil, err
	}

	imaLogBytes, err := s.cvm.FullIMALog()
	if err != nil {
		return nil, err
	}

	el := cctrusted_base.NewEventLogger(eventLogBytes, imaLogBytes, cctrusted_base.TCG_PCCLIENT_FORMAT)
	return el, nil
}

// Report implements CCTrustedAPI.
func (s *SDK) GetCCReport(nonce, userData string, _ any) (cctrusted_base.Report, error) {
	if s.cvm == nil {
		return nil, errors.New("no available cvm in sdk")
	}

	reportBytes, err := s.cvm.Report(nonce, userData)
	if err != nil {
		return nil, err
	}

	vmCtx := s.cvm.CVMContext()
	switch vmCtx.VMType {
	case cctrusted_base.TYPE_CC_TDX:
		report, err := tdx.NewTdxReportFromBytes(reportBytes)
		if err != nil {
			return nil, err
		}
		return report, nil
	default:
	}
	return nil, errors.New("parse" + vmCtx.VMType.String() + "report failed")
}

var once sync.Once

var instance SDK
var sdkErr error

func GetSDKInstance(args *cctrusted_vm.CVMInitArgs) (SDK, error) {
	once.Do(func() {
		cvm, err := cctrusted_vm.GetCVMInstance(args)
		if err != nil {
			sdkErr = err
		} else {
			instance.cvm = cvm
		}
	})
	return instance, sdkErr
}
