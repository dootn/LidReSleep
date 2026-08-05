# syntax=docker/dockerfile:1.7

# Builder image: cross-compiles the Windows single-file EXE.
# Pure stdlib, no third-party deps -> no network needed at build time.
FROM golang:1.23-alpine

WORKDIR /workspace
