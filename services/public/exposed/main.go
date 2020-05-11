package main

import (
    "fmt"
    "log"
    "io"
    "io/ioutil"
    "net/http"
    "strconv"
    "time"
    //"crypto/ecdsa"
    "crypto/x509"
    "encoding/pem"
//    "crypto/elliptic"
//    "crypto/rand"
    //"crypto/sha256"

    "github.com/dgrijalva/jwt-go"
    "github.com/julienschmidt/httprouter"
    "google.golang.org/protobuf/proto"
    "google.golang.org/protobuf/encoding/protojson"
)

var data ProtoExposedList


func exposed(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    // TODO: Validate date & retrieve data based on it from key-value store
    // var date string
    // date = ps.ByName("date")

    //type Token struct {
    //  Raw       string                 // The raw token.  Populated when you Parse a token
    //  Method    SigningMethod          // The signing method used or to be used
    //  Header    map[string]interface{} // The first segment of the token
    //  Claims    Claims                 // The second segment of the token
    //  Signature string                 // The third segment of the token.  Populated when you Parse a token
    //  Valid     bool                   // Is the token valid?  Populated when you Parse/Verify a token
    //0
    // https://godoc.org/github.com/dgrijalva/jwt-go#example-New--Hmac}

    w.Header().Set("Content-Type", "application/x-protobuf")
    w.Header().Set("x-public-key", "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUZrd0V3WUhLb1pJemowQ0FRWUlLb1pJemowREFRY0RRZ0FFTWl5SEU4M1lmRERMeWg5R3dCTGZsYWZQZ3pnNgpJanhySjg1ejRGWjlZV3krU2JpUDQrWW8rL096UFhlbDhEK0o5TWFrMXpvT2FJOG4zRm90clVnM2V3PT0KLS0tLS1FTkQgUFVCTElDIEtFWS0tLS0t")
    w.Header().Set("x-batch-release-time",strconv.FormatInt(makeTimestampMillis(),10))

    //mySigningKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
    //mySigningKey := []byte("MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQgevZzL1gdAFr88hb2OF/2NxApJCzGCEDdfSp6VQO30hyhRANCAAQRWz+jn65BtOMvdyHKcvjBeBSDZH2r1RTwjmYSi9R/zpBnuQ4EiMnCqfMPWiZqB4QdbAd0E7oH50VpuZ1P087G")
    //fmt.Println(mySigningKey, err)
    mySigningKey, err0 := ioutil.ReadFile("/ec256-key")
    if err0 != nil {
		fmt.Println("Unable to load ECDSA private key: ", err0)
		w.WriteHeader(http.StatusInternalServerError)
		return
    }
    //fmt.Println("key: ",mySigningKey)
    block, _ := pem.Decode(mySigningKey)
    if block == nil || block.Type != "EC PRIVATE KEY" {
		fmt.Println("failed to decode PEM block containing EC privatekey")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

    //fmt.Println("key: ",mySigningKey)
    //var ecdsaKey *ecdsa.PrivateKey
    //ecdsaKey, err2 := x509.ParsePKCS8PrivateKey(block.Bytes)
    ecdsaKey, err2 := x509.ParseECPrivateKey(block.Bytes)
    //
    if err2 != nil {
		fmt.Println("Unable to parse ECDSA private key: ", err2)
		w.WriteHeader(http.StatusInternalServerError)
		return
    }

    token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims {
		"content-hash": "digesthereSHA256",
		"hash-alg": "sha256",
		"iss": "d3pt",
		"iat": strconv.FormatInt(makeTimestampSeconds(),10),
		"exp": strconv.FormatInt(makeTimestampSeconds()+1814400,10),
		"batch-release-time": strconv.FormatInt(makeTimestampMillis(),10),
             })
    tokenString, err := token.SignedString(ecdsaKey)
    fmt.Println(tokenString, err)
    if err != nil {
	// If there is an error in creating the JWT return an internal server error
	w.WriteHeader(http.StatusInternalServerError)
	return
    }

    w.Header().Set("Signature", tokenString)
    w.WriteHeader(http.StatusOK)


    m, err := proto.Marshal(&data)
    if err != nil {
        log.Fatal("Failed to encode ProtoExposedList: ", err)
    }

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
    data = ProtoExposedList{
        BatchReleaseTime: 123456789,
        Exposed: []*ProtoExposee{},
    }

    router := httprouter.New()
    router.GET("/", hello)
    router.GET("/exposed/:date", exposed)
    router.POST("/exposed", expose)

    log.Fatal(http.ListenAndServe(":8080", router))
}
