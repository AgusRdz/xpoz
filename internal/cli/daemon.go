package cli

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
)

// startDetached re-execs the current binary as a detached background process
// with the given port and pre-resolved subdomain. Writes a PID file when name
// is non-empty and prints status to out.
func startDetached(out io.Writer, port int, sub, fullURL, name string) error {
	args := buildChildArgs(port, sub, name)
	proc := exec.Command(os.Args[0], args...)
	proc.Stdin = nil
	proc.Stdout = nil
	proc.Stderr = nil
	setSysProcAttr(proc)

	if err := proc.Start(); err != nil {
		return fmt.Errorf("starting background process: %w", err)
	}

	if name != "" {
		if err := writePID(name, proc.Process.Pid); err != nil {
			printWarn(out, "could not write PID file: "+err.Error())
		}
	}

	stopHint := "kill " + strconv.Itoa(proc.Process.Pid)
	if name != "" {
		stopHint = "xpoz stop " + name
	}

	printSuccess(out, "Tunnel started in background")
	fmt.Fprintln(out, "  "+urlStyle.Render(fullURL))
	fmt.Fprintf(out, "  %s  localhost:%d\n", dimStyle.Render("forwarding to"), port)
	fmt.Fprintf(out, "  %s  %d\n", dimStyle.Render("PID"), proc.Process.Pid)
	fmt.Fprintln(out, dimStyle.Render("  Stop with: "+stopHint))
	fmt.Fprintln(out)
	return nil
}
