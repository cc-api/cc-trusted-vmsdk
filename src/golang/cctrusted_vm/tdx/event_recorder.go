package tdx

import (
	"errors"
	"fmt"
	"os"

	cctrusted_vm "github.com/cc-api/cc-trusted-vmsdk/src/golang/cctrusted_vm"
)

const (
	DEFAULT_ACPI_TABLE_FILE      = "/sys/firmware/acpi/tables/CCEL"
	DEFAULT_ACPI_TABLE_DATA_FILE = "/sys/firmware/acpi/tables/data/CCEL"
)

var _ cctrusted_vm.EventRecorder = (*TDXEventLogRecorder)(nil)

type TDXEventLogRecorder struct {
	acpiTableFile     string
	acpiTableDataFile string
	rawEventLog       []byte
	redirected        bool
	// redirectedAcpiTableFile is the alternative
	// of the original `acpiTableFile`, if which
	// can not be accessed
	redirectedAcpiTableFile string
	// redirectedAcpiTableDataFile is the alternative
	// of the original `acpiTableDataFile`, if which
	// can not be accessed
	redirectedAcpiTableDataFile string
}

// FullEventLog implements cctrusted_vm.EventRecorder.
func (t *TDXEventLogRecorder) FullEventLog() ([]byte, error) {
	path := t.acpiTableDataFile
	if t.redirected {
		path = t.redirectedAcpiTableDataFile
	}
	log, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	t.rawEventLog = log
	return t.rawEventLog, nil
}

// ProbeRecorder implements cctrusted_vm.EventRecorder.
func (t *TDXEventLogRecorder) ProbeRecorder() error {
	t.acpiTableFile = DEFAULT_ACPI_TABLE_FILE
	t.acpiTableDataFile = DEFAULT_ACPI_TABLE_DATA_FILE
	if _, err := os.Stat(t.acpiTableFile); err != nil {
		if _, e := os.Stat(t.redirectedAcpiTableFile); e != nil {
			return fmt.Errorf("event log file open file: %v & %v", err, e)
		}
		t.redirected = true
	}
	if _, err := os.Stat(t.acpiTableDataFile); err != nil {
		if _, e := os.Stat(t.redirectedAcpiTableDataFile); e != nil {
			return fmt.Errorf("event log file open file: %v & %v", err, e)
		}
	}

	path := t.acpiTableFile
	if t.redirected {
		path = t.redirectedAcpiTableFile
	}

	if err := t.isValidCCELTable(path); err != nil {
		return err
	}

	return nil
}

func (t *TDXEventLogRecorder) isValidCCELTable(file string) error {
	d, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	if len(d) < 4 || string(d[:4]) != "CCEL" {
		return errors.New("invalid CCEL table")
	}
	return nil
}

func (t *TDXEventLogRecorder) RedirectAcpiTableFile(file string) {
	t.redirectedAcpiTableFile = file
}

func (t *TDXEventLogRecorder) RedirectAcpiTableDataFile(file string) {
	t.redirectedAcpiTableDataFile = file
}
