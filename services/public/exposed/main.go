package main

import (
    "fmt"
    "log"
    "io"
    "io/ioutil"
    "net/http"

    "github.com/julienschmidt/httprouter"
    "google.golang.org/protobuf/proto"
)

var data ProtoExposedList

func exposed(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    // TODO: Validate date & select data based on it
    // var date string
    // date = ps.ByName("date")

    w.Header().Set("Content-Type", "application/x-protobuf")
    w.WriteHeader(http.StatusCreated)

    m, err := proto.Marshal(&data)
    if err != nil {
        log.Fatal("Failed to encode ProtoExposedList: ", err)
    }

    w.Write(m)
}

func expose(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    exposee := &ProtoExposee{}

    in, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024))
    if err != nil {
        log.Fatal("Error reading request:", err)
    }

    if err := proto.Unmarshal(in, exposee); err != nil {
        log.Fatal("Failed to parse Exposee: ", err)
    }

    data.Exposed = append(data.Exposed, exposee)

    w.WriteHeader(http.StatusOK)
    fmt.Fprint(w, "OK\n")
}

func hello(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    w.WriteHeader(http.StatusOK)
    fmt.Fprint(w, "hello\n")
}

func main() {
    data = ProtoExposedList{
        BatchReleaseTime: 123456789,
        Exposed: []*ProtoExposee{},
    }

    router := httprouter.New()
    router.GET("/v1", hello)
    router.GET("/v1/exposed/:date", exposed)
    router.POST("/v1/exposed", expose)

    log.Fatal(http.ListenAndServe(":8080", router))
}
