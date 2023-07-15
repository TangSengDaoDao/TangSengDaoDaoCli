package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstall(t *testing.T) {
	ctx := NewTangSengDaoDaoContext(NewTangSengDaoDao())
	install := newInstallCMD(ctx)

	err := install.run(nil, nil)
	assert.NoError(t, err)
}
