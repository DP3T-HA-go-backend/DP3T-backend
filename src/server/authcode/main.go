package main

import (
	"dp3t-backend/store"

	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"github.com/spf13/viper"
)

type AuthCode struct {
	Code string `json:"code"`
}

var code AuthCode
var data store.Store

func gencode(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("x-public-key", "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUZrd0V3WUhLb1pJemowQ0FRWUlLb1pJemowREFRY0RRZ0FFTWl5SEU4M1lmRERMeWg5R3dCTGZsYWZQZ3pnNgpJanhySjg1ejRGWjlZV3krU2JpUDQrWW8rL096UFhlbDhEK0o5TWFrMXpvT2FJOG4zRm90clVnM2V3PT0KLS0tLS1FTkQgUFVCTElDIEtFWS0tLS0t")

	mySigningKey, err0 := ioutil.ReadFile("/ec256-key")
	if err0 != nil {
		fmt.Println("Unable to load ECDSA private key: ", err0)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	block, _ := pem.Decode(mySigningKey)
	if block == nil || block.Type != "EC PRIVATE KEY" {
		fmt.Println("failed to decode PEM block containing EC privatekey")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ecdsaKey, err2 := x509.ParseECPrivateKey(block.Bytes)
	if err2 != nil {
		fmt.Println("Unable to parse ECDSA private key: ", err2)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	var generatedcode string

	// Make sure we generate a new key
	for {
		generatedcode = ""
		for i := 0; i < 12; i++ {
			generatedcode = generatedcode + strconv.Itoa(r1.Intn(10))
		}

		if data.AddAuthCode(generatedcode) == nil {
			break
		}
	}

	code = AuthCode{Code: generatedcode}

	m, err := json.Marshal(code)

	if err != nil {
		log.Fatal("Failed to encode AuthCode: ", err)
	}

	h := sha256.Sum256([]byte(m))
	digest := base64.StdEncoding.EncodeToString(h[:])
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"content-hash": digest,
		"hash-alg":     "sha256",
		"iss":          "d3pt",
		"iat":          strconv.FormatInt(makeTimestampSeconds(), 10),
		"exp":          strconv.FormatInt(makeTimestampSeconds()+1814400, 10),
	})

	tokenString, err := token.SignedString(ecdsaKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Digest", "sha-256="+digest)
	w.Header().Set("Signature", tokenString)
	w.WriteHeader(http.StatusOK)
	fmt.Println("GET:", r.URL)

	w.Write(m)
}

func makeTimestampMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func makeTimestampSeconds() int64 {
	return time.Now().UnixNano() / int64(time.Second)
}

func main() {
	viper.SetDefault("core.port", 8080)

	viper.SetConfigName("config")
	viper.SetConfigType("ini")
	viper.AddConfigPath("/service/etc/authcode/")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Failed to read config file: ", err)
	}

	router := httprouter.New()
	router.GET("/", gencode)

	addr := fmt.Sprint(":", viper.GetInt("core.port"))
	fmt.Println("Listening on", addr)

	log.Fatal(http.ListenAndServe(addr, router))
}
