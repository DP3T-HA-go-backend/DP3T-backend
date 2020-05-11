package main

import (
	"fmt"

	"./lib"
)


func main() {
	kvs.KVPutTTL("foo", "boo", 14)
	kvs.KVPutTTL("fooo", "booo", 14)
	kvs.KVPutTTL("fooooo", "se fueeee", 14)

	//kvs.KVDelete("fooo")

	/*resp := kvs.KVGet("foo")
	for _, ev := range resp.Kvs {
                fmt.Printf("%s : %s\n", ev.Key, ev.Value)
        }*/

	resp := kvs.KVGetWithPrefix("fo")
	for _, ev := range resp.Kvs {
                fmt.Printf("%s : %s\n", ev.Key, ev.Value)
        }

}
