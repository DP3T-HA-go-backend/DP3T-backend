package main

import (
	"dp3t-backend/api"
	"dp3t-backend/store"

	"crypto/ecdsa"
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
	"os"
	"strconv"
	"time"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/dgrijalva/jwt-go.v3"
	"gopkg.in/ini.v1"
)

const PUBLIC_KEY string = "" +
	"LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUZrd0V3WUhLb1pJemowQ0FRWUlL" +
	"b1pJemowREFRY0RRZ0FFTWl5SEU4M1lmRERMeWg5R3dCTGZsYWZQZ3pnNgpJanhy" +
	"Sjg1ejRGWjlZV3krU2JpUDQrWW8rL096UFhlbDhEK0o5TWFrMXpvT2FJOG4zRm90" +
	"clVnM2V3PT0KLS0tLS1FTkQgUFVCTElDIEtFWS0tLS0t"

type Config struct {
	Port           int    `ini:"port"`
	PrivateKeyFile string `ini:"private-key-file"`
	PrivateKey     *ecdsa.PrivateKey
}

var conf Config

var data store.Store

func exposed(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO: Validate date & retrieve data based on it from key-value store
	// var date string
	// date = ps.ByName("date")

	// TODO: Need to pass an appropriate time value
	exposed, _ := data.GetExposed(1234)

	m, err := proto.Marshal(exposed)
	if err != nil {
		log.Println("ERROR: Encoding protobuf:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ts_s := time.Now().UnixNano() / int64(time.Second)
	ts_ms := time.Now().UnixNano() / int64(time.Millisecond)

	time := strconv.FormatInt(ts_s, 10)
	time_exp := strconv.FormatInt(ts_s + 1814400, 10)
	time_ms := strconv.FormatInt(ts_ms, 10)

	h := sha256.Sum256([]byte(m))
	digest := base64.StdEncoding.EncodeToString(h[:])
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"content-hash":       digest,
		"hash-alg":           "sha256",
		"iss":                "d3pt",
		"iat":                time,
		"exp":                time_exp,
		"batch-release-time": time_ms,
	})

	signature, err := token.SignedString(conf.PrivateKey)
	if err != nil {
		log.Println("ERROR: Token signature:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println("INFO: GET", r.URL)

	w.Header().Set("Content-Type", "application/x-protobuf")
	w.Header().Set("Digest", "sha-256=" + digest)
	w.Header().Set("Signature", signature)
	w.Header().Set("x-public-key", PUBLIC_KEY)
	w.Header().Set("x-batch-release-time", time_ms)
	w.Header().Set("x-protobuf-message", "org.dpppt.backend.sdk.model.proto.ProtoExposedList")
	w.Header().Set("x-protobuf-schema", "exposed.proto")
	w.WriteHeader(http.StatusOK)
	w.Write(m)
}

func expose(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	in, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024))
	if err != nil {
		log.Println("ERROR: Reading request:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	exposee := &api.ProtoExposee{}
	if err := protojson.Unmarshal(in, exposee); err != nil {
		log.Println("ERROR: Decoding JSON:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data.AddExposee(exposee)

	log.Println("INFO: POST", r.URL, string(in))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "OK\n")
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
	conf_file_p := flag.String("config", "./config/exposed.ini", "path to config file")
	flag.Parse()

	err := initConfig(*conf_file_p)
	if err != nil {
		log.Fatal("ERROR: ", err)
	}

	data = &store.InMem{}
	data.Init()

	router := httprouter.New()
	router.GET("/:date", exposed)
	router.POST("/", expose)

	addr := fmt.Sprint(":", conf.Port)

	log.Println("INFO: Config file:", *conf_file_p)
	log.Println("INFO: Key file:", conf.PrivateKeyFile)
	log.Println("INFO: Listening on:", addr)

	log.Fatal(http.ListenAndServe(addr, router))
}
