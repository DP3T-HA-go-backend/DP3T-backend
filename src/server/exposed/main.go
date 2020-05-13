package main

import (
	"dp3t-backend/api"
	"dp3t-backend/store"
	"dp3t-backend/server"

	"crypto/sha256"
	"encoding/base64"
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

	"github.com/julienschmidt/httprouter"
	"gopkg.in/dgrijalva/jwt-go.v3"
)

var conf *server.Config
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
	w.Header().Set("x-public-key", server.PUBLIC_KEY)
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

func main() {
	conf_file_p := flag.String("config", "./config/exposed.ini", "path to config file")
	flag.Parse()

	var err error
	conf, err = server.InitConfig(*conf_file_p)
	if err != nil {
		log.Fatal("ERROR: ", err)
	}

	data, err = store.InitStore(conf)
	if err != nil {
		log.Fatal("ERROR: ", err)
	}

	err = data.Init(conf)
	if err != nil {
		log.Fatal("ERROR: ", err)
	}

	router := httprouter.New()
	router.GET("/:date", exposed)
	router.POST("/", expose)

	addr := fmt.Sprint(":", conf.Port)

	log.Println("INFO: Config file:", *conf_file_p)
	log.Println("INFO: Key file:", conf.PrivateKeyFile)
	log.Println("INFO: Store:", conf.StoreType)
	log.Println("INFO: Listening on:", addr)

	log.Fatal(http.ListenAndServe(addr, router))
}
