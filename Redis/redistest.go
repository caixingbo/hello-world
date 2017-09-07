package main

import (
	"github.com/go-redis/redis"
	"fmt"
	"time"
)
/*
DENIED Redis is running in protected mode because protected mode is enabled, no bind address was specified,
no authentication password is requested to clients. In this mode connections are only accepted from the loopback interface.
If you want to connect from external computers to Redis you may adopt one of the following solutions:
1) Just disable protected mode sending the command 'CONFIG SET protected-mode no' from the loopback interface by connecting
to Redis from the same host the server is running, however MAKE SURE Redis is not publicly accessible from internet if you do so.
Use CONFIG REWRITE to make this change permanent.
2) Alternatively you can just disable the protected mode by editing the Redis
configuration file, and setting the protected mode option to 'no', and then restarting the server.
3) If you started the server manually just for testing, restart it with the '--protected-mode no' option.
4) Setup a bind address or an authentication password. NOTE: You only need to do one of the above things in order
for the server to start accepting connections from the outside.
*/

var client *redis.Client

func Init() {
	client = redis.NewClient(&redis.Options{
		Addr:         "116.62.245.54:6379",
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	})
	//client.FlushDb()
}

func ExampleNewClient() {
	client := redis.NewClient(&redis.Options{
		Addr:     "116.62.245.54:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>
}


func ExampleClient() {
	err := client.Set("key3", "temo", 0).Err()
	if err != nil {
		panic(err)
	}

	//val, err := client.Get("name").Result()
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("key", val)

	val1, err := client.Get("key1").Result()
	if err == redis.Nil {
		fmt.Println("key1 does not exists")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key1", val1)
	}

	val2, err := client.Get("key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exists")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
	// Output: key value
	// key2 does not exists
}

func main(){
	//ExampleNewClient()
	Init()
	ExampleClient()
	//client.Del("name")
	//client.Save()
}

