package main

import (
	"fmt"
	"net"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func scanPort(host string, port int) bool {
	address := fmt.Sprintf("%s:%d", host, port)

	conn, err := net.DialTimeout(
		"tcp",
		address,
		500*time.Millisecond,
	)

	if err != nil {
		return false
	}

	conn.Close()

	return true
}

func main() {
	a := app.New()

	w := a.NewWindow("Port Scanner")

	// Fixed window size
	w.Resize(fyne.NewSize(500, 400))
	w.SetFixedSize(true)

	// --------------------
	// State
	// --------------------

	openPorts := []int{}

	// --------------------
	// UI
	// --------------------

	hostEntry := widget.NewEntry()
	hostEntry.SetText("127.0.0.1")
	hostEntry.SetPlaceHolder("127.0.0.1")

	portEntry := widget.NewEntry()
	portEntry.SetText("1-100")

	result := widget.NewMultiLineEntry()
	result.Disable()

	// Declare first because callback needs scanButton
	var scanButton *widget.Button

	scanButton = widget.NewButton("Scan", func() {

		// Reset state
		openPorts = []int{}

		result.SetText("Scanning...\n")

		scanButton.Disable()

		host := strings.TrimSpace(hostEntry.Text)

		parts := strings.Split(portEntry.Text, "-")

		if len(parts) != 2 {
			result.SetText("Invalid port range")
			scanButton.Enable()
			return
		}

		start, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
		end, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))

		if err1 != nil || err2 != nil {
			result.SetText("Invalid port range")
			scanButton.Enable()
			return
		}

		if start < 1 || end > 65535 || start > end {
			result.SetText("Invalid port range")
			scanButton.Enable()
			return
		}

		// Run scanner in background
		go func() {

			var wg sync.WaitGroup
			var mu sync.Mutex

			for port := start; port <= end; port++ {

				wg.Add(1)

				go func(p int) {
					defer wg.Done()

					if scanPort(host, p) {

						mu.Lock()

						openPorts = append(openPorts, p)

						mu.Unlock()
					}

				}(port)
			}

			wg.Wait()

			// Protect state while sorting/copying
			mu.Lock()

			sort.Ints(openPorts)

			ports := append([]int(nil), openPorts...)

			mu.Unlock()

			// Build result
			output := ""

			if len(ports) == 0 {
				output = "No open ports found."
			} else {
				output = "Open Ports:\n\n"

				for _, port := range ports {
					output += fmt.Sprintf(
						"Port %d OPEN\n",
						port,
					)
				}
			}

			result.SetText(output)
			scanButton.Enable()
		}()
	})

	// --------------------
	// Layout
	// --------------------

	content := container.NewVBox(
		widget.NewLabel("Host:"),
		hostEntry,

		widget.NewLabel("Ports:"),
		portEntry,

		scanButton,

		widget.NewLabel("Results:"),
		result,
	)

	w.SetContent(content)

	w.ShowAndRun()
}
