package app

import sdk "github.com/cc-api/cc-trusted-vmsdk/src/golang/cctrusted_vm/sdk"

var instance sdk.SDK
var cvmErr error

func GetSDK() (sdk.SDK, error) {
	return instance, cvmErr
}

func init() {
	instance, cvmErr = sdk.GetSDKInstance(nil)
}
