package main

import (
	"bitbucket.org/chrj/smtpd"
	"bufio"
	"flag"
	"fmt"
	"os"
	"path"
	"sync/atomic"
)

var dir string
var counter uint64 = 0

func handler(peer smtpd.Peer, env smtpd.Envelope) error {
	atomic.AddUint64(&counter, 1)
	filePath := path.Join(dir, fmt.Sprintf("%s.%06d", env.Recipients[0], counter))
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not create file %s", file)
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	defer writer.Flush()

	_, err = writer.Write(env.Data)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not write to file %s", file)
		return err
	}
	fmt.Printf("Wrote %s\n", filePath)
	return nil
}

func main() {
	var portPtr = flag.Int("port", 2500, "the port to listen to")
	var hostPtr = flag.String("host", "0.0.0.0", "the host ip to bind to")
	flag.StringVar(&dir, "output", os.TempDir(), "the directory to write mails to")
	flag.Parse()

	fmt.Printf("binding to port %d on %s\nWriting emails to %s\n", *portPtr, *hostPtr, dir)

	server := &smtpd.Server{
		WelcomeMessage: "Fakeserver 1.0",
		Handler:        handler,
	}

	server.ListenAndServe(fmt.Sprintf("%s:%d", *hostPtr, *portPtr))
}
