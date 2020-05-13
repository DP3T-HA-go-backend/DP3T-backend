package main

import (
	"fmt"

	//	"./lib"
	kvs "kvs/lib"
)

func main() {
	kvs.KVPutTTL("foo", "boo", 14)
	kvs.KVPutTTL("foo", "booo", 14)
	kvs.KVPutTTL("fooooo", "se fueeee", 14)
	//var resp string
	resp := kvs.KVGet("foo37")
	if resp.Count == 0 {
		fmt.Printf("didnt exist\n")
	} else {
		fmt.Printf("existed\n")
	}

	for _, ev := range resp.Kvs {
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	}
	//kvs.KVDelete("fooo")

	kvs.KVDelete("footoput")
	kvs.KVPutTTL("footodelete", "boo", 14)
	kvs.KVPutAndDelete("footodelete", "footoput", "valueput")
	resp = kvs.KVGet("footodelete")
	if resp.Count != 0 {
		fmt.Printf("footodelete was not deleted\n")
		kvs.KVDelete("footodelete")
	} else {
		fmt.Printf("footodelete was deleted\n")
	}
	resp = kvs.KVGet("footoput")
	if resp.Count != 0 {
		fmt.Printf("footoput was put\n")
		kvs.KVDelete("footoput")
	} else {
		fmt.Printf("footoput was not put\n")
	}

	for _, ev := range resp.Kvs {
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	}

	resp = kvs.KVGetWithPrefix("fo")
	for _, ev := range resp.Kvs {
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	}

}
