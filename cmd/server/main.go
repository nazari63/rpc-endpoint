package main

import (
	"crypto/ecdsa"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/flashbots/rpc-endpoint/server"
)

var (
	defaultListenAddress = "127.0.0.1:9000"
	defaultProxyUrl      = "http://127.0.0.1:8545"
	defaultRelayUrl      = "https://relay.flashbots.net"
	defaultRedisUrl      = "localhost:6379"

	version = "dev" // is set during build process
)

var versionPtr = flag.Bool("version", false, "just print the program version")
var listenAddress = flag.String("listen", getEnvOrDefault("LISTEN_ADDR", defaultListenAddress), "Listen address")
var proxyUrl = flag.String("proxy", getEnvOrDefault("PROXY_URL", defaultProxyUrl), "URL for default JSON-RPC proxy target (eth node, Infura, etc.)")
var redisUrl = flag.String("redis", getEnvOrDefault("REDIS_URL", defaultRedisUrl), "URL for Redis")

// Flags for using the relay
var relayUrl = flag.String("relayUrl", getEnvOrDefault("RELAY_URL", defaultRelayUrl), "URL for relay")
var relaySigningKey = flag.String("signingKey", os.Getenv("RELAY_SIGNING_KEY"), "Signing key for relay requests")

func main() {
	var key *ecdsa.PrivateKey
	var err error

	flag.Parse()

	// Perhaps print only the version
	if *versionPtr {
		fmt.Printf("rpc-endpoint %s\n", version)
		return
	}

	log.Printf("rpc-endpoint %s\n", version)

	if *relaySigningKey == "" {
		log.Fatal("Cannot use the relay without a signing key.")
	}

	if strings.HasPrefix(*relaySigningKey, "0x") {
		*relaySigningKey = (*relaySigningKey)[2:]
	}

	key, err = crypto.HexToECDSA(*relaySigningKey)
	if err != nil {
		log.Fatal("Invalid relay signing key", err)
	}

	// Start the endpoint
	s, err := server.NewRpcEndPointServer(version, *listenAddress, *proxyUrl, *relayUrl, key, *redisUrl)
	if err != nil {
		log.Fatal("Server init error:", err)
	}
	s.Start()
}

func getEnvOrDefault(key string, defaultValue string) string {
	ret := os.Getenv(key)
	if ret == "" {
		ret = defaultValue
	}
	return ret
}
