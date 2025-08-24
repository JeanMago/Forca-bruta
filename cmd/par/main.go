package main

import (
	"flag"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var password string

// Brute-force recursivo local (sem atomic ou canais)
func bruteRecursive(attempt []byte, depth, maxLen int) bool {
	if depth == maxLen {
		if string(attempt[:depth]) == password {
			return true
		}
		return false
	}

	for i := 0; i < len(charset); i++ {
		attempt[depth] = charset[i]
		if bruteRecursive(attempt, depth+1, maxLen) {
			return true
		}
	}
	return false
}

// Worker paralelo: cada thread fixa os dois primeiros caracteres
func worker(start1, start2 int, wg *sync.WaitGroup, maxLen int, found *bool, mu *sync.Mutex) {
	defer wg.Done()
	attempt := make([]byte, maxLen)
	attempt[0] = charset[start1]
	if maxLen > 1 {
		attempt[1] = charset[start2]
	}

	if bruteRecursive(attempt, 2, maxLen) {
		mu.Lock()
		*found = true
		mu.Unlock()
	}
}

func main() {
	flag.StringVar(&password, "password", "a1B9zAB", "Senha alvo (até 7 chars)")
	var threads int
	flag.IntVar(&threads, "threads", runtime.NumCPU(), "Número de goroutines")
	flag.Parse()

	if len(password) == 0 || len(password) > 7 {
		panic("A senha deve ter entre 1 e 7 caracteres")
	}
	for _, c := range password {
		if !strings.ContainsRune(charset, c) {
			panic("Senha contém caracteres fora de [A-Za-z0-9]")
		}
	}

	maxLen := len(password)
	runtime.GOMAXPROCS(runtime.NumCPU())
	start := time.Now()

	var wg sync.WaitGroup
	var found bool
	var mu sync.Mutex

	// Distribui threads pelos dois primeiros caracteres
	for i := 0; i < len(charset); i++ {
		for j := 0; j < len(charset); j++ {
			wg.Add(1)
			go worker(i, j, &wg, maxLen, &found, &mu)
		}
	}
	wg.Wait()
	elapsed := time.Since(start).Seconds()

	if found {
		fmt.Printf("🔑 Senha encontrada (paralela): %s\n", password)
	} else {
		fmt.Println("❌ Senha não encontrada")
	}
	fmt.Printf("⏱ Tempo paralelo: %.6f s (threads=%d)\n", elapsed, threads)
}
