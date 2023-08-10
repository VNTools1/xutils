/*
 * @Date: 2023-07-20 09:38:29
 * @LastEditTime: 2023-07-20 09:44:24
 * @Description:
 */
package xhashring

import (
	"fmt"
	"testing"
)

func TestInitHashRing(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "device1"}, {name: "device2"}, {name: "device3"}, {name: "device4"}, {name: "device5"}, {name: "device6机加工"},
	}
	hashRing := NewHashring("device", 10)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, _ := hashRing.GetNode(tt.name)
			fmt.Println(node)
			t.Fatal("failed")
		})
	}
}
