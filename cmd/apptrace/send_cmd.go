package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"sourcegraph.com/sourcegraph/apptrace"
)

func init() {
	_, err := CLI.AddCommand("send",
		"send sample data to a remote collector",
		"The send command sends sample data to a remote collector server.",
		&sendCmd,
	)
	if err != nil {
		log.Fatal(err)
	}
}

type SendCmd struct {
	CollectorAddr  string `short:"c" long:"collector" description:"collector listen address" default:":7701"`
	CollectorProto string `short:"p" long:"proto" description:"collector protocol (tcp or tls)" default:"tcp"`
	ServerName     string `short:"s" long:"server-name" description:"server name (required for TLS)"`
	Debug          bool   `short:"d" long:"debug" description:"debug log"`
}

var sendCmd SendCmd

func (c *SendCmd) Execute(args []string) error {
	var rc *apptrace.RemoteCollector
	switch c.CollectorProto {
	case "tcp":
		rc = apptrace.NewRemoteCollector(c.CollectorAddr)
	case "tls":
		rc = apptrace.NewTLSRemoteCollector(c.CollectorAddr, &tls.Config{ServerName: c.ServerName})
	default:
		return fmt.Errorf("unknown proto: %q", c.CollectorProto)
	}
	rc.Debug = c.Debug

	rcc := &apptrace.ChunkedCollector{
		Collector:   rc,
		MinInterval: time.Second,
	}

	log.Println("Sending sample data...")
	if err := sampleData(rc); err != nil {
		return err
	}
	log.Println("Done sending sample data.")

	log.Println("Flushing chunked collector...")
	if err := rcc.Flush(); err != nil {
		return err
	}
	log.Println("Done flushing chunked collector.")

	return nil
}
