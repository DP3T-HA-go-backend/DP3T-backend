package main

import (
	"crypto/sha256"
	"crypto/x509"
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
	Port    int    `ini:"port"`
	KeyFile string `ini:"keyfile"`
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

	mySigningKey, err0 := ioutil.ReadFile(conf.KeyFile)
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

	m, err := proto.Marshal(&data)
	if err != nil {
		log.Fatal("Failed to encode ProtoExposedList: ", err)
	}

	h := sha256.Sum256([]byte(m))
	digest := base64.StdEncoding.EncodeToString(h[:])
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"content-hash":       digest,
		"hash-alg":           "sha256",
		"iss":                "d3pt",
		"iat":                strconv.FormatInt(makeTimestampSeconds(), 10),
		"exp":                strconv.FormatInt(makeTimestampSeconds()+1814400, 10),
		"batch-release-time": strconv.FormatInt(makeTimestampMillis(), 10),
	})

	tokenString, err := token.SignedString(ecdsaKey)
	//fmt.Println(tokenString, err)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Digest", "sha-256="+digest)
	w.Header().Set("x-protobuf-message", "org.dpppt.backend.sdk.model.proto.ProtoExposedList")
	w.Header().Set("x-protobuf-schema", "exposed.proto")
	w.Header().Set("Signature", tokenString)
	w.WriteHeader(http.StatusOK)
	fmt.Println("GET:", r.URL)

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

func hello(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Hello\n")
}

func makeTimestampMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func makeTimestampSeconds() int64 {
	return time.Now().UnixNano() / int64(time.Second)
}

func main() {
	config_p := flag.String("config", "./config.ini", "path to config file")
	flag.Parse()

	i, err := ini.Load(*config_p)
	if err != nil {
		log.Fatal("Failed to read config file: ", err)
	}

	if err := i.MapTo(&conf); err != nil {
		log.Fatal("Failed to decode config: ", err)
	}

	// Initialize exposed data
	data = ProtoExposedList{
		BatchReleaseTime: 123456789,
		Exposed:          []*ProtoExposee{},
	}

	router := httprouter.New()
	router.GET("/", hello)
	router.GET("/exposed/:date", exposed)
	router.POST("/exposed", expose)

	addr := fmt.Sprint(":", conf.Port)

	fmt.Println("Key file:", conf.KeyFile)
	fmt.Println("Listening on:", addr)

	log.Fatal(http.ListenAndServe(addr, router))
}
