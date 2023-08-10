package xdraw

import (
	"image"
	"testing"

	"github.com/VNTools1/xutils/ximg"
	"github.com/stretchr/testify/assert"
)

func TestWatermark(t *testing.T) {
	src := ximg.Read("../../logo.png")
	dst := WatermarkText(src, "2023-04-07", "#f00", 64, 100, 100)
	assert.NotNil(t, dst)

	assert.NotNil(t, WatermarkImage(src, dst, image.Point{X: 0, Y: 0}))
}
