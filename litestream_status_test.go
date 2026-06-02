package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLatestLitestreamSnapshotUsesNewestCreatedAt(t *testing.T) {
	output := `replica  generation        index  size   created
s3       old-generation    0      1234   2026-06-01T01:00:00Z
s3       new-generation    4      57073  2026-06-01T03:40:45Z`

	snapshot, err := latestLitestreamSnapshot(output)
	if err != nil {
		t.Fatalf("latest snapshot: %v", err)
	}
	if snapshot.CreatedAt == nil {
		t.Fatal("expected created at")
	}
	if got := snapshot.CreatedAt.Format("2006-01-02T15:04:05Z"); got != "2026-06-01T03:40:45Z" {
		t.Fatalf("expected newest snapshot time, got %s", got)
	}
	if snapshot.Size != "57073" {
		t.Fatalf("expected newest snapshot size, got %q", snapshot.Size)
	}
}

func TestLatestLitestreamSnapshotReportsMissingSnapshots(t *testing.T) {
	_, err := latestLitestreamSnapshot("replica  generation  index  size  created\n")
	if err == nil {
		t.Fatal("expected missing snapshots error")
	}
}

func TestLitestreamConfigReadoutReturnsSanitizedConfig(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "litestream.yml")
	config := `dbs:
  - path: /mnt/volume/production.sqlite3
    replicas:
      - type: s3
        endpoint: ${DB_REPLICA_URL}
        bucket: ${R2_BUCKET}
        access-key-id: ${R2_ACCESS_KEY_ID}
        secret-access-key: ${R2_SECRET_ACCESS_KEY}
`
	if err := os.WriteFile(configPath, []byte(config), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	readable, readout, readErr := litestreamConfigReadout(configPath)
	if !readable {
		t.Fatal("expected config to be readable")
	}
	if readErr != "" {
		t.Fatalf("expected no read error, got %q", readErr)
	}
	if strings.Contains(readout, "${R2_ACCESS_KEY_ID}") || strings.Contains(readout, "${R2_SECRET_ACCESS_KEY}") {
		t.Fatalf("expected credential placeholders to be redacted, got:\n%s", readout)
	}
	if !strings.Contains(readout, "access-key-id: <redacted>") || !strings.Contains(readout, "secret-access-key: <redacted>") {
		t.Fatalf("expected redacted credential lines, got:\n%s", readout)
	}
	if !strings.Contains(readout, "path: /mnt/volume/production.sqlite3") {
		t.Fatalf("expected database path to remain visible, got:\n%s", readout)
	}
}

func TestLitestreamConfigReadoutReportsReadError(t *testing.T) {
	readable, readout, readErr := litestreamConfigReadout(filepath.Join(t.TempDir(), "missing.yml"))
	if readable {
		t.Fatal("expected missing config to be unreadable")
	}
	if readout != "" {
		t.Fatalf("expected empty readout, got %q", readout)
	}
	if readErr == "" {
		t.Fatal("expected read error")
	}
}
