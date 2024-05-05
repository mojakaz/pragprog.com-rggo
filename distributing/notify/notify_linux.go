package notify

import "os/exec"

var command = exec.Command

func (n *Notify) Send() error {
	notifyCmdName := "wsl-notify-send.exe"
	notifyCmd, err := exec.LookPath(notifyCmdName)
	if err != nil {
		return err
	}
	notifyCommand := command(notifyCmd, "--category $WSL_DISTRO_NAME", "-u", n.severity.String(),
		n.title, n.message)
	return notifyCommand.Run()
}
