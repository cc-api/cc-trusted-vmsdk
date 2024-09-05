package vmsdk

import (
	"errors"
	"fmt"
	"log"
	"sync"

	cctrusted_vm "github.com/cc-api/cc-trusted-vmsdk/src/golang/cctrusted_vm"
	_ "github.com/cc-api/cc-trusted-vmsdk/src/golang/cctrusted_vm/tdx"

	"github.com/cc-api/evidence-api/common/golang/evidence_api"
	"github.com/cc-api/evidence-api/common/golang/evidence_api/tdx"
)

var _ evidence_api.EvidenceAPI = (*SDK)(nil)

type SDK struct {
	cvm cctrusted_vm.ConfidentialVM
}

// DumpCCReport implements evidence_api.EvidenceAPI.
func (s *SDK) DumpCCReport(reportBytes []byte) error {
	vmCtx := s.cvm.CVMContext()
	switch vmCtx.VMType {
	case evidence_api.TYPE_CC_TDX:
		report, err := tdx.NewTdxReportFromBytes(reportBytes)
		if err != nil {
			return err
		}
		report.Dump(evidence_api.QuoteDumpFormatHuman)
	default:
	}
	return nil
}

// GetCCMeasurement implements evidence_api.EvidenceAPI.
func (s *SDK) GetCCMeasurement(index int, alg evidence_api.TCG_ALG) (evidence_api.TcgDigest, error) {
	emptyRet := evidence_api.TcgDigest{}
	report, err := s.GetCCReport(nil, nil, nil)
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

// GetMeasurementCount implements evidence_api.EvidenceAPI.
func (s *SDK) GetMeasurementCount() (int, error) {
	return s.cvm.MaxImrIndex() + 1, nil
}

// ReplayCCEventLog implements evidence_api.EvidenceAPI.
func (s *SDK) ReplayCCEventLog(formatedEventLogs []evidence_api.FormatedTcgEvent) map[int]map[evidence_api.TCG_ALG][]byte {
	return evidence_api.ReplayFormatedEventLog(formatedEventLogs)
}

// GetDefaultAlgorithm implements evidence_api.EvidenceAPI.
func (s *SDK) GetDefaultAlgorithm() (evidence_api.TCG_ALG, error) {
	return s.cvm.DefaultAlgorithm(), nil
}

// SelectEventlog implements EvidenceAPI.
func (s *SDK) GetCCEventLog(params ...int32) ([]evidence_api.FormatedTcgEvent, error) {
	el, err := s.internelEventlog()
	if err != nil {
		return nil, err
	}
	el.Parse()

	var start int32
	var count int32

	// Fetch optional params according to user specification
	if len(params) > 2 || len(params) == 0 {
		log.Fatalf("Invalid params specified. Using default values.")
		start = 0
		count = int32(len(el.EventLog()))
	} else if len(params) == 2 {
		start = params[0]
		count = params[1]
	} else {
		start = params[0]
		count = int32(len(el.EventLog())) - start
	}

	if start != 0 || count != 0 {
		el, err = el.Select(int(start), int(count))
		if err != nil {
			return nil, err
		}
	}

	return el.EventLog(), nil
}

func (s *SDK) internelEventlog() (*evidence_api.EventLogger, error) {
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

	el := evidence_api.NewEventLogger(eventLogBytes, imaLogBytes, evidence_api.TCG_PCCLIENT_FORMAT)
	return el, nil
}

// Report implements EvidenceAPI.
func (s *SDK) GetCCReport(nonce, userData []byte, extraArgs map[string]any) (evidence_api.Report, error) {
	if s.cvm == nil {
		return nil, errors.New("no available cvm in sdk")
	}

	reportStruct, err := s.cvm.Report(nonce, userData, extraArgs)
	if err != nil {
		return nil, err
	}

	vmCtx := s.cvm.CVMContext()
	switch vmCtx.VMType {
	case evidence_api.TYPE_CC_TDX:
		report, err := tdx.NewTdxReportFromBytes(reportStruct.Outblob)
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
