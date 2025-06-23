//go:build !(arm || mips || s390x)

package kdf

import (
	"testing"

	"github.com/yunmoon/gmsm/sm3"
)

// This case should be failed on 32bits system.
func TestKdfPanic(t *testing.T) {
	shouldPanic(t, func() {
		Kdf(sm3.New, []byte("123456"), 1<<37)
	})
}
