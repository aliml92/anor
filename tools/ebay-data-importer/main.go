package main

import (
	"context"
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Parse()

	if err := run(context.Background(), os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
