package main

import (
	"flag"
	"fmt"
	"strings"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var password string

// Brute-force recursivo sequencial
func bruteSeq(attempt []byte, depth, maxLen int) bool {
	if depth == maxLen {
		if string(attempt[:depth]) == password {
			return true
		}
		return false
	}

	for i := 0; i < len(charset); i++ {
		attempt[depth] = charset[i]
		if bruteSeq(attempt, depth+1, maxLen) {
			return true
		}
	}
	return false
}

func main() {
	flag.StringVar(&password, "password", "a1B9zAB", "Senha alvo (até 7 chars)")
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
	start := time.Now()
	attempt := make([]byte, maxLen)

	ok := bruteSeq(attempt, 0, maxLen)
	elapsed := time.Since(start).Seconds()

	if ok {
		fmt.Printf("🔑 Senha encontrada (sequencial): %s\n", password)
	} else {
		fmt.Println("❌ Senha não encontrada")
	}
	fmt.Printf("⏱ Tempo sequencial: %.6f s\n", elapsed)
}
