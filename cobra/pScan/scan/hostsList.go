package scan

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"slices"
)

var (
	ErrExists    = errors.New("Host already in the list")
	ErrNotExists = errors.New("Host not in the list")
)

// HostsList represents a list of hosts to run port scan
type HostsList struct {
	Hosts []string
}

// search searches for a host in the list.
// If the host is in the list, search returns true and the index of the host in the list.
// Else, search returns false and -1.
func (hl *HostsList) search(host string) (bool, int) {
	slices.Sort(hl.Hosts)
	i := slices.Index(hl.Hosts, host)
	if i >= 0 && i < len(hl.Hosts) && hl.Hosts[i] == host {
		return true, i
	}
	return false, -1
}

// Add adds a host to the list
func (hl *HostsList) Add(host string) error {
	if found, _ := hl.search(host); found {
		return fmt.Errorf("%w: %s", ErrExists, host)
	}
	hl.Hosts = append(hl.Hosts, host)
	return nil
}

// Remove deletes a host from the list
func (hl *HostsList) Remove(host string) error {
	if found, i := hl.search(host); found {
		hl.Hosts = append(hl.Hosts[:i], hl.Hosts[i+1:]...)
		return nil
	}
	return fmt.Errorf("%w: %s", ErrNotExists, host)
}

// Load obtains hosts from a hosts file
func (hl *HostsList) Load(hostsFile string) error {
	f, err := os.Open(hostsFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		hl.Hosts = append(hl.Hosts, scanner.Text())
	}
	return nil
}

// Save saves hosts to a hosts file
func (hl *HostsList) Save(hostsFile string) error {
	output := ""
	for _, h := range hl.Hosts {
		output += fmt.Sprintln(h)
	}
	return os.WriteFile(hostsFile, []byte(output), 0644)
}
