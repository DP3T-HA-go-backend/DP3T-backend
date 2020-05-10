package main

import (
    "fmt"
    "log"
    "io"
    "io/ioutil"
    "net/http"
    "strconv"
    "time"

    "github.com/julienschmidt/httprouter"
    "google.golang.org/protobuf/proto"
    "google.golang.org/protobuf/encoding/protojson"
)

var data ProtoExposedList

func exposed(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    // TODO: Validate date & retrieve data based on it from key-value store
    // var date string
    // date = ps.ByName("date")

    w.Header().Set("Content-Type", "application/x-protobuf")
    w.Header().Set("x-public-key", "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUZrd0V3WUhLb1pJemowQ0FRWUlLb1pJemowREFRY0RRZ0FFdkxXZHVFWThqcnA4aWNSNEpVSlJaU0JkOFh2UgphR2FLeUg2VlFnTXV2Zk1JcmxrNk92QmtKeHdhbUdNRnFWYW9zOW11di9rWGhZdjF1a1p1R2RjREJBPT0KLS0tLS1FTkQgUFVCTElDIEtFWS0tLS0tCg")
    w.Header().Set("x-batch-release-time",strconv.FormatInt(makeTimestamp(),10))
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

func makeTimestamp() int64 {
    return time.Now().UnixNano() / int64(time.Millisecond)
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
