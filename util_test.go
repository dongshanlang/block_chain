/**
 * @Author: lixiumin
 * @E-Mail: lixiuminmxl@163.com
 * @Date: 2022/8/4 5:16 PM
 * @Desc:
 */

package main

import "testing"

func TestIsFileExist(t *testing.T) {
	exist := IsFileExist(BlockChainDBName)
	t.Log(exist)
}
