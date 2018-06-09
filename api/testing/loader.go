package testing

import (
	"encoding/json"
	"io/ioutil"

	. "github.com/onsi/gomega"
)

//File load a file verifying possible errors
func File(path string) []byte {
	blob, err := ioutil.ReadFile(path)
	Expect(err).Should(BeNil())
	return blob
}

//LoadJSON loads a json file
func LoadJSON(path string, target interface{}) {
	LoadJSONFromBytes(File(path), target)
}

//LoadJSONFromBytes loads a json verifying possible errors
func LoadJSONFromBytes(blob []byte, target interface{}) {
	err := json.Unmarshal(blob, &target)
	Expect(err).Should(BeNil())
}

//JSONToBytes marshal a json verifying possible errors
func JSONToBytes(data interface{}) []byte {
	blob, err := json.Marshal(data)
	Expect(err).Should(BeNil())
	return blob
}
