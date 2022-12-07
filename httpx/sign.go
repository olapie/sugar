package httpx

import (
	"bytes"
	"code.olapie.com/sugar/conv"
	"code.olapie.com/sugar/errorx"
	"code.olapie.com/sugar/mathx"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const (
	KeyTimestamp = "X-Timestamp"
	KeySignature = "X-Signature"
)

type Signer interface {
	Sign(req *http.Request) error
}

type SignerFunc func(req *http.Request) error

func (h SignerFunc) Sign(req *http.Request) error {
	return h(req)
}

type PrivateKey interface {
	ecdsa.PrivateKey | rsa.PrivateKey
}

func Sign[K PrivateKey](req *http.Request, priv *K) error {
	ts := fmt.Sprint(time.Now().Unix())
	req.Header.Set(KeyTimestamp, ts)
	hash := getMessageHashForSigning(req)
	var sign []byte
	var err error
	switch k := any(priv).(type) {
	case *ecdsa.PrivateKey:
		sign, err = ecdsa.SignASN1(rand.Reader, k, hash)
	case *rsa.PrivateKey:
		sign, err = rsa.SignPKCS1v15(rand.Reader, k, crypto.SHA256, hash)
	default:
		err = fmt.Errorf("invalid private key: %T", k)
	}

	if err != nil {
		return err
	}

	req.Header.Set(KeySignature, base64.StdEncoding.EncodeToString(sign))
	return nil
}

func GetSigner[K PrivateKey](priv *K) Signer {
	return SignerFunc(func(req *http.Request) error {
		return Sign(req, priv)
	})
}

type Verifier interface {
	Verify(req *http.Request) bool
}

type VerifierFunc func(req *http.Request) bool

func (h VerifierFunc) Verify(req *http.Request) bool {
	return h(req)
}

type PublicKey interface {
	ecdsa.PublicKey | rsa.PublicKey
}

func Verify[K PublicKey](req *http.Request, pub *K) bool {
	ts := req.Header.Get(KeyTimestamp)
	if ts == "" {
		fmt.Printf("[sugar/httpx] missing %s in header\n", KeyTimestamp)
	}
	t, err := strconv.ParseInt(ts, 0, 64)
	if err != nil {
		fmt.Printf("[sugar/httpx] strconv.ParseInt: %s, %v\n", ts, err)
		return false
	}

	if time.Now().Unix()-t > 5 {
		return false
	}

	signature := req.Header.Get(KeySignature)
	sign, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		fmt.Printf("[sugar/httpx] base64.DecodeString: %s, %v\n", signature, err)
		return false
	}
	hash := getMessageHashForSigning(req)
	switch k := any(pub).(type) {
	case *ecdsa.PublicKey:
		return ecdsa.VerifyASN1(k, hash, sign)
	case *rsa.PublicKey:
		return rsa.VerifyPKCS1v15(k, crypto.SHA256, hash, sign) == nil
	default:
		return false
	}
}

func GetVerifier[K PublicKey](pub *K) Verifier {
	return VerifierFunc(func(req *http.Request) bool {
		return Verify(req, pub)
	})
}

func getMessageHashForSigning(req *http.Request) []byte {
	var buf bytes.Buffer
	buf.WriteString(req.Method)
	path := req.URL.Path
	if path == "" {
		path = req.URL.RawPath
	}
	buf.WriteString(path)
	buf.WriteString(req.URL.RawQuery)
	buf.WriteString(GetHeader(req.Header, KeyContentType))
	buf.WriteString(GetHeader(req.Header, KeyAppID))
	buf.WriteString(GetHeader(req.Header, KeyClientID))
	buf.WriteString(GetHeader(req.Header, KeyTimestamp))
	buf.WriteString(GetHeader(req.Header, KeyAuthorization))
	hash := sha256.Sum256(buf.Bytes())
	return hash[:]
}

func CheckTimestamp[H HeaderTypeSet](h H) error {
	ts := GetHeader(h, KeyTimestamp)
	if ts == "" {
		return errorx.BadRequest("missing %s", KeyTimestamp)
	}
	t, err := conv.ToInt64(ts)
	if err != nil {
		return errorx.BadRequest("invalid timestamp")
	}
	now := time.Now().Unix()
	if mathx.Abs(now-t) > 10 {
		return errorx.NotAcceptable("outdated request")
	}
	return nil
}

func example() {
	// generate and marshal
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pubData, _ := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	priData, _ := x509.MarshalECPrivateKey(privateKey)

	// parse
	pri, _ := x509.ParseECPrivateKey(priData)
	anyPub, _ := x509.ParsePKIXPublicKey(pubData)
	pub := anyPub.(*ecdsa.PublicKey)
	_ = pri
	_ = pub
}
