package main

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"strings"
	"time"
)

const defaultLitestreamConfigPath = "/etc/litestream.yml"

type LitestreamBackupStatus struct {
	DatabasePath       string
	ConfigPath         string
	ConfigReadable     bool
	ConfigReadout      string
	ConfigError        string
	EndpointSet        bool
	BucketSet          bool
	AccessKeySet       bool
	SecretKeySet       bool
	LitestreamPresent  bool
	StartedWithApp     bool
	SnapshotsReachable bool
	LastSnapshotAt     *time.Time
	LastSnapshotSize   string
	Error              string
}

func (s LitestreamBackupStatus) Configured() bool {
	return s.DatabasePath != "" && s.EndpointSet && s.BucketSet && s.AccessKeySet && s.SecretKeySet
}

func (s LitestreamBackupStatus) Healthy() bool {
	return s.Configured() && s.LitestreamPresent && s.StartedWithApp && s.SnapshotsReachable && s.LastSnapshotAt != nil && s.Error == ""
}

func litestreamBackupStatus() LitestreamBackupStatus {
	status := LitestreamBackupStatus{
		DatabasePath:   databasePathFromURL(os.Getenv("DATABASE_URL")),
		ConfigPath:     envOrDefault("LITESTREAM_CONFIG", defaultLitestreamConfigPath),
		EndpointSet:    strings.TrimSpace(os.Getenv("DB_REPLICA_URL")) != "",
		BucketSet:      strings.TrimSpace(os.Getenv("R2_BUCKET")) != "",
		AccessKeySet:   strings.TrimSpace(os.Getenv("R2_ACCESS_KEY_ID")) != "",
		SecretKeySet:   strings.TrimSpace(os.Getenv("R2_SECRET_ACCESS_KEY")) != "",
		StartedWithApp: strings.EqualFold(strings.TrimSpace(os.Getenv("LITESTREAM_ACTIVE")), "true"),
	}
	status.ConfigReadable, status.ConfigReadout, status.ConfigError = litestreamConfigReadout(status.ConfigPath)

	litestreamPath, err := exec.LookPath("litestream")
	status.LitestreamPresent = err == nil
	if !status.Configured() {
		status.Error = "Litestream environment is incomplete."
		return status
	}
	if !status.LitestreamPresent {
		status.Error = "Litestream binary is not available."
		return status
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	output, err := exec.CommandContext(ctx, litestreamPath, "snapshots", "-config", status.ConfigPath, status.DatabasePath).CombinedOutput()
	if ctx.Err() != nil {
		status.Error = "Timed out checking Litestream snapshots."
		return status
	}
	if err != nil {
		status.Error = strings.TrimSpace(string(output))
		if status.Error == "" {
			status.Error = err.Error()
		}
		return status
	}
	snapshot, err := latestLitestreamSnapshot(string(output))
	if err != nil {
		status.Error = err.Error()
		return status
	}
	status.SnapshotsReachable = true
	status.LastSnapshotAt = snapshot.CreatedAt
	status.LastSnapshotSize = snapshot.Size
	return status
}

func litestreamConfigReadout(configPath string) (bool, string, string) {
	if strings.TrimSpace(configPath) == "" {
		return false, "", "Litestream config path is empty."
	}
	config, err := os.ReadFile(configPath)
	if err != nil {
		return false, "", err.Error()
	}
	return true, sanitizeLitestreamConfig(string(config)), ""
}

func sanitizeLitestreamConfig(config string) string {
	lines := strings.Split(config, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(trimmed, "access-key-id:"):
			lines[i] = redactConfigLineValue(line)
		case strings.HasPrefix(trimmed, "secret-access-key:"):
			lines[i] = redactConfigLineValue(line)
		}
	}
	return strings.TrimSpace(strings.Join(lines, "\n"))
}

func redactConfigLineValue(line string) string {
	separator := strings.Index(line, ":")
	if separator == -1 {
		return line
	}
	return line[:separator+1] + " <redacted>"
}

func envOrDefault(name string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return fallback
	}
	return value
}

func databasePathFromURL(databaseURL string) string {
	dbPath := strings.TrimSpace(databaseURL)
	switch {
	case strings.HasPrefix(dbPath, "file://"):
		return strings.TrimPrefix(dbPath, "file://")
	case strings.HasPrefix(dbPath, "file:"):
		return strings.TrimPrefix(dbPath, "file:")
	default:
		return dbPath
	}
}

type litestreamSnapshot struct {
	CreatedAt *time.Time
	Size      string
}

func latestLitestreamSnapshot(output string) (litestreamSnapshot, error) {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	var latest litestreamSnapshot
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 5 || fields[0] == "replica" {
			continue
		}
		createdAt, err := time.Parse(time.RFC3339, fields[4])
		if err != nil {
			continue
		}
		if latest.CreatedAt == nil || createdAt.After(*latest.CreatedAt) {
			latest.CreatedAt = &createdAt
			latest.Size = fields[3]
		}
	}
	if latest.CreatedAt == nil {
		return latest, errors.New("No Litestream snapshots found.")
	}
	return latest, nil
}
