package mailer

import "testing"

func TestConnectClient(t *testing.T) {
	if err := ConnectClient("127.0.0.1", 27111); err != nil {
		t.Error(err)
	}
}
