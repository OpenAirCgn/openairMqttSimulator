package openairMqttSimulator

import "testing"

const testfn = "test/client_list.json"

func TestLoadClientList(t *testing.T) {
	list, err := LoadClientCertConfigList(testfn)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 3 {
		t.Fatalf("did not load 3 entries: %d", len(list))
	}

	if "test/7722745105e9e02e8f1aaf17f7b3aac5c56cd805.crt" != list[2].CRTFilename {
		t.Fatalf("fn incorrect: %s", list[1].CRTFilename)
	}
}
