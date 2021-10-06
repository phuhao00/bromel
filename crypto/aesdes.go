package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
	"encoding/hex"
	"errors"
)

/*
加密模式："ECB", "CBC", "CTR", "OFB", "CFB"
填充方式："PKCS5", "PKCS7", "ZERO"
输出格式："BASE64", "HEX"
*/

type EncryptWay uint

const (
	EncryptWayAES EncryptWay = iota + 1
	EncryptWatDES
)

type EncryptMode uint

const (
	AESDESModeECB EncryptMode = iota + 1
)

const (
	defaultAESDESECBBlockSize = 8
)

type PaddingWay uint

const (
	PaddingWayCS7 PaddingWay = iota + 1
	PaddingWayCS5
	PaddingWayZero
)

type EncodingWay uint

const (
	base64StdEncoding EncodingWay = iota + 1
	hexEncoding
)

type Encrypt struct {
	encryptWay  EncryptWay
	mode        EncryptMode
	paddingWay  PaddingWay
	encodingWay EncodingWay
}

//defaultAesEncrypt ...
func defaultAesEncrypt() *Encrypt {
	return &Encrypt{
		encryptWay:  EncryptWayAES,
		mode:        AESDESModeECB,
		paddingWay:  PaddingWayCS7,
		encodingWay: base64StdEncoding}
}

//NewEncrypt ...
func NewEncrypt(opts ...NewEncryptOption) *Encrypt {
	aes := defaultAesEncrypt()
	for _, opt := range opts {
		opt(aes)
	}
	return aes
}

//NewEncryptOption ...
type NewEncryptOption func(encrypt *Encrypt)

//WithEncryptWay ...
func WithEncryptWay(way EncryptWay) NewEncryptOption {
	return func(encrypt *Encrypt) {
		encrypt.encryptWay = way
	}
}

//WithAesEncryptMode  ...
func WithAesEncryptMode(mode EncryptMode) NewEncryptOption {
	return func(encrypt *Encrypt) {
		encrypt.mode = mode
	}
}

//WithAesEncryptModePaddingWay ...
func WithAesEncryptModePaddingWay(way PaddingWay) NewEncryptOption {
	return func(encrypt *Encrypt) {
		encrypt.paddingWay = way
	}
}

//WithEncodingWay ...
func WithEncodingWay(way EncodingWay) NewEncryptOption {
	return func(encrypt *Encrypt) {
		encrypt.encodingWay = way
	}
}

//Encrypt ...
func (e Encrypt) Encrypt(orig string, key string) (string, error) {
	origData := []byte(orig)
	k := []byte(key)
	block, err := e.GetCipherBlock(k)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	origData, err = e.Padding(origData, blockSize)
	if err != nil {
		return "", err
	}
	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])
	encrypted := make([]byte, len(origData))
	blockMode.CryptBlocks(encrypted, origData)
	return e.EnCodding(encrypted), nil
}

//Decrypt ...
func (e Encrypt) Decrypt(encrypted string, key string) (string, error) {
	encryptedByte, err := e.UnEnCodding(encrypted)
	if err != nil {
		return "", err
	}
	k := []byte(key)
	block, err := e.GetCipherBlock(k)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
	original := make([]byte, len(encryptedByte))
	blockMode.CryptBlocks(original, encryptedByte)
	original, err = e.UnPadding(original)
	if err != nil {
		return "", err
	}
	return string(original), nil
}

//GetCipherBlock ...
func (e Encrypt) GetCipherBlock(encrypted []byte) (cipher.Block, error) {
	if e.encryptWay == EncryptWayAES {
		return aes.NewCipher(encrypted)
	}
	if e.encryptWay == EncryptWatDES {
		return des.NewCipher(encrypted)
	}
	return nil, errors.New("unknown encrypt way ")
}

//EnCodding ...
func (e Encrypt) EnCodding(encrypted []byte) string {
	switch e.encodingWay {
	case base64StdEncoding:
		return base64.StdEncoding.EncodeToString(encrypted)
	case hexEncoding:
		return hex.EncodeToString(encrypted)
	}
	return ""
}

//UnEnCodding ...
func (e *Encrypt) UnEnCodding(encrypted string) ([]byte, error) {
	switch e.encodingWay {
	case base64StdEncoding:
		return base64.StdEncoding.DecodeString(encrypted)
	case hexEncoding:
		return hex.DecodeString(encrypted)
	}
	return nil, errors.New("unknown encoding way ")
}

//Padding ...
func (e Encrypt) Padding(ciphertext []byte, blockSize int) ([]byte, error) {
	if e.paddingWay == PaddingWayCS7 {
		return CS7Padding(ciphertext, blockSize), nil
	}
	if e.paddingWay == PaddingWayCS5 {
		return CS7Padding(ciphertext, defaultAESDESECBBlockSize), nil
	}
	if e.paddingWay == PaddingWayZero {
		return ZeroPadding(ciphertext, blockSize)
	}
	return nil, errors.New("unknown padding way ")
}

//UnPadding ...
func (e Encrypt) UnPadding(origData []byte) ([]byte, error) {
	switch e.paddingWay {
	case PaddingWayCS5, PaddingWayCS7:
		return CS7UnPadding(origData)
	case PaddingWayZero:
		return ZeroUnPadding(origData)
	}
	return nil, errors.New("unknown padding way ")
}

//ZeroPadding ...
func ZeroPadding(ciphertext []byte, blockSize int) ([]byte, error) {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padText...), nil
}

//ZeroUnPadding ...
func ZeroUnPadding(origData []byte) ([]byte, error) {
	return bytes.TrimFunc(origData,
		func(r rune) bool {
			return r == rune(0)
		}), nil
}

//CS7Padding 补码
//AES加密数据块分组长度必须为128bit(byte[16])，密钥长度可以是128bit(byte[16])、192bit(byte[24])、256bit(byte[32])中的任意一个。
func CS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

//CS7UnPadding 去码
func CS7UnPadding(origData []byte) ([]byte, error) {
	length := len(origData)
	unPadded := int(origData[length-1])
	return origData[:(length - unPadded)], nil
}
