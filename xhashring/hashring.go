/*
 * @Date: 2023-07-20 09:38:29
 * @LastEditTime: 2023-07-20 09:43:16
 * @Description:
 */
package xhashring

import (
	"fmt"

	"github.com/serialx/hashring"
)

func NewHashring(prefix string, slots int) *hashring.HashRing {
	if slots == 0 {
		slots = 20
	}
	if prefix == "" {
		prefix = "hashring_test"
	}
	hashKeys := []string{}
	for i := 0; i < slots; i++ {
		key := fmt.Sprintf("%s_%d", prefix, i)
		hashKeys = append(hashKeys, key)
	}

	ring := hashring.New(hashKeys)
	return ring
}
