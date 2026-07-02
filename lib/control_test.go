package lib

import "testing"

func TestNormalizeRateStatName(t *testing.T) {
	if got := normalizeRateStatName("ntcp.activePeers"); got != "tcp.activePeers" {
		t.Fatalf("normalizeRateStatName(ntcp.activePeers) = %q, want %q", got, "tcp.activePeers")
	}

	if got := normalizeRateStatName("NTCP.activePeers"); got != "tcp.activePeers" {
		t.Fatalf("normalizeRateStatName(NTCP.activePeers) = %q, want %q", got, "tcp.activePeers")
	}

	if got := normalizeRateStatName("udp.activePeers"); got != "udp.activePeers" {
		t.Fatalf("normalizeRateStatName(udp.activePeers) = %q, want unchanged", got)
	}
}
