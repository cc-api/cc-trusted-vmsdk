package tdx

import (
	cctrusted_vm "github.com/cc-api/cc-trusted-vmsdk/src/golang/cctrusted_vm"

	"github.com/cc-api/cc-trusted-api/common/golang/cctrusted_base"
	"github.com/cc-api/cc-trusted-api/common/golang/cctrusted_base/tdx"
)

var _ cctrusted_vm.ConfidentialVM = (*TdxVM)(nil)

type TdxVM struct {
	cctrusted_vm.Device
	cctrusted_vm.EventRecorder
	cctrusted_base.IMARecorder
}

func NewTdxVM(args *cctrusted_vm.CVMInitArgs) *TdxVM {
	vm := &TdxVM{
		Device:      &TDXDevice{},
		IMARecorder: &cctrusted_base.DefaultIMARecorder{},
	}
	r := &TDXEventLogRecorder{}
	if args != nil {
		if args.RedirectedAcpiTableFile != "" {
			r.RedirectAcpiTableFile(args.RedirectedAcpiTableFile)
		}
		if args.RedirectedAcpiTableDataFile != "" {
			r.RedirectAcpiTableDataFile(args.RedirectedAcpiTableDataFile)
		}
	}
	vm.EventRecorder = r
	return vm
}

// DefaultAlgorithm implements cctrusted_vm.ConfidentialVM.
func (t *TdxVM) DefaultAlgorithm() cctrusted_base.TCG_ALG {
	return cctrusted_base.TPM_ALG_SHA384
}

// MaxImrIndex implements cctrusted_vm.ConfidentialVM.
func (t *TdxVM) MaxImrIndex() int {
	return tdx.RTMRMaxIndex
}

// CVMContext implements cctrusted_vm.ConfidentialVM.
func (t *TdxVM) CVMContext() cctrusted_vm.CVMContext {
	return cctrusted_vm.CVMContext{
		VMType:  t.CCType(),
		Version: t.Version(),
	}
}

// Probe implements cctrusted_vm.ConfidentialVM,
// probing tdx device, eventlog and ima
func (t *TdxVM) Probe() error {
	// probe tdx device
	if err := t.ProbeDevice(); err != nil {
		return err
	}

	// probe eventlog
	if err := t.ProbeRecorder(); err != nil {
		return err
	}

	// probe ima
	if err := t.ProbeIMARecorder(); err != nil {
		return err
	}
	return nil
}

// TdxVMInitFunc creates and inits a tdx confidential VM
func TdxVMInitFunc(args *cctrusted_vm.CVMInitArgs) (cctrusted_vm.ConfidentialVM, error) {
	tdx := NewTdxVM(args)
	if err := tdx.Probe(); err != nil {
		return nil, err
	}
	return tdx, nil
}

func init() {
	cctrusted_vm.RegisterCVMInitFunc(TdxVMInitFunc)
}
