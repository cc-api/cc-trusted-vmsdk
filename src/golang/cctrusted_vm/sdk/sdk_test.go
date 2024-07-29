package vmsdk

import (
	"bytes"
	"log"
	"testing"

	"github.com/cc-api/evidence-api/common/golang/evidence_api"

	"github.com/stretchr/testify/assert"
)

func TestSDKReport(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)

	sdk, err := GetSDKInstance(nil)
	assert.Nil(t, err)
	report, err := sdk.GetCCReport("", "", nil)
	assert.Nil(t, err)
	report.Dump(evidence_api.QuoteDumpFormatHuman)

}

func TestSDKFullEventLog(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)

	sdk, err := GetSDKInstance(nil)
	assert.Nil(t, err)

	el, err := sdk.GetCCEventLog(0, 0)
	assert.Nil(t, err)
	el.Dump(evidence_api.QuoteDumpFormatHuman)

}
