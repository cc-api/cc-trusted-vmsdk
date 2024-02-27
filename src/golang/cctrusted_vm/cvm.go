package vmsdk

import (
	"errors"

	"github.com/cc-api/cc-trusted-api/common/golang/cctrusted_base"
)

type Device interface {
	ProbeDevice() error
	Report(nonce, userDate string) ([]byte, error)
	Name() string
	CCType() cctrusted_base.CC_Type
	Version() cctrusted_base.DeviceVersion
}

type EventRecorder interface {
	ProbeRecorder() error
	FullEventLog() ([]byte, error)
}

type CVMContext struct {
	VMType  cctrusted_base.CC_Type
	Version cctrusted_base.DeviceVersion
}

type ConfidentialVM interface {
	Probe() error
	CVMContext() CVMContext
	MaxImrIndex() int
	DefaultAlgorithm() cctrusted_base.TCG_ALG
	Device
	EventRecorder
	cctrusted_base.IMARecorder
}

type CVMInitArgs struct {
	// RedirectedAcpiTableFile is the alternative
	// of the original `DEFAULT_ACPI_TABLE_FILE`, if which
	// can not be accessed
	RedirectedAcpiTableFile string
	// RedirectedAcpiTableDataFile is the alternative
	// of the original `DEFAULT_ACPI_TABLE_DATA_FILE`, if which
	// can not be accessed
	RedirectedAcpiTableDataFile string
}

type CVMInitFunc func(*CVMInitArgs) (ConfidentialVM, error)

var cvmInitFuncs []CVMInitFunc

func RegisterCVMInitFunc(fn CVMInitFunc) {
	cvmInitFuncs = append(cvmInitFuncs, fn)
}

func GetCVMInstance(args *CVMInitArgs) (ConfidentialVM, error) {
	for _, fn := range cvmInitFuncs {
		cvm, err := fn(args)
		if err != nil {
			continue
		}
		return cvm, nil
	}
	return nil, errors.New("no available confidential vm")
}
