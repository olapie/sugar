package mobilex

import "code.olapie.com/sugar/cryptox"

func Encrypt(data []byte, passphrase string) []byte {
	content, _ := cryptox.Encrypt(data, passphrase)
	return content
}

func Decrypt(data []byte, passphrase string) []byte {
	content, _ := cryptox.Decrypt(data, passphrase)
	return content
}

func EncryptFile(dst, src, passphrase string) bool {
	err := cryptox.EncryptFile(cryptox.Destination(dst), cryptox.Source(src), passphrase)
	return err == nil
}

func DecryptFile(dst, src, passphrase string) bool {
	err := cryptox.DecryptFile(cryptox.Destination(dst), cryptox.Source(src), passphrase)
	return err == nil
}
