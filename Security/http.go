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

var SSL_CERT_PATH = Config.Dir("..", ".pki", "cert.crt")
var SSL_KEY_PATH = Config.Dir( "..", ".pki", "key.pem")
var SSL_BACKUP_PUBLICKEY_PATH = Config.Dir( "..", ".pki", "pubkeybackup.pem")

func SetHeaders(fn func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-XSS-Protection", "1; mode=block")
		w.Header().Add("X-Content-Type-Options", "nosniff")
		w.Header().Add("X-Frame-Options", "DENY")

		w.Header().Add("Referrer-Policy", "same-origin")
		w.Header().Add("Content-Security-Policy", "default-src 'self'; font-src https://fonts.gstatic.com; connect-src wss://*.nanollet.org; upgrade-insecure-requests; block-all-mixed-content; disown-opener; require-sri-for script style;")
		w.Header().Add("Feature-Policy", "accelerometer 'none'; ambient-light-sensor 'none'; autoplay 'none'; camera 'none'; encrypted-media 'none'; fullscreen 'none'; geolocation 'none'; gyroscope 'none'; magnetometer 'none'; microphone 'none'; midi 'none'; payment 'none'; picture-in-picture 'none'; speaker 'none'; usb 'none'; vr 'none';")

		w.Header().Add("Public-Key-Pins", "pin-sha256=\""+ CreateKeyPinning(SSL_CERT_PATH) +"\"; pin-sha256=\""+ CreateKeyBackupPinning(SSL_BACKUP_PUBLICKEY_PATH) +"\"; max-age=15768000")
		w.Header().Add("Expect-CT", "enforce, max-age=31536000")
		w.Header().Add("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")

		fn(w, r)
	}
}

func CreateSRI(filepath string) string {
	hash := sha512.New384()
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	io.Copy(hash, file)
	return "sha384-" + base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

var KeyPinningCache = make(map[string]string)

func CreateKeyPinning(file string) string {
	if hash, ok := KeyPinningCache[file]; ok {
		return hash
	}

	certpem, err := ioutil.ReadFile(file)
	if err != nil {
		log.Panic("You don't have one SSL certificate!", err)
	}

	block, _ := pem.Decode(certpem)
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Panic(err)
	}

	hash := sha256.Sum256(cert.RawSubjectPublicKeyInfo)
	KeyPinningCache[file] = base64.StdEncoding.EncodeToString(hash[:])

	return KeyPinningCache[file]
}

func CreateKeyBackupPinning(file string) string {
	if hash, ok := KeyPinningCache[file]; ok {
		return hash
	}

	certpem, err := ioutil.ReadFile(file)
	if err != nil {
		log.Panic("You don't have one backup key, set your public-key backup into ", file, ":", err)
	}

	block, _ := pem.Decode(certpem)
	hash := sha256.Sum256(block.Bytes)
	KeyPinningCache[file] = base64.StdEncoding.EncodeToString(hash[:])

	return KeyPinningCache[file]
}
