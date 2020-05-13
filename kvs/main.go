package main

import (
	"fmt"

	//	"./lib"
	kvs "kvs/lib"
)

func main() {

	// Test Write if not exists
	fmt.Println("Test Write if not exists")
	kvs.KVPutTTL("foo", "boo", 14)
	r1 := kvs.KVPutIfNotExists("foo", "bee")
	if r1 != nil {
		fmt.Println("Ok. Existing key could not be overwritten")
	} else {
		fmt.Println("Error. Existing key was overwritten")
	}
	kvs.KVDelete("foo")

	r2 := kvs.KVPutIfNotExists("foo", "bee")
	if r2 != nil {
		fmt.Println("Error. Non-existing key could not be created")
	} else {
		fmt.Println("Ok. None-existing key was rwritten")

	}
	kvs.KVDelete("foo")
	fmt.Println()

	// Test write and delete PASS
	fmt.Println("Test write and delete PASS")
	kvs.KVDelete("footoput")
	kvs.KVPut("footodelete", "boo")
	kvs.KVPutAndDelete("footodelete", "footoput", "valueput")
	r3 := kvs.KVGet("footodelete")
	if r3.Count != 0 {
		fmt.Println("Err. footodelete was not deleted")
		kvs.KVDelete("footodelete")
	} else {
		fmt.Println("Ok. footodelete was deleted")
	}
	r4 := kvs.KVGet("footoput")
	if r4.Count != 0 {
		fmt.Println("Ok. footoput was put")
		kvs.KVDelete("footoput")
	} else {
		fmt.Println("Err. footoput was not put")
	}
	fmt.Println()

	// Test write and delete FAIL
	fmt.Println("Test write and delete FAIL")
	kvs.KVPut("footoput", "foo")
	kvs.KVPut("footodelete", "boo")
	kvs.KVPutAndDelete("footodelete", "footoput", "valueput")
	r5 := kvs.KVGet("footodelete")
	if r5.Count != 0 {
		fmt.Println("Ok. footodelete was not deleted")
		kvs.KVDelete("footodelete")
	} else {
		fmt.Println("Error. footodelete was deleted")
	}
	r6 := kvs.KVGet("footoput")
	if r6.Count != 0 {
		fmt.Println("Ok. footoput already existed")
		kvs.KVDelete("footoput")
	} else {
		fmt.Println("Err. footoput was not put")
	}
	fmt.Println()

	// Previous tests

	/*kvs.KVPutTTL("foo", "booo", 14)
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

	for _, ev := range resp.Kvs {
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	}

	resp = kvs.KVGetWithPrefix("fo")
	for _, ev := range resp.Kvs {
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	}*/

}
