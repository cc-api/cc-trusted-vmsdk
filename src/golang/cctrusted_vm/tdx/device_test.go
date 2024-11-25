package tdx

import (
	"bytes"
	"fmt"
	"log"
	"testing"

	"github.com/cc-api/evidence-api/common/golang/evidence_api/tdx"

	"github.com/stretchr/testify/assert"
)

func TestReport15(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)

	spec := tdx.TdxDeviceSpecs[tdx.TDX_VERSION_1_5_DEVICE]
	device := TDXDevice{spec: spec}
	res := device.ProbeDevice()
	t.Log(buf.String())
	fmt.Println(res)
	assert.Equal(t, true, res)

	nonce := []byte{"IXUKoBO1UM3c1wopN4sY"}
	userData := []byte{"MTIzNDU2NzgxMjM0NTY3ODEyMzQ1Njc4MTIzNDU2NzgxMjM0NTY3ODEyMzQ1Njc4"}
	tdreport, err := device.TdReport(nonce, userData)
	assert.Nil(t, err)
	t.Log(tdreport)

	quote, err := device.Quote(tdreport)
	assert.Nil(t, err)
	assert.NotNil(t, quote)
}
