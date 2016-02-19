package jwt

import (
	"github.com/RangelReale/osin"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-errors/errors"
	"github.com/pborman/uuid"
	"io"
	"io/ioutil"
	"os"
	"regexp"
)

var TestCertificates = [][]string{
	{"../example/cert/rs256-private.pem",
		`-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA4f5wg5l2hKsTeNem/V41fGnJm6gOdrj8ym3rFkEU/wT8RDtn
SgFEZOQpHEgQ7JL38xUfU0Y3g6aYw9QT0hJ7mCpz9Er5qLaMXJwZxzHzAahlfA0i
cqabvJOMvQtzD6uQv6wPEyZtDTWiQi9AXwBpHssPnpYGIn20ZZuNlX2BrClciHhC
PUIIZOQn/MmqTD31jSyjoQoV7MhhMTATKJx2XrHhR+1DcKJzQBSTAGnpYVaqpsAR
ap+nwRipr3nUTuxyGohBTSmjJ2usSeQXHI3bODIRe1AuTyHceAbewn8b462yEWKA
Rdpd9AjQW5SIVPfdsz5B6GlYQ5LdYKtznTuy7wIDAQABAoIBAQCwia1k7+2oZ2d3
n6agCAbqIE1QXfCmh41ZqJHbOY3oRQG3X1wpcGH4Gk+O+zDVTV2JszdcOt7E5dAy
MaomETAhRxB7hlIOnEN7WKm+dGNrKRvV0wDU5ReFMRHg31/Lnu8c+5BvGjZX+ky9
POIhFFYJqwCRlopGSUIxmVj5rSgtzk3iWOQXr+ah1bjEXvlxDOWkHN6YfpV5ThdE
KdBIPGEVqa63r9n2h+qazKrtiRqJqGnOrHzOECYbRFYhexsNFz7YT02xdfSHn7gM
IvabDDP/Qp0PjE1jdouiMaFHYnLBbgvlnZW9yuVf/rpXTUq/njxIXMmvmEyyvSDn
FcFikB8pAoGBAPF77hK4m3/rdGT7X8a/gwvZ2R121aBcdPwEaUhvj/36dx596zvY
mEOjrWfZhF083/nYWE2kVquj2wjs+otCLfifEEgXcVPTnEOPO9Zg3uNSL0nNQghj
FuD3iGLTUBCtM66oTe0jLSslHe8gLGEQqyMzHOzYxNqibxcOZIe8Qt0NAoGBAO+U
I5+XWjWEgDmvyC3TrOSf/KCGjtu0TSv30ipv27bDLMrpvPmD/5lpptTFwcxvVhCs
2b+chCjlghFSWFbBULBrfci2FtliClOVMYrlNBdUSJhf3aYSG2Doe6Bgt1n2CpNn
/iu37Y3NfemZBJA7hNl4dYe+f+uzM87cdQ214+jrAoGAXA0XxX8ll2+ToOLJsaNT
OvNB9h9Uc5qK5X5w+7G7O998BN2PC/MWp8H+2fVqpXgNENpNXttkRm1hk1dych86
EunfdPuqsX+as44oCyJGFHVBnWpm33eWQw9YqANRI+pCJzP08I5WK3osnPiwshd+
hR54yjgfYhBFNI7B95PmEQkCgYBzFSz7h1+s34Ycr8SvxsOBWxymG5zaCsUbPsL0
4aCgLScCHb9J+E86aVbbVFdglYa5Id7DPTL61ixhl7WZjujspeXZGSbmq0Kcnckb
mDgqkLECiOJW2NHP/j0McAkDLL4tysF8TLDO8gvuvzNC+WQ6drO2ThrypLVZQ+ry
eBIPmwKBgEZxhqa0gVvHQG/7Od69KWj4eJP28kq13RhKay8JOoN0vPmspXJo1HY3
CKuHRG+AP579dncdUnOMvfXOtkdM4vk0+hWASBQzM9xzVcztCa+koAugjVaLS9A+
9uQoqEeVNTckxx0S2bYevRy7hGQmUJTyQm3j1zEUR5jpdbL83Fbq
-----END RSA PRIVATE KEY-----
`},
	{"../example/cert/rs256-public.pem",
		`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA4f5wg5l2hKsTeNem/V41
fGnJm6gOdrj8ym3rFkEU/wT8RDtnSgFEZOQpHEgQ7JL38xUfU0Y3g6aYw9QT0hJ7
mCpz9Er5qLaMXJwZxzHzAahlfA0icqabvJOMvQtzD6uQv6wPEyZtDTWiQi9AXwBp
HssPnpYGIn20ZZuNlX2BrClciHhCPUIIZOQn/MmqTD31jSyjoQoV7MhhMTATKJx2
XrHhR+1DcKJzQBSTAGnpYVaqpsARap+nwRipr3nUTuxyGohBTSmjJ2usSeQXHI3b
ODIRe1AuTyHceAbewn8b462yEWKARdpd9AjQW5SIVPfdsz5B6GlYQ5LdYKtznTuy
7wIDAQAB
-----END PUBLIC KEY-----
`,
	},
}

type JWT struct {
	privateKey []byte
	publicKey  []byte
}

func New(privateKey, publicKey []byte) *JWT {
	return &JWT{
		privateKey: privateKey,
		publicKey:  publicKey,
	}
}

// Helper func: Read certificate from specified file
func LoadCertificate(path string) ([]byte, error) {
	if path == "" {
		return nil, errors.Errorf("No path specified")
	} else if _, err := os.Stat(path); os.IsNotExist(err) {
		return []byte(path), nil
	}

	var rdr io.Reader
	if f, err := os.Open(path); err == nil {
		rdr = f
		defer f.Close()
	} else {
		return nil, err
	}
	return ioutil.ReadAll(rdr)
}

// Verify a token and output the claims.
func (j *JWT) VerifyToken(tokenData []byte) (*jwt.Token, error) {
	// trim possible whitespace from token
	tokenData = regexp.MustCompile(`\s*$`).ReplaceAll(tokenData, []byte{})

	// Parse the token.  Load the key from command line option
	token, err := jwt.Parse(string(tokenData), func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		return jwt.ParseRSAPublicKeyFromPEM(j.publicKey)
	})
	if err != nil {
		return nil, errors.Errorf("Couldn't parse token: %v", err)
	} else if !token.Valid {
		return nil, errors.Errorf("Token is invalid")
	}

	claims := ClaimsCarrier(token.Claims)
	if claims.AssertExpired() {
		token.Valid = false
		return token, errors.Errorf("Token expired at %v", claims.GetExpiresAt())
	}
	//if claims.AssertNotYetValid() {
	//	token.Valid = false
	//	return token, errors.Errorf("Token validates in the future: %v", claims.GetNotBefore())
	//}
	return token, nil
}

// Create, sign, and return a token.
func (j *JWT) SignToken(claims map[string]interface{}, header map[string]interface{}) (string, error) {
	if _, ok := header["alg"]; ok {
		return "", errors.New("You may not override the alg header key.")
	}

	if _, ok := header["typ"]; ok {
		return "", errors.New("You may not override the typ header key.")
	}

	token := jwt.New(jwt.SigningMethodRS256)
	token.Claims = claims
	token.Header = merge(token.Header, header)
	ecdsaKey, err := jwt.ParseRSAPrivateKeyFromPEM(j.privateKey)
	if err != nil {
		return "", err
	}
	return token.SignedString(ecdsaKey)
}

func merge(a, b map[string]interface{}) map[string]interface{} {
	for k, w := range b {
		if _, ok := a[k]; ok {
			continue
		}
		a[k] = w
	}
	return a
}

func (j *JWT) GenerateAccessToken(data *osin.AccessData, generateRefresh bool) (accessToken string, refreshToken string, err error) {
	claims, ok := data.UserData.(ClaimsCarrier)
	if !ok {
		return "", "", errors.Errorf("Could not assert claims to ClaimsCarrier: %v", claims)
	}

	claims["exp"] = data.ExpireAt().Unix()
	if accessToken, err = j.SignToken(claims, map[string]interface{}{}); err != nil {
		return "", "", err
	} else if !generateRefresh {
		return
	}

	claims = ClaimsCarrier{}
	claims["id"] = uuid.New()
	claims["aud"] = claims.GetAudience()
	if refreshToken, err = j.SignToken(claims, map[string]interface{}{}); err != nil {
		return "", "", err
	}
	return
}
