package main

// This is an example of using go-iapclient for a purpose that is not usage by
// a golang http.Client. It uses GetToken() to retrieve the token, then inserts
// it into a curl commandline
// Example usage:
// go run iapwrap.go --client-id $CLIENT_ID -- -v -k https://$HOST

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/vitaminmoo/go-iapclient"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	cid  = kingpin.Flag("client-id", "OAuth Client ID").Required().String()
	args = kingpin.Arg("args", "remaining args").Strings()
)

func main() {
	kingpin.Parse()

	iap, err := iapclient.NewIAP(*cid, nil)
	if err != nil {
		log.Fatalf("Failed to create new IAP object: %v", err)
	}

	token, err := iap.GetToken(context.Background())
	if err != nil {
		log.Fatalf("Failed to get token: %v", err)
	}

	curl, err := exec.LookPath("curl")
	args := append([]string{"curl", "-H", fmt.Sprintf("Authorization: Bearer %s", token)}, *args...)
	env := os.Environ()

	if err := syscall.Exec(curl, args, env); err != nil {
		log.Fatal(err)
	}

}
