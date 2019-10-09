package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/mitchellh/mapstructure"
)

func TestDecode(t *testing.T) {
	var cloudTrailEvent CloudTrailEventRecord

	testfileContents, err := ioutil.ReadFile("testdata/test.json")
	if err != nil {
		t.Errorf("Error: %s", err)
		return
	}

	err = json.Unmarshal(testfileContents, &cloudTrailEvent)
	if err != nil {
		t.Errorf("Error: %s", err)
		return
	}

	var securityGroupEvent CloudTrailSecurityGroupEvent

	if err = mapstructure.Decode(cloudTrailEvent.RequestParameters, &securityGroupEvent); err != nil {
		t.Errorf("Error: %s", err)
		return
	}

	if securityGroupEvent.GroupID != "sg-05ffcaf1d3252d12d" {
		t.Errorf("Error: %s", err)
		return
	}

}
