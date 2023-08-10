package xdraw

import (
	"testing"

	"github.com/VNTools1/xutils/ximg"
	"github.com/stretchr/testify/assert"
)

func TestResize(t *testing.T) {
	src := ximg.Read("../../logo.png")
	dst := Resize(src, 100, 100)
	assert.NotNil(t, dst)
}
