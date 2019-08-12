//             ,%%%%%%%%,
//           ,%%/\%%%%/\%%
//          ,%%%\c "" J/%%%
// %.       %%%%/ o  o \%%%
// `%%.     %%%%    _  |%%%
//  `%%     `%%%%(__Y__)%%'
//  //       ;%%%%`\-/%%%'
// ((       /  `%%%%%%%'
//  \\    .'          |
//   \\  /       \  | |
//    \\/攻城狮保佑) | |
//     \         /_ | |__
//     (___________)))))))                   `\/'
/*
 * 修订记录:
 * long.qian 2019-08-12 17:29 创建
 */

/**
 * @author long.qian
 */

package compress

import (
	"bytes"
	"compress/zlib"
	"io"
)

//进行zlib压缩
func ZlibCompress(src []byte) ([]byte,error) {
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	_,err := w.Write(src)
	if err != nil {
		return nil,err
	}
	_ = w.Close()
	return in.Bytes(),nil
}

//进行zlib解压缩
func ZlibUnCompress(compressSrc []byte) ([]byte,error) {
	b := bytes.NewReader(compressSrc)
	var out bytes.Buffer
	r, _ := zlib.NewReader(b)
	_,err := io.Copy(&out, r)
	if err != nil {
		return nil,err
	}
	return out.Bytes(),nil
}