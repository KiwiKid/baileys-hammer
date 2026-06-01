package main

import "testing"

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
