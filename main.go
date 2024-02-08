package main

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"
)

const (
	packetSize = 10
	numWorkers = 3

	envInterval = "N"
	envWorkers  = "M"
	envDuration = "K"
)

func main() {
	interval, _ := strconv.Atoi(os.Getenv(envInterval))
	workers, _ := strconv.Atoi(os.Getenv(envWorkers))
	duration, _ := strconv.Atoi(os.Getenv(envDuration))

	packets := make(chan []int) // Канал для передачи пакетов данных

	// Packet Generator
	go func() {
		for {
			packet := generatePacket()
			packets <- packet
			time.Sleep(time.Millisecond * time.Duration(interval))
		}
	}()

	// Workers
	results := make(chan []int, 100) // Канал для передачи результатов обработки пакетов

	wg := &sync.WaitGroup{}
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			for packet := range packets {
				result := processPacket(packet)
				results <- result // Отправка результата в канал для аккумуляции
			}
			wg.Done()
		}()
	}

	// Accumulator
	go func() {
		var accumulator int
		for result := range results {
			accumulator += sum(result)
			fmt.Println("accumulator:", accumulator)
			time.Sleep(time.Second * time.Duration(duration)) // Ожидание указанной длительности
		}
	}()

	wg.Wait()

	close(packets)
	close(results)
}

func generatePacket() []int {
	packet := make([]int, packetSize)
	for i := 0; i < packetSize; i++ {
		packet[i] = rand.Intn(10)
	}
	return packet
}

func processPacket(packet []int) []int {
	sort.Slice(packet, func(i, j int) bool {
		return packet[i] > packet[j]
	})
	return packet[:3]
}

func sum(ints []int) int {
	result := 0
	for _, num := range ints {
		result += num
	}
	return result
}
