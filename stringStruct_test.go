/**
 * @Author: lixiumin
 * @E-Mail: lixiuminmxl@163.com
 * @Date: 2022/8/9 11:51 AM
 * @Desc:
 */

package main

import (
	"fmt"
	"testing"
)

type TestS struct {
	str string
}

func (s TestS) String() string {
	return "hello"
}
func TestHain(t *testing.T) {
	t1 := TestS{"hello"}
	fmt.Printf("%v\n", t1)
}
