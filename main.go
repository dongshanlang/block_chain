/**
 * @Author: lixiumin
 * @E-Mail: lixiuminmxl@163.com
 * @Date: 2022/7/26 2:04 PM
 * @Desc:
 */

package main

var (
	gensisInfo = "this is the first block"
)

func main() {
	bc := NewBlockChain("monitor")
	cli := CLI{bc: bc}
	cli.Run()
}
