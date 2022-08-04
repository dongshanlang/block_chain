/**
 * @Author: lixiumin
 * @E-Mail: lixiuminmxl@163.com
 * @Date: 2022/7/26 5:55 PM
 * @Desc:
 */

package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"
)

func uintToByte(num uint64) []byte {
	//use binary.Writer encode
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.BigEndian, num)
	if err != nil {
		log.Panicln(err)
	}
	return buffer.Bytes()
}

func IsFileExist(fileName string) bool {
	_, err := os.Stat(fileName)
	if err == nil {
		return true
	}
	return false
}
