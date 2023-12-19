package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"selfupdate.blockthrough.com"
)

func main() {
	fmt.Println("Hello, playground", os.Args)
	selfupdate.NewCliRunner(os.Args[0], os.Args[1:]...).Run(context.Background())
}
