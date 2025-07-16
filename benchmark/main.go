package main

import (
	"fmt"
	"gocache"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type RunnerOpts struct {
	numWorkers      int
	numCommands     int
	runTimeDuration time.Duration
	pauseDuration   time.Duration
}

type Params struct {
	endTime    time.Time
	wg         *sync.WaitGroup
	keysCh     chan int
	RunnerOpts RunnerOpts
}

var cache gocache.GoCache
var params Params

func main() {
	cache = gocache.New()

	// Inserts 1M items into the cache per second with a TTL of 5s
	// Due to key randomization - around 1.5-2M elements will exist in the cache
	// (up to 8GB of memory used to hold that many elements)
	start, ttlDuration := time.Now(), time.Second*5
	var wg sync.WaitGroup

	runningOpts := RunnerOpts{
		numWorkers:      10,
		numCommands:     100_000,
		runTimeDuration: time.Minute,
		pauseDuration:   time.Second,
	}
	params = Params{
		endTime:    start.Add(runningOpts.runTimeDuration),
		wg:         &wg,
		keysCh:     make(chan int, 1000),
		RunnerOpts: runningOpts,
	}

	params.wg.Add(3)
	go generateKeys()
	go populateCache(ttlDuration)
	go launchRequests()
	params.wg.Wait()

	stats := cache.GetStats()
	fmt.Printf("\nExecution time: %fs\n", time.Since(start).Seconds())
	fmt.Printf("Total: %d, Gets: %d, Sets: %d, Has: %d, Deletes:%d\n",
		stats.TotalOperations, stats.NumGets, stats.NumSets, stats.NumHasChecks, stats.NumDeletes)

}

func populateCache(ttl time.Duration) {
	defer params.wg.Done()

	valueMap := map[int]string{
		0: "zero",
		1: "one",
		2: "two",
		3: "three",
		4: "four",
		5: "five",
		6: "six",
		7: "seven",
		8: "eight",
		9: "nine",
	}

	runner(func() {
		defer params.wg.Done()

		for range params.RunnerOpts.numCommands {
			key := <-params.keysCh
			cache.Set(strconv.Itoa(key), valueMap[key%10], ttl)
		}
	})
}

func launchRequests() {
	defer params.wg.Done()

	runner(func() {
		defer params.wg.Done()

		for range params.RunnerOpts.numCommands {
			keyStr := strconv.Itoa(<-params.keysCh)
			req, numRequests := []string{"GET", "HAS", "DELETE"}, 3
			requestIndex := rand.Intn(numRequests - 1)

			switch req[requestIndex] {
			case "GET":
				cache.Get(keyStr)
			case "HAS":
				cache.Has(keyStr)
			case "DELETE":
				cache.Delete(keyStr)
			default:
				panic("invalid request type")
			}
		}
	})
}

func runner(fn func()) {
	for time.Now().Before(params.endTime) {
		stepStart := time.Now()

		for range params.RunnerOpts.numWorkers {
			params.wg.Add(1)
			go fn()
			fmt.Printf("\rCache Size: %d ", cache.Size()) // Display size of cache while running
		}
		time.Sleep(params.RunnerOpts.pauseDuration - time.Since(stepStart))
	}
}

func generateKeys() {
	defer params.wg.Done()

	randomRange := params.RunnerOpts.numWorkers * params.RunnerOpts.numCommands * 10
	for time.Now().Before(params.endTime) {
		params.keysCh <- rand.Intn(randomRange)
	}
	close(params.keysCh)
}
