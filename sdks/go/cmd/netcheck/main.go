package main

import (
	"context"
	"fmt"
	"log"
	"time"

	boxlite "github.com/RussellLuo/boxlite/sdks/go"
)

func runCheck(box *boxlite.Box, label string, args ...string) {
	ctx := context.Background()
	result, err := box.Exec(ctx, args[0], args[1:]...)
	if err != nil {
		fmt.Printf("=== %s ===\nexec error: %v\n\n", label, err)
		return
	}
	fmt.Printf("=== %s ===\nexit=%d\nstdout:\n%s\nstderr:\n%s\n\n",
		label, result.ExitCode, result.Stdout, result.Stderr)
}

func main() {
	ctx := context.Background()

	fmt.Printf("boxlite.Version()=%s\n", boxlite.Version())

	rt, err := boxlite.NewRuntime()
	if err != nil {
		log.Fatalf("new runtime: %v", err)
	}
	defer rt.Close()

	box, err := rt.Create(ctx, "alpine:latest",
		boxlite.WithName("go-netcheck"),
		boxlite.WithMemory(512),
	)
	if err != nil {
		log.Fatalf("create: %v", err)
	}
	defer box.Close()

	if err := box.Start(ctx); err != nil {
		log.Fatalf("start: %v", err)
	}
	defer func() {
		if err := box.Stop(ctx); err != nil {
			fmt.Printf("stop error: %v\n", err)
		}
	}()

	runCheck(box, "whoami", "sh", "-lc", "id && uname -a")
	runCheck(box, "resolv.conf", "sh", "-lc", "cat /etc/resolv.conf")
	runCheck(box, "route+gateway", "sh", "-lc", "ip route 2>/dev/null || route -n 2>/dev/null || true; ping -c 1 -W 2 192.168.127.1 2>/dev/null || true")
	runCheck(box, "http example.com", "sh", "-lc", "wget -O- -T 10 http://example.com >/tmp/example.out 2>/tmp/example.err; code=$?; cat /tmp/example.out; echo; echo '---ERR---'; cat /tmp/example.err; exit $code")

	time.Sleep(500 * time.Millisecond)
	if metrics, err := box.Metrics(ctx); err == nil {
		fmt.Printf("metrics: sent=%d recv=%d tcp_conns=%d tcp_errs=%d\n",
			metrics.NetworkBytesSent,
			metrics.NetworkBytesReceived,
			metrics.NetworkTCPConns,
			metrics.NetworkTCPErrors,
		)
	}
}
