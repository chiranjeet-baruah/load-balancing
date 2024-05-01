package consistent

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"testing"
	"time"
)

// Maximum number of nodes in the ring
const MaxNodes = 10

func TestConsistentHashingRingSimulation(t *testing.T) {
	consistentHashingRing := NewRing()
	done := make(chan bool)

	mu := sync.Mutex{}

	var wg sync.WaitGroup
	// Add 3 for the three goroutines
	wg.Add(3)

	// Start a goroutine for requesting nodes
	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				return
			default:
				// generate a random request name
				requestName := fmt.Sprintf("request%d", rand.Intn(100))
				serverId, alternateServerId, err := consistentHashingRing.Get(requestName)
				if err != nil {
					if errors.Is(err, errors.New("no nodes available")) {
						log.Printf("No nodes available for request: %s\n", requestName)
					} else {
						log.Printf("Error: %s\n", err.Error())
					}
					continue
				}
				if alternateServerId == "" {
					log.Printf("Request: %s, Server: %s\n", requestName, serverId)
				} else {
					log.Printf("Request: %s, Server: %s, Alternate Server: %s\n", requestName, serverId, alternateServerId)
				}
				// random sleep time
				time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			}
		}
	}()

	// Start a goroutine for adding nodes
	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				return
			default:
				// generate a random server name
				serverName := fmt.Sprintf("server%d", rand.Intn(100))
				mu.Lock()
				if len(consistentHashingRing.Nodes) < MaxNodes {
					consistentHashingRing.AddNode(serverName)
					log.Printf("Added node: %s\n", serverName)
				}
				// print the current size of the ring
				log.Printf("Current size of the ring: %d\n", len(consistentHashingRing.Nodes))
				mu.Unlock()
				// random sleep time
				time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			}
		}
	}()

	// Start a goroutine for removing nodes
	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				return
			default:
				mu.Lock()
				if len(consistentHashingRing.Nodes) > 0 {
					randomIndex := rand.Intn(len(consistentHashingRing.Nodes))
					serverName := consistentHashingRing.Nodes[randomIndex].Id
					err := consistentHashingRing.RemoveNode(serverName)
					if err == nil {
						log.Printf("Removed node: %s\n", serverName)
					}
				}
				mu.Unlock()
				// random sleep time
				time.Sleep(time.Duration(rand.Intn(5000)) * time.Millisecond)
			}
		}
	}()

	// Run the simulation for n minutes
	time.Sleep(1 * time.Minute)

	// Signal the end of the simulation
	close(done)

	// Wait for all goroutines to finish
	wg.Wait()
}
