package crypto

import (
	"encoding/hex"
	"fmt"
	"testing"
)

//指定了key长度大小
var (
	aesTestKey = "X7WBOELqgn6dc8CN"
	desTestKey = "gn6dc8CN"
)

func TestAesEncryptAndDeEncrypt(t *testing.T) {
	aes := NewEncrypt(
		WithEncryptWay(EncryptWatDES),
		WithAesEncryptMode(AESDESModeECB),
		WithAesEncryptModePaddingWay(PaddingWayCS7),
		WithEncodingWay(hexEncoding))
	encrypted, err := aes.Encrypt("abc", desTestKey)
	if err != nil {
		panic(err)
	}
	original, err := aes.Decrypt(encrypted, desTestKey)
	if err != nil {
		panic(err)
	}
	fmt.Println(original)
}

func TestHex(t *testing.T) {
	tmp := []byte("huhao555555555506575hgfdhg府大院  有try啊【】【】【奋斗684324 那就好 gfdgfsfh3wrwer vxcxbgf    53......,,,bvbcnbncvg隔热人忒特 特瑞特人他5464554````.//,,m,m,m")
	i := fmt.Sprintf("%x", tmp)
	fmt.Println(i)
	i1, err := hex.DecodeString(i)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(i1))
}
