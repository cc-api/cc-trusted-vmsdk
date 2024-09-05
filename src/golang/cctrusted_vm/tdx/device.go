package tdx

import (
	"errors"
	"os"

	"github.com/cc-api/evidence-api/common/golang/evidence_api/tdx"

	"github.com/cc-api/evidence-api/common/golang/evidence_api"

	cctrusted_vm "github.com/cc-api/cc-trusted-vmsdk/src/golang/cctrusted_vm"
)

var _ cctrusted_vm.Device = (*TDXDevice)(nil)

type TDXDevice struct {
	spec tdx.TDXDeviceSpec
	cctrusted_vm.GenericDevice
	QuoteHandler
}

// Version implements cctrusted_vm.Device.
func (t *TDXDevice) Version() evidence_api.DeviceVersion {
	return t.spec.Version
}

// CCType implements cctrusted_vm.Device.
func (t *TDXDevice) CCType() evidence_api.CC_Type {
	return evidence_api.TYPE_CC_TDX
}

// Name implements cctrusted_vm.Device.
func (t *TDXDevice) Name() string {
	return t.spec.DevicePath
}

// ProbeDevice implements cctrusted_vm.Device, probe valid tdx device.
func (t *TDXDevice) ProbeDevice() error {
	for _, spec := range tdx.TdxDeviceSpecs {
		_, err := os.Stat(spec.DevicePath)
		if err != nil {
			continue
		}

		t.spec = spec
		err = t.initDevice()
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("no valid tdx device found")
}

func (t *TDXDevice) initDevice() error {
	h, err := GetQuoteHandler(t.spec)
	if err != nil {
		return err
	}
	t.QuoteHandler = h
	return nil
}

// Report implements cctrusted_vm.Device, get CC report
func (t *TDXDevice) Report(nonce, userData []byte, extraArgs map[string]any) (evidence_api.CcReport, error) {
	var resp evidence_api.CcReport
	var err error

	// call parent Report() func to retrieve cc report using Configfs-tsm
	resp, err = t.GenericDevice.Report(nonce, userData, extraArgs)
	if err == nil {
		return resp, nil
	}

	// get tdx report
	tdreport, err := t.TdReport(nonce, userData)
	if err != nil {
		return evidence_api.CcReport{}, err
	}
	// get tdx quote, aka. CC report
	quote, err := t.Quote(tdreport)
	if err != nil {
		return evidence_api.CcReport{}, err
	}

	resp = evidence_api.CcReport{
		Outblob: quote,
	}

	return resp, nil
}
