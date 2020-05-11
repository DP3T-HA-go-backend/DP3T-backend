package main

import (
    "fmt"
    "log"
    "io"
    "io/ioutil"
    "net/http"
    "strconv"
    "time"
    "crypto/sha256"
    "crypto/x509"
    "encoding/pem"
    "encoding/base64"
//    "crypto/elliptic"
//    "crypto/rand"
    //"crypto/sha256"

    "github.com/dgrijalva/jwt-go"
    "github.com/julienschmidt/httprouter"
    "google.golang.org/protobuf/proto"
    "google.golang.org/protobuf/encoding/protojson"
    "github.com/spf13/viper"
)

var data ProtoExposedList

func exposed(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO: Validate date & retrieve data based on it from key-value store
	// var date string
	// date = ps.ByName("date")

	w.Header().Set("Content-Type", "application/x-protobuf")
	w.Header().Set("x-public-key", "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUZrd0V3WUhLb1pJemowQ0FRWUlLb1pJemowREFRY0RRZ0FFTWl5SEU4M1lmRERMeWg5R3dCTGZsYWZQZ3pnNgpJanhySjg1ejRGWjlZV3krU2JpUDQrWW8rL096UFhlbDhEK0o5TWFrMXpvT2FJOG4zRm90clVnM2V3PT0KLS0tLS1FTkQgUFVCTElDIEtFWS0tLS0t")
	w.Header().Set("x-batch-release-time", strconv.FormatInt(makeTimestampMillis(), 10))

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

    m, err := proto.Marshal(&data)
    if err != nil {
        log.Fatal("Failed to encode ProtoExposedList: ", err)
    }

    h := sha256.Sum256([]byte(m))
    digest := base64.StdEncoding.EncodeToString(h[:])
    token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims {
		"content-hash": digest,
		"hash-alg": "sha256",
		"iss": "d3pt",
		"iat": strconv.FormatInt(makeTimestampSeconds(),10),
		"exp": strconv.FormatInt(makeTimestampSeconds()+1814400,10),
		"batch-release-time": strconv.FormatInt(makeTimestampMillis(),10),
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
	viper.SetDefault("core.port", 8080)

	viper.SetConfigName("config")
	viper.SetConfigType("ini")
	viper.AddConfigPath("/service/etc/exposed/")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Failed to read config file: ", err)
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

	addr := fmt.Sprint(":", viper.GetInt("core.port"))
	fmt.Println("Listening on", addr)

	log.Fatal(http.ListenAndServe(addr, router))
}
