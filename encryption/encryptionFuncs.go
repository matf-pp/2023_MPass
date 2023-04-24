package encryption

type Encryption interface {
	DeriveKey(passphrase string, salt []byte) ([]byte, []byte)
	Encrypt(passphrase, plaintext string) string
	Decrypt(passphrase, ciphertext string) string
	CreateAuthKey(passphrase, cipthertext string) ([]byte, []byte)
	StoreAuthKey(pathname string, authKey []byte)
	RetreiveAuthKey(pathname string) []byte
	StoreEncryptedData(pathname, encryptedData string)
	RetreiveEncryptedData(pathname string) string
	ValidatePassword(passphrase, cipthertext string, authKey []byte) bool
	Init()
	AssertAvailablePRNG()
	GenerateRandomBytes(n int) ([]byte, error)
	GenerateRandomString(n int) (string, error)
	GenerateRandomStringURLSafe(n int) (string, error)
}
