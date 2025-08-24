package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"time"
)

var (
	runs     int
	threads  int
	password string
	seqPath  string
	parPath  string
	outFile  string
)

func parseTimeSeconds(out string) (float64, error) {
	re := regexp.MustCompile(`([0-9]+\.[0-9]+)\s*s`)
	m := re.FindStringSubmatch(out)
	if len(m) < 2 {
		return 0, fmt.Errorf("tempo não encontrado em: %q", out)
	}
	return strconv.ParseFloat(m[1], 64)
}

func runCmd(path string, args ...string) (string, error) {
	cmd := exec.Command(path, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil { return "", err }
	stderr, err := cmd.StderrPipe()
	if err != nil { return "", err }
	if err := cmd.Start(); err != nil { return "", err }
	b := bufio.NewScanner(stdout)
	var out string
	for b.Scan() { out += b.Text() + "\\n" }
	be := bufio.NewScanner(stderr)
	for be.Scan() { out += be.Text() + "\\n" }
	if err := cmd.Wait(); err != nil { return out, err }
	return out, nil
}

func main() {
	flag.IntVar(&runs, "runs", 20, "quantidade de execuções de cada versão")
	flag.IntVar(&threads, "threads", 8, "threads para versão paralela")
	flag.StringVar(&password, "password", "a1B9z", "senha alvo")
	flag.StringVar(&seqPath, "seq", "./bin/seq", "caminho do binário sequencial")
	flag.StringVar(&parPath, "par", "./bin/par", "caminho do binário paralelo")
	flag.StringVar(&outFile, "out", "resultados.csv", "arquivo de saída (CSV)")
	flag.Parse()

	f, err := os.Create(outFile)
	if err != nil { panic(err) }
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	w.Write([]string{"timestamp", "run", "password", "threads", "seq_time_s", "par_time_s", "speedup_seq_over_par"})

	for i := 1; i <= runs; i++ {
		fmt.Printf("Execução %d/%d...\\n", i, runs)

		outSeq, err := runCmd(seqPath, "-password", password)
		if err != nil { fmt.Println("[seq] erro:", err) }
		seqT, _ := parseTimeSeconds(outSeq)

		outPar, err := runCmd(parPath, "-password", password, "-threads", strconv.Itoa(threads))
		if err != nil { fmt.Println("[par] erro:", err) }
		parT, _ := parseTimeSeconds(outPar)

		speedup := ""
		if parT > 0 {
			speedup = fmt.Sprintf("%.6f", seqT/parT)
		}

		w.Write([]string{
			time.Now().Format(time.RFC3339),
			strconv.Itoa(i),
			password,
			strconv.Itoa(threads),
			fmt.Sprintf("%.6f", seqT),
			fmt.Sprintf("%.6f", parT),
			speedup,
		})
		w.Flush()
	}

	fmt.Println("Arquivo gerado:", outFile)
}
