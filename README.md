# Brute Force 5 caracteres (A–Z, a–z, 0–9) em Go

## Build
make all

## Uso
./bin/seq -password a1B9z
./bin/par -password a1B9z -threads 8
./bin/bench -runs 20 -threads 8 -password a1B9z
