package snapshot

import (
	"fmt"
	"io"
	"text/tabwriter"
	"time"
)

// Report writes a human-readable summary of a PortSnapshot to w.
func Report(w io.Writer, snap *PortSnapshot) error {
	if snap == nil {
		_, err := fmt.Fprintln(w, "no snapshot available")
		return err
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	defer tw.Flush()

	fmt.Fprintf(tw, "Snapshot at %s\n", snap.Timestamp.Format(time.RFC3339))
	fmt.Fprintf(tw, "PROTOCOL\tADDRESS\tSTATE\n")
	fmt.Fprintf(tw, "--------\t-------\t-----\n")

	for _, p := range snap.Ports {
		fmt.Fprintf(tw, "%s\t%s\t%s\n", p.Protocol, p.LocalAddress, p.State)
	}

	return nil
}

// DiffReport writes a human-readable diff between two consecutive snapshots.
func DiffReport(w io.Writer, appeared, disappeared []portscanner.PortEntry) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	defer tw.Flush()

	if len(appeared) == 0 && len(disappeared) == 0 {
		_, err := fmt.Fprintln(tw, "no changes detected")
		return err
	}

	for _, p := range appeared {
		fmt.Fprintf(tw, "+ %s\t%s\t%s\n", p.Protocol, p.LocalAddress, p.State)
	}

	for _, p := range disappeared {
		fmt.Fprintf(tw, "- %s\t%s\t%s\n", p.Protocol, p.LocalAddress, p.State)
	}

	return nil
}
