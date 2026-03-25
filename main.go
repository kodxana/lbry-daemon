package main

import (
	"lbry/daemon/dht"
	"lbry/daemon/rpc"
	"lbry/daemon/stream"
	"log"
)

func main() {
	// Start DHT node
	dhtNode, err := dht.NewNode(4444, 3333)
	if err != nil {
		log.Printf("Warning: DHT init failed: %v (P2P streaming disabled)", err)
	}

	var streamMgr *stream.Manager

	if dhtNode != nil {
		if err := dhtNode.Start(); err != nil {
			log.Printf("Warning: DHT start failed: %v", err)
		} else {
			// Bootstrap in background
			go func() {
				if err := dhtNode.Bootstrap(); err != nil {
					log.Printf("Warning: DHT bootstrap failed: %v", err)
				}
			}()

			// Start streaming server
			streamMgr = stream.NewManager(dhtNode)
			if err := streamMgr.StartHTTP(5280); err != nil {
				log.Printf("Warning: streaming server failed: %v", err)
			}
		}
	}

	// Set stream manager for RPC handlers
	rpc.SetStreamManager(streamMgr)

	// Start RPC server (blocks)
	rpc.StartServer(rpc.CreateServer(), 5279)
}
