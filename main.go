package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func calculateSMA(closes []float64, period int) []float64 {
	sma := make([]float64, len(closes)-period+1)
	sum := 0.0
	for i := 0; i < period; i++ {
		sum += closes[i]
	}
	sma[0] = sum / float64(period)
	for i := period; i < len(closes); i++ {
		sum += closes[i] - closes[i-period]
		sma[i-period+1] = sum / float64(period)
	}
	return sma
}

func movingAverageCrossover(closes []float64, shortPeriod, longPeriod int, signals chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	shortSMA := calculateSMA(closes, shortPeriod)
	longSMA := calculateSMA(closes, longPeriod)

	for i := 1; i < len(shortSMA) && i < len(longSMA) && i+longPeriod-1 < len(closes); i++ {
		if shortSMA[i] > longSMA[i] && shortSMA[i-1] <= longSMA[i-1] {
			signals <- fmt.Sprintf("Buy at index %d, Price: %.2f", i+longPeriod-1, closes[i+longPeriod-1])
		} else if shortSMA[i] < longSMA[i] && shortSMA[i-1] >= longSMA[i-1] {
			signals <- fmt.Sprintf("Sell at index %d, Price: %.2f", i+longPeriod-1, closes[i+longPeriod-1])
		}
	}
}

func main() {
	start := time.Now()

	file, err := os.Open("data.csv")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var closes []float64

	// Skip header
	scanner.Scan()

	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), ",")
		if len(fields) > 4 {
			if close, err := strconv.ParseFloat(fields[4], 64); err == nil {
				closes = append(closes, close)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	signals := make(chan string, 1000)
	var wg sync.WaitGroup

	// Number of goroutines to use
	numGoroutines := 4

	// Split the data into chunks
	chunkSize := len(closes) / numGoroutines

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		start := i * chunkSize
		end := (i + 1) * chunkSize
		if i == numGoroutines-1 {
			end = len(closes)
		}
		go movingAverageCrossover(closes[start:end], 5, 10, signals, &wg)
	}

	// Close channel when all goroutines are done
	go func() {
		wg.Wait()
		close(signals)
	}()

	buySignals := []string{}
	sellSignals := []string{}

	for signal := range signals {
		if strings.HasPrefix(signal, "Buy") {
			buySignals = append(buySignals, signal)
		} else {
			sellSignals = append(sellSignals, signal)
		}
	}

	fmt.Println("Buy Signals:")
	for _, signal := range buySignals {
		fmt.Println(signal)
	}

	fmt.Println("\nSell Signals:")
	for _, signal := range sellSignals {
		fmt.Println(signal)
	}

	fmt.Printf("Execution time: %s\n", time.Since(start))
}
