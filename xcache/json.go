// Package xcache
/*
 * @Date: 2023-07-20 09:46:36
 * @LastEditTime: 2023-07-20 09:47:05
 * @Description:
 */
package xcache

import (
	"encoding/json"
)

type DefaultJSONSerializer struct{}

// Serialize
// @param v
// @date 2022-07-02 08:12:26
func (d *DefaultJSONSerializer) Serialize(v any) ([]byte, error) {
	return json.Marshal(v)
}

// Deserialize
// @param data
// @param v
// @date 2022-07-02 08:12:25
func (d *DefaultJSONSerializer) Deserialize(data []byte, v any) error {

	return json.Unmarshal(data, v)
}
