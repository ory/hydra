package config

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	r "gopkg.in/dancannon/gorethink.v2"
)

var cert1 = `-----BEGIN CERTIFICATE-----
MIIB/zCCAamgAwIBAgIJAPKYmr9KduB0MA0GCSqGSIb3DQEBCwUAMFsxCzAJBgNV
BAYTAlVTMRMwEQYDVQQIDApDYWxpZm9ybmlhMRQwEgYDVQQHDAtMb3MgQW5nZWxl
czEhMB8GA1UECgwYSW50ZXJuZXQgV2lkZ2l0cyBQdHkgTHRkMB4XDTE2MDYyMzEx
NTM0NloXDTI2MDYyMTExNTM0NlowWzELMAkGA1UEBhMCVVMxEzARBgNVBAgMCkNh
bGlmb3JuaWExFDASBgNVBAcMC0xvcyBBbmdlbGVzMSEwHwYDVQQKDBhJbnRlcm5l
dCBXaWRnaXRzIFB0eSBMdGQwXDANBgkqhkiG9w0BAQEFAANLADBIAkEA29Rsxzzh
ZkN6b1UZ8eQJcJBLsSxEDJPdHaztJbL1azYr+pdSCBx7QGL+8odxC6Vur9Y3keZl
fkUza9YoKTxUpwIDAQABo1AwTjAdBgNVHQ4EFgQUMlXwDvfVUxBftSsr3Hl9HMDo
tY0wHwYDVR0jBBgwFoAUMlXwDvfVUxBftSsr3Hl9HMDotY0wDAYDVR0TBAUwAwEB
/zANBgkqhkiG9w0BAQsFAANBAIPOzp8x4rS0sOFfUaXacKgFBm7wh+ski5P35chZ
0BZAXTka0jqHhLRo+qsD/mNGQWwUkDbtLSkMZVxxO5HMJXg=
-----END CERTIFICATE-----`

var cert2 = `-----BEGIN CERTIFICATE-----
MIIB/TCCAaegAwIBAgIJAPywvLC/ldr6MA0GCSqGSIb3DQEBCwUAMFoxCzAJBgNV
BAYTAlVTMRMwEQYDVQQIDApXYXNoaW5ndG9uMRMwEQYDVQQHDApXYXNoaW5ndG9u
MSEwHwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQwHhcNMTYwNjIzMTE1
NDI3WhcNMjYwNjIxMTE1NDI3WjBaMQswCQYDVQQGEwJVUzETMBEGA1UECAwKV2Fz
aGluZ3RvbjETMBEGA1UEBwwKV2FzaGluZ3RvbjEhMB8GA1UECgwYSW50ZXJuZXQg
V2lkZ2l0cyBQdHkgTHRkMFwwDQYJKoZIhvcNAQEBBQADSwAwSAJBAMMXumS376Qp
yK/Jaa83MOcyIGwrmLxuu0vKWQU5u1USSuFGGnIQ08YhBcklQI8+t4aJY+hCsfDx
5PKNIX8SxNsCAwEAAaNQME4wHQYDVR0OBBYEFF6FjQx8T3JOU46GAdGNQ7UWDQEj
MB8GA1UdIwQYMBaAFF6FjQx8T3JOU46GAdGNQ7UWDQEjMAwGA1UdEwQFMAMBAf8w
DQYJKoZIhvcNAQELBQADQQC4aXmUbyBiVaXWnkKCGzbf5Uxs6y3td+togd1La1mn
d9ahfgeVHNG0Dz9VJu2LA3aB4pAWPlbxd2m0frzK1Sx4
-----END CERTIFICATE-----`

var TestImportRethinkDBRootCAData = []struct {
	Name     string
	EnvKey   string
	EnvValue string
	Cert     string
}{
	{"Without Certificate", "", "", ""},
	{"Certificate string 1", "RETHINK_TLS_CERT", cert1, cert1},
	{"Certificate string 2", "RETHINK_TLS_CERT", cert2, cert2},
	{"Certificate path 1", "RETHINK_TLS_CERT_PATH", "cert1.pem", cert1},
	{"Certificate path 2", "RETHINK_TLS_CERT_PATH", "cert2.pem", cert2},
}

func TestImportRethinkDBRootCA(t *testing.T) {
	ioutil.WriteFile("cert1.pem", []byte(cert1), 0644)
	ioutil.WriteFile("cert2.pem", []byte(cert2), 0644)

	for _, test := range TestImportRethinkDBRootCAData {
		viper.Reset()
		if test.EnvKey != "" {
			viper.Set(test.EnvKey, test.EnvValue)
		}

		opts := r.ConnectOpts{}

		importRethinkDBRootCA(&opts)

		if test.Cert != "" {
			require.NotNil(t, opts.TLSConfig, test.Name)

			// We try to add the same certificate twice to see if it has added the correct one
			assert.Equal(t, 1, len(opts.TLSConfig.RootCAs.Subjects()), test.Name)
			opts.TLSConfig.RootCAs.AppendCertsFromPEM([]byte(test.Cert))
			assert.Equal(t, 1, len(opts.TLSConfig.RootCAs.Subjects()), test.Name)
		}
	}

	// Cleanup
	viper.Reset()
	os.Remove("cert1.pem")
	os.Remove("cert2.pem")
}
