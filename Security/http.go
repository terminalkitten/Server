package Security

import (
	"net/http"
	"crypto/sha512"
	"os"
	"encoding/base64"
	"io"
	"io/ioutil"
	"encoding/pem"
	"log"
	"crypto/x509"
	"crypto/sha256"
	"github.com/brokenbydefault/Server/Config"
)

func SetHeaders(fn func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-XSS-Protection", "1; mode=block")
		w.Header().Add("X-Content-Type-Options", "nosniff")
		w.Header().Add("X-Frame-Options", "DENY")

		w.Header().Add("Referrer-Policy", "same-origin")
		w.Header().Add("Content-Security-Policy", "default-src 'self'; font-src https://fonts.gstatic.com; connect-src wss://*.nanollet.org; upgrade-insecure-requests; block-all-mixed-content; disown-opener; require-sri-for script style;")

		w.Header().Add("Public-Key-Pins", "pin-sha256=\""+ CreateKeyPinning("/../.pki/cert.crt") +"\"; pin-sha256=\""+ CreateKeyBackupPinning("/../.pki/pubkeybackup.pem") +"\"; max-age=15768000")
		w.Header().Add("Expect-CT", "1")
		w.Header().Add("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		fn(w, r)
	}
}

func CreateSRI(filepath string) string {
	hash := sha512.New384()
	file, err := os.Open(Config.Dir() + filepath)
	if err != nil {
		panic(err)
	}
	io.Copy(hash, file)
	return "sha384-" + base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func CreateKeyPinning(file string) string {
	certpem, err := ioutil.ReadFile(Config.Dir() + file)
	if err != nil {
		log.Panic("You don't have one SSL certificate!", err)
	}

	block, _ := pem.Decode(certpem)
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Panic(err)
	}

	hash := sha256.Sum256(cert.RawSubjectPublicKeyInfo)
	return base64.StdEncoding.EncodeToString(hash[:])
}

func CreateKeyBackupPinning(file string) string {
	certpem, err := ioutil.ReadFile(Config.Dir() + file)
	if err != nil {
		log.Panic("You don't have one backup key, set your public-key backup into ../pki", err)
	}

	block, _ := pem.Decode(certpem)
	hash := sha256.Sum256(block.Bytes)
	return base64.StdEncoding.EncodeToString(hash[:])
}
