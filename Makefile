BIN_DIR := bin
SEQ := $(BIN_DIR)/seq
PAR := $(BIN_DIR)/par
BENCH := $(BIN_DIR)/bench

GO := go
RUNS ?= 20
THREADS ?= 8
PASSWORD ?= a1B9z

all: $(SEQ) $(PAR) $(BENCH)

$(SEQ): cmd/seq/main.go go.mod
	@mkdir -p $(BIN_DIR)
	$(GO) build -o $(SEQ) ./cmd/seq

$(PAR): cmd/par/main.go go.mod
	@mkdir -p $(BIN_DIR)
	$(GO) build -o $(PAR) ./cmd/par

$(BENCH): tools/bench/main.go go.mod
	@mkdir -p $(BIN_DIR)
	$(GO) build -o $(BENCH) ./tools/bench

run-seq: $(SEQ)
	$(SEQ) -password $(PASSWORD)

run-par: $(PAR)
	$(PAR) -password $(PASSWORD) -threads $(THREADS)

bench: $(BENCH) $(SEQ) $(PAR)
	$(BENCH) -runs $(RUNS) -threads $(THREADS) -password $(PASSWORD)

clean:
	rm -rf $(BIN_DIR)
