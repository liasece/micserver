/**
 * \file GBEncode.go
 * \version
 * \author wzy
 * \date  2018年01月31日 11:27:16
 * \brief
 *
 */

package util

import (
	"bytes"
	"compress/zlib"
	"github.com/liasece/micserver/log"
	"io"
)

//进行zlib压缩
func ZlibCompress(src []byte) []byte {
	analysictime := FunctionTimeAnalysic{}
	analysictime.Start()
	defer analysictime.Stop()

	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	n, err := w.Write(src)
	w.Close()
	if err != nil {
		log.Error("[ZlibCompress] w.Write Error[%s] N[%d]",
			err.Error(), n)
	} else {
		log.Debug("[ZlibCompress] 压缩数据正常,%d-->>%d", len(src),
			len(in.Bytes()))
	}
	return in.Bytes()
}

//进行zlib解压缩
func ZlibUnCompress(compressSrc []byte) []byte {
	analysictime := FunctionTimeAnalysic{}
	analysictime.Start()
	defer analysictime.Stop()

	b := bytes.NewReader(compressSrc)
	out := new(bytes.Buffer)
	r, err := zlib.NewReader(b)
	if err != nil {
		log.Error("[zlib压缩]解压数据异常,%d,%s", len(compressSrc), err.Error())
		return []byte("")
	}
	n, err := io.Copy(out, r)

	if err != nil {
		log.Error("[ZlibUnCompress] io.Copy Error[%s] N[%d]",
			err.Error(), n)
	} else {
		log.Debug("[ZlibUnCompress] 解压数据正常,%d-->>%d", len(compressSrc),
			len(out.Bytes()))
	}
	return out.Bytes()
}
