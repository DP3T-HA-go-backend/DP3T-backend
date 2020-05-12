package main

import (
	"os"
	"crypto/sha256"
	"crypto/x509"
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/pem"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"gopkg.in/ini.v1"
	"gopkg.in/dgrijalva/jwt-go.v3"
	"github.com/julienschmidt/httprouter"
)

type Config struct {
	Port           int    `ini:"port"`
	PrivateKeyFile string `ini:"private-key-file"`
	PrivateKey     *ecdsa.PrivateKey
}

var conf Config

var data ProtoExposedList

func exposed(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO: Validate date & retrieve data based on it from key-value store
	// var date string
	// date = ps.ByName("date")

	w.Header().Set("Content-Type", "application/x-protobuf")
	w.Header().Set("x-public-key", "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUZrd0V3WUhLb1pJemowQ0FRWUlLb1pJemowREFRY0RRZ0FFTWl5SEU4M1lmRERMeWg5R3dCTGZsYWZQZ3pnNgpJanhySjg1ejRGWjlZV3krU2JpUDQrWW8rL096UFhlbDhEK0o5TWFrMXpvT2FJOG4zRm90clVnM2V3PT0KLS0tLS1FTkQgUFVCTElDIEtFWS0tLS0t")
	w.Header().Set("x-batch-release-time", strconv.FormatInt(makeTimestampMillis(), 10))

	m, err := proto.Marshal(&data)
	if err != nil {
		fmt.Println("ERROR:", "protobuf encoding:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h := sha256.Sum256([]byte(m))
	digest := base64.StdEncoding.EncodeToString(h[:])
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"content-hash":       digest,
		"hash-alg":           "sha256",
		"iss":                "d3pt",
		"iat":                strconv.FormatInt(makeTimestampSeconds(), 10),
		"exp":                strconv.FormatInt(makeTimestampSeconds() + 1814400, 10),
		"batch-release-time": strconv.FormatInt(makeTimestampMillis(), 10),
	})

	signature, err := token.SignedString(conf.PrivateKey)
	if err != nil {
		fmt.Println("ERROR:", "token signature:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println("GET:", r.URL)

	w.Header().Set("Digest", "sha-256="+digest)
	w.Header().Set("x-protobuf-message", "org.dpppt.backend.sdk.model.proto.ProtoExposedList")
	w.Header().Set("x-protobuf-schema", "exposed.proto")
	w.Header().Set("Signature", signature)
	w.WriteHeader(http.StatusOK)
	w.Write(m)
}

func expose(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	in, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024))
	if err != nil {
		log.Fatal("Error reading request:", err)
	}

	fmt.Println("POST:", r.URL, string(in))

	exposee := &ProtoExposee{}
	if err := protojson.Unmarshal(in, exposee); err != nil {
		log.Fatal("Failed to parse Exposee: ", err)
	}

	// TODO: Add data to key-value store
	data.Exposed = append(data.Exposed, exposee)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "OK\n")
}

func makeTimestampMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func makeTimestampSeconds() int64 {
	return time.Now().UnixNano() / int64(time.Second)
}

func initConfig(conf_file string) error {
	i, e := ini.Load(conf_file)
	if e != nil {
		return fmt.Errorf("Failed to read config file: %s", e)
	}

	if e := i.MapTo(&conf); e != nil {
		return fmt.Errorf("Failed to decode config: %s", e)
	}

	if _, e := os.Stat(conf.PrivateKeyFile); e != nil {
		return fmt.Errorf("Failed to read private key: %s", e)
	}

	keyfile, e := ioutil.ReadFile(conf.PrivateKeyFile)
	if e != nil {
		return fmt.Errorf("Failed to read private key: %s", e)
	}

	block, _ := pem.Decode(keyfile)
	if block == nil || block.Type != "EC PRIVATE KEY" {
		return fmt.Errorf("Failed to decode PEM block containing EC private key")
	}

	conf.PrivateKey, e = x509.ParseECPrivateKey(block.Bytes)
	if e != nil {
		return fmt.Errorf("Failed to parse EC private key: %s", e)
	}

	return nil
}

func main() {
	conf_file_p := flag.String("config", "./config.ini", "path to config file")
	flag.Parse()

	err := initConfig(*conf_file_p)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize exposed data
	data = ProtoExposedList{
		BatchReleaseTime: 123456789,
		Exposed:          []*ProtoExposee{},
	}

	router := httprouter.New()
	router.GET("/:date", exposed)
	router.POST("/", expose)

	addr := fmt.Sprint(":", conf.Port)
	fmt.Println("Key file:", conf.PrivateKeyFile)
	fmt.Println("Listening on:", addr)

	log.Fatal(http.ListenAndServe(addr, router))
}
