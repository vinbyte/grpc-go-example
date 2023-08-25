package main

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"sync"
	"time"

	pb "client/helloworld"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	call = flag.Int("call", 100, "number of grpc calls")
	kacp = keepalive.ClientParameters{
		Time:                10 * time.Second,       // send pings every 10 seconds if there is no activity
		Timeout:             100 * time.Millisecond, // wait 1 second for ping ack before considering the connection dead
		PermitWithoutStream: true,                   // send pings even without active streams
	}
	retryPolicy = `{
		"methodConfig": [{
		  "name": [{"service": "helloworld.Greeter"}],
		  "waitForReady": true,
		  "retryPolicy": {
			  "MaxAttempts": 3,
			  "InitialBackoff": ".01s",
			  "MaxBackoff": ".01s",
			  "BackoffMultiplier": 1.0,
			  "RetryableStatusCodes": [ "UNAVAILABLE" ]
		  }
		}]}`
)

func main() {
	flag.Parse()

	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(retryPolicy),
		grpc.WithKeepaliveParams(kacp),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// call many grpc server at the same time
	wg := sync.WaitGroup{}
	for i := 0; i < *call; i++ {
		wg.Add(1)
		go func() {
			callSayHello(c)
			wg.Done()
		}()
	}
	wg.Wait()

	select {} // Block forever; run with GODEBUG=http2debug=2 to observe ping frames and GOAWAYs due to idleness.
}

func callSayHello(c pb.GreeterClient) {
	// Contact the server and print out its response.
	r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: "Jhon Doe"})
	if err != nil {
		log.Printf("cannot call server: %v \n", err)
	}
	log.Printf("RPC response: %s", r.GetMessage())
}

func randomSleepTime() (sleep time.Duration) {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(3) // n will be between 0 and 3
	sleep = time.Duration(n) * time.Second

	return
}
