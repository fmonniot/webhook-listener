// dns
package webhookListener

import (
	"fmt"
	"os/exec"
)

type DnsHandler struct {
	dnsHandlerConfig DnsHandlerConfig
}

type DnsHandlerConfig struct {
}

func (d *DnsHandler) Call(message GitlabPushMessage) {
	cmd := exec.Command("/usr/bin/git", "pull")
	run(cmd)

	cmd = exec.Command("/usr/sbin/named-checkconf")
	run(cmd)

	cmd = exec.Command("/etc/init.d/bind9", "restart")
	run(cmd)
}

func run(cmd *exec.Cmd) {
	cmd.Dir = "/etc/bind"
	err := cmd.Run()
	if err != nil {
		fmt.Print(err)
	}
}
