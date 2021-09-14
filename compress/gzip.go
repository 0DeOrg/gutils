package compress

/**
 * @Author: lee
 * @Description:
 * @File: gzip
 * @Date: 2021/9/13 4:58 下午
 */

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

//func GZipCompress(text string) []byte {
//
//
//}

func GZipUnCompress(msg []byte) string {
	reader := bytes.NewReader(msg)
	r, _ := gzip.NewReader(reader)

	defer r.Close()
	undatas, _ := ioutil.ReadAll(r)
	return string(undatas)
}