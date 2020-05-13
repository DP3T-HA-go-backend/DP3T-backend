package main

import (
	"dp3t-backend/api"
	"dp3t-backend/store"
	"dp3t-backend/server"

	"crypto/sha256"
	"flag"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/dgrijalva/jwt-go.v3"
)

var conf *server.Config
var data store.Store

func gencode(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	var generatedcode string

	// Make sure we generate a new key
	for {
		generatedcode = ""
		for i := 0; i < 12; i++ {
			generatedcode = generatedcode + strconv.Itoa(r1.Intn(10))
		}

		if err := data.AddAuthCode(generatedcode); err == nil {
			break
		}
	}

	code := api.ProtoAuthData{Value: generatedcode}

	m, err := json.Marshal(code)
	if err != nil {
		log.Println("ERROR: Encoding JSON:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ts_s := time.Now().UnixNano() / int64(time.Second)
	time := strconv.FormatInt(ts_s, 10)
	time_exp := strconv.FormatInt(ts_s + 1814400, 10)

	h := sha256.Sum256([]byte(m))
	digest := base64.StdEncoding.EncodeToString(h[:])
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"content-hash":       digest,
		"hash-alg":           "sha256",
		"iss":                "d3pt",
		"iat":                time,
		"exp":                time_exp,
	})

	signature, err := token.SignedString(conf.PrivateKey)
	if err != nil {
		log.Println("ERROR: Token signature:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println("INFO: GET", r.URL)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Digest", "sha-256=" + digest)
	w.Header().Set("Signature", signature)
	w.Header().Set("x-public-key", server.PUBLIC_KEY)
	w.WriteHeader(http.StatusOK)
	w.Write(m)
}

func main() {
	conf_file_p := flag.String("config", "./config/authcode.ini", "path to config file")
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

	data.Init(conf)

	router := httprouter.New()
	router.GET("/", gencode)

	addr := fmt.Sprint(":", conf.Port)

	log.Println("INFO: Config file:", *conf_file_p)
	log.Println("INFO: Key file:", conf.PrivateKeyFile)
	log.Println("INFO: Store:", conf.StoreType)
	log.Println("INFO: Listening on:", addr)

	log.Fatal(http.ListenAndServe(addr, router))
}
