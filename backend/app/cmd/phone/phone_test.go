package phone

import (
	"testing"
)

func TestSendSms(t *testing.T) {
	if err := SendSms("  8 950-592 68-09   ", "5777"); err != nil {
		t.Error(err)
		return
	}
}
