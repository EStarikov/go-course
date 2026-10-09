package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(
        context.Background(), syscall.SIGINT, syscall.SIGTERM)
    defer stop()

	from := flag.Int("from", -1, "id of first film")
	to := flag.Int("to", -1, "id of last film")
	workers := flag.Int("workers", 10, "num of workers")
	timeout := flag.Int("timeout", 5, "timeout for http request")
	flag.Parse()

	if *from < 0 {
		fmt.Println("id of first film must be positive")
		os.Exit(-1)
	}

	if *to < 0 {
		fmt.Println("id of last film must be positive")
		os.Exit(-1)
	}

	if *from > *to {
		fmt.Println("last index must be bigger than first index")
		os.Exit(-1)
	}

	if *workers < 0 {
		fmt.Println("number of workers must be positive")
		os.Exit(-1)
	}

	if *timeout < 0 {
		fmt.Println("timeout must be positive")
		os.Exit(-1)
	}

	load(*from, *to, *workers, *timeout, ctx)
}
