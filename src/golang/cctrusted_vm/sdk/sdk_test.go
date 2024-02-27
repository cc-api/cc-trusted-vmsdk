package vmsdk

import (
	"bytes"
	"log"
	"testing"

	"github.com/cc-api/cc-trusted-api/common/golang/cctrusted_base"

	"github.com/stretchr/testify/assert"
)

func TestSDKReport(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)

	sdk, err := GetSDKInstance(nil)
	assert.Nil(t, err)
	report, err := sdk.GetCCReport("", "", nil)
	assert.Nil(t, err)
	report.Dump(cctrusted_base.QuoteDumpFormatHuman)

}

func TestSDKFullEventLog(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)

	sdk, err := GetSDKInstance(nil)
	assert.Nil(t, err)

	el, err := sdk.GetCCEventLog(0, 0)
	assert.Nil(t, err)
	el.Dump(cctrusted_base.QuoteDumpFormatHuman)

}
