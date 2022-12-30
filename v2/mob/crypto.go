package mob

import "code.olapie.com/sugar/olasec"

func Encrypt(data []byte, passphrase string) []byte {
	content, _ := olasec.Encrypt(data, passphrase)
	return content
}

func Decrypt(data []byte, passphrase string) []byte {
	content, _ := olasec.Encrypt(data, passphrase)
	return content
}

func EncryptFile(src, dst, passphrase string) bool {
	err := olasec.EncryptFile(olasec.SF(src), olasec.DF(dst), passphrase)
	return err == nil
}

func DecryptFile(src, dst, passphrase string) bool {
	err := olasec.DecryptFile(olasec.SF(src), olasec.DF(dst), passphrase)
	return err == nil
}
