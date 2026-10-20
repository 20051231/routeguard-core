package mobile

import "github.com/hiddify/hiddify-core/v2/hcore"

// In-process only. The Android caller checks its physical network epoch before
// and after each call; no plaintext loopback mutation endpoint is registered.
func RouteGuardPoolSnapshot(epoch int64) (string, error) { return hcore.RouteGuardPoolSnapshot(epoch) }
func RouteGuardWarmPool(epoch int64) (string, error)     { return hcore.RouteGuardWarmPool(epoch) }
func RouteGuardSelectKeep(epoch int64, tag string) error {
	return hcore.RouteGuardSelectKeep(epoch, tag)
}
