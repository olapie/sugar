package xhttp

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
	"code.olapie.com/sugar/v2/xcontext"
	"code.olapie.com/sugar/v2/xerror"
	"code.olapie.com/sugar/v2/xmath"
)

const (
	KeyTimestamp = "X-Timestamp"
	KeySignature = "X-Signature"
)

type Signer interface {
	Sign(ctx context.Context, req *http.Request) error
}

type SignerFunc func(ctx context.Context, req *http.Request) error

func (h SignerFunc) Sign(ctx context.Context, req *http.Request) error {
	return h(ctx, req)
}

type PrivateKey interface {
	ecdsa.PrivateKey | rsa.PrivateKey
}

func Sign[K PrivateKey](ctx context.Context, req *http.Request, priv *K) error {
	if xcontext.HasLogin(ctx) {
		if login := xcontext.GetLogin[string](ctx); login != "" {
			SetHeaderNX(req.Header, KeyUserID, login)
		} else if login := xcontext.GetLogin[int64](ctx); login != 0 {
			SetHeaderNX(req.Header, KeyUserID, fmt.Sprint(login))
		}
	}
	SetHeaderNX(req.Header, KeyAppID, xcontext.GetAppID(ctx))
	SetHeaderNX(req.Header, KeyClientID, xcontext.GetClientID(ctx))
	traceID := xcontext.GetTraceID(ctx)
	if traceID == "" {
		traceID = base62.NewUUIDString()
	}
	SetHeaderNX(req.Header, KeyTraceID, traceID)
	SetHeaderNX(req.Header, KeyTimestamp, fmt.Sprint(time.Now().Unix()))

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
	return SignerFunc(func(ctx context.Context, req *http.Request) error {
		return Sign(ctx, req, priv)
	})
}

type Verifier interface {
	Verify(ctx context.Context, req *http.Request) bool
}

type VerifierFunc func(ctx context.Context, req *http.Request) bool

func (h VerifierFunc) Verify(ctx context.Context, req *http.Request) bool {
	return h(ctx, req)
}

type PublicKey interface {
	ecdsa.PublicKey | rsa.PublicKey
}

func Verify[K PublicKey](ctx context.Context, req *http.Request, pub *K) bool {
	ts := req.Header.Get(KeyTimestamp)
	if ts == "" {
		fmt.Printf("[sugar/v2/xhttp] missing %s in header\n", KeyTimestamp)
	}
	t, err := strconv.ParseInt(ts, 0, 64)
	if err != nil {
		fmt.Printf("[sugar/v2/xhttp] strconv.ParseInt: %s, %v\n", ts, err)
		return false
	}

	if time.Now().Unix()-t > 5 {
		fmt.Printf("[sugar/v2/xhttp] outdated timestamp %s in header\n", ts)
		return false
	}

	signature := req.Header.Get(KeySignature)
	sign, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		fmt.Printf("[sugar/v2/xhttp] base64.DecodeString: %s, %v\n", signature, err)
		return false
	}
	hash := getMessageHashForSigning(req)
	switch k := any(pub).(type) {
	case *ecdsa.PublicKey:
		verified := ecdsa.VerifyASN1(k, hash, sign)
		if !verified {
			fmt.Printf("[sugar/v2/xhttp] failed verifying: %s\n", signature)
		}
		return verified
	case *rsa.PublicKey:
		err = rsa.VerifyPKCS1v15(k, crypto.SHA256, hash, sign)
		if err != nil {
			fmt.Printf("[sugar/v2/xhttp] failed verifying: %s, %v\n", signature, err)
			return false
		}
		return true
	default:
		fmt.Printf("[sugar/v2/xhttp] unsupported pub key: %T", pub)
		return false
	}
}

func GetVerifier[K PublicKey](pub *K) Verifier {
	return VerifierFunc(func(ctx context.Context, req *http.Request) bool {
		return Verify(ctx, req, pub)
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
	buf.WriteString(GetHeader(req.Header, KeyTraceID))
	buf.WriteString(GetHeader(req.Header, KeyTimestamp))
	hash := md5.Sum(buf.Bytes())
	return hash[:]
}

func CheckTimestamp[H Headerxtypeet](h H) error {
	ts := GetHeader(h, KeyTimestamp)
	if ts == "" {
		return xerror.BadRequest("missing %s", KeyTimestamp)
	}
	t, err := conv.ToInt64(ts)
	if err != nil {
		return xerror.BadRequest("invalid timestamp")
	}
	now := time.Now().Unix()
	if xmath.Abs(now-t) > 60 {
		return xerror.NotAcceptable("outdated request")
	}
	return nil
}

func DecodeSign[H Headerxtypeet](h H) ([]byte, error) {
	signature := GetHeader(h, KeySignature)
	if signature == "" {
		return nil, xerror.BadRequest("missing %s", KeySignature)
	}
	sign, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return nil, xerror.NotAcceptable("malformed signature")
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
