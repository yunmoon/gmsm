package smx509

import (
	"sync"

	"github.com/yunmoon/gmsm/internal/godebug"
)

var (
	once           sync.Once
	systemRootsMu  sync.RWMutex
	systemRoots    *CertPool
	systemRootsErr error
	fallbacksSet   bool
)

func systemRootsPool() *CertPool {
	once.Do(initSystemRoots)
	systemRootsMu.RLock()
	defer systemRootsMu.RUnlock()
	return systemRoots
}

func initSystemRoots() {
	systemRootsMu.Lock()
	defer systemRootsMu.Unlock()
	systemRoots, systemRootsErr = loadSystemRoots()
	if systemRootsErr != nil {
		systemRoots = nil
	}
}

// SetFallbackRoots sets the roots to use during certificate verification, if no
// custom roots are specified and a platform verifier or a system certificate
// pool is not available (for instance in a container which does not have a root
// certificate bundle). SetFallbackRoots will panic if roots is nil.
//
// SetFallbackRoots may only be called once, if called multiple times it will
// panic.
//
// The fallback behavior can be forced on all platforms, even when there is a
// system certificate pool, by setting GODEBUG=x509usefallbackroots=1 (note that
// on Windows and macOS this will disable usage of the platform verification
// APIs and cause the pure Go verifier to be used). Setting
// x509usefallbackroots=1 without calling SetFallbackRoots has no effect.
func SetFallbackRoots(roots *CertPool) {
	if roots == nil {
		panic("roots must be non-nil")
	}

	// trigger initSystemRoots if it hasn't already been called before we
	// take the lock
	_ = systemRootsPool()

	systemRootsMu.Lock()
	defer systemRootsMu.Unlock()

	if fallbacksSet {
		panic("SetFallbackRoots has already been called")
	}
	fallbacksSet = true
	if systemRoots != nil && (systemRoots.len() > 0 || systemRoots.systemPool) && (godebug.Get("x509usefallbackroots") != "1") {
		return
	}
	systemRoots, systemRootsErr = roots, nil
}
