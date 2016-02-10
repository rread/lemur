package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/net/context"

	"github.intel.com/hpdd/logging/alert"
	"github.intel.com/hpdd/logging/audit"
	"github.intel.com/hpdd/logging/debug"
	"github.intel.com/hpdd/policy/pdm/lhsmd/agent"

	_ "github.intel.com/hpdd/policy/pdm/lhsmd/transport/grpc"
	// _ "github.intel.com/hpdd/policy/pdm/lhsmd/transport/queue"
)

func init() {
	flag.Var(debug.FlagVar())
}

func interruptHandler(once func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM)

	go func() {
		stopping := false
		for sig := range c {
			debug.Printf("signal received: %s", sig)
			if !stopping {
				stopping = true
				once()
			}
		}
	}()

}

func main() {
	flag.Parse()

	if debug.Enabled() {
		// Set this so that plugins can use it without needing
		// to mess around with plugin args.
		os.Setenv(debug.EnableEnvVar, "true")
	}
	// Setting the prefix helps us to track down deprecated calls to log.*
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(audit.Writer().Prefix("DEPRECATED "))

	conf := agent.ConfigInitMust()

	debug.Printf("current configuration:\n%v", conf.String())

	ct, err := agent.New(conf)
	if err != nil {
		alert.Fatalf("Error creating agent: %s", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	interruptHandler(func() {
		ct.Stop()
		cancel()
	})

	if err := ct.Start(ctx); err != nil {
		alert.Fatalf("Error in HsmAgent.Start(): %s", err)
	}
}
