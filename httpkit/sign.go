package httpkit

import (
	"bytes"
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"code.olapie.com/sugar/v2/base62"
	"code.olapie.com/sugar/v2/conv"
	"code.olapie.com/sugar/v2/ctxutil"
	"code.olapie.com/sugar/v2/maths"
	"code.olapie.com/sugar/v2/xerror"
)

const (
	KeyTimestamp = "X-Timestamp"
	KeySignature = "X-Signature"
)

type Signer interface {
	Sign(ctx context.Context, header http.Header) error
}

type SignerFunc func(ctx context.Context, header http.Header) error

func (h SignerFunc) Sign(ctx context.Context, header http.Header) error {
	return h(ctx, header)
}

type PrivateKey interface {
	ecdsa.PrivateKey | rsa.PrivateKey
}

func Sign[K PrivateKey](ctx context.Context, header http.Header, priv *K) error {
	SetHeaderNX(header, keyAppID, ctxutil.GetAppID(ctx))
	SetHeaderNX(header, keyClientID, ctxutil.GetClientID(ctx))
	traceID := ctxutil.GetTraceID(ctx)
	if traceID == "" {
		traceID = base62.NewUUIDString()
	}
	SetHeaderNX(header, keyTraceID, traceID)
	SetHeaderNX(header, KeyTimestamp, fmt.Sprint(time.Now().Unix()))

	hash := getMessageHashForSigning(header)
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

	header.Set(KeySignature, base64.StdEncoding.EncodeToString(sign))
	return nil
}

func GetSigner[K PrivateKey](priv *K) Signer {
	return SignerFunc(func(ctx context.Context, header http.Header) error {
		return Sign(ctx, header, priv)
	})
}

type Verifier interface {
	Verify(ctx context.Context, header http.Header) bool
}

type VerifierFunc func(ctx context.Context, header http.Header) bool

func (h VerifierFunc) Verify(ctx context.Context, header http.Header) bool {
	return h(ctx, header)
}

type PublicKey interface {
	ecdsa.PublicKey | rsa.PublicKey
}

func Verify[K PublicKey](ctx context.Context, header http.Header, pub *K) bool {
	ts := header.Get(KeyTimestamp)
	if ts == "" {
		fmt.Printf("[sugar/v2/httpkit] missing %s in header\n", KeyTimestamp)
		return false
	}
	t, err := strconv.ParseInt(ts, 0, 64)
	if err != nil {
		fmt.Printf("[sugar/v2/httpkit] strconv.ParseInt: %s, %v\n", ts, err)
		return false
	}

	if time.Now().Unix()-t > 5 {
		fmt.Printf("[sugar/v2/httpkit] outdated timestamp %s in header\n", ts)
		return false
	}

	if GetTraceID(header) == "" {
		fmt.Printf("[sugar/v2/httpkit] missing %s in header\n", keyTraceID)
		return false
	}

	signature := GetHeader(header, KeySignature)
	sign, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		fmt.Printf("[sugar/v2/httpkit] base64.DecodeString: %s, %v\n", signature, err)
		return false
	}
	hash := getMessageHashForSigning(header)
	switch k := any(pub).(type) {
	case *ecdsa.PublicKey:
		verified := ecdsa.VerifyASN1(k, hash, sign)
		if !verified {
			fmt.Printf("[sugar/v2/httpkit] failed verifying: %s\n", signature)
		}
		return verified
	case *rsa.PublicKey:
		err = rsa.VerifyPKCS1v15(k, crypto.SHA256, hash, sign)
		if err != nil {
			fmt.Printf("[sugar/v2/httpkit] failed verifying: %s, %v\n", signature, err)
			return false
		}
		return true
	default:
		fmt.Printf("[sugar/v2/httpkit] unsupported pub key: %T", pub)
		return false
	}
}

func GetVerifier[K PublicKey](pub *K) Verifier {
	return VerifierFunc(func(ctx context.Context, header http.Header) bool {
		return Verify(ctx, header, pub)
	})
}

func getMessageHashForSigning(h http.Header) []byte {
	var buf bytes.Buffer
	buf.WriteString(GetHeader(h, keyTraceID))
	buf.WriteString(GetHeader(h, KeyTimestamp))
	hash := md5.Sum(buf.Bytes())
	return hash[:]
}

func CheckTimestamp[H HeaderTypes](h H) error {
	ts := GetHeader(h, KeyTimestamp)
	if ts == "" {
		return xerror.New(http.StatusBadRequest, "missing %s", KeyTimestamp)
	}
	t, err := conv.ToInt64(ts)
	if err != nil {
		return xerror.New(http.StatusBadRequest, "invalid timestamp")
	}
	now := time.Now().Unix()
	if maths.Abs(now-t) > 60 {
		return xerror.New(http.StatusNotAcceptable, "outdated request")
	}
	return nil
}

func DecodeSign[H HeaderTypes](h H) ([]byte, error) {
	signature := GetHeader(h, KeySignature)
	if signature == "" {
		return nil, xerror.New(http.StatusBadRequest, "missing %s", KeySignature)
	}
	sign, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return nil, xerror.New(http.StatusNotAcceptable, "malformed signature")
	}
	return sign, nil
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
