package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type DB struct{ db *sql.DB }

// Artifact is a single build output (binary, package, release archive)
// tracked in the registry. Status is one of: available, archived, broken,
// staging. SizeBytes is stored as an integer.
type Artifact struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	BuildID     string `json:"build_id"`
	Version     string `json:"version"`
	Platform    string `json:"platform"`
	SizeBytes   int    `json:"size_bytes"`
	Checksum    string `json:"checksum"`
	Status      string `json:"status"`
	DownloadURL string `json:"download_url"`
	CreatedAt   string `json:"created_at"`
}

func Open(d string) (*DB, error) {
	if err := os.MkdirAll(d, 0755); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", filepath.Join(d, "ironworks.db")+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, err
	}
	db.Exec(`CREATE TABLE IF NOT EXISTS artifacts(
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		build_id TEXT DEFAULT '',
		version TEXT DEFAULT '',
		platform TEXT DEFAULT '',
		size_bytes INTEGER DEFAULT 0,
		checksum TEXT DEFAULT '',
		status TEXT DEFAULT 'available',
		download_url TEXT DEFAULT '',
		created_at TEXT DEFAULT(datetime('now'))
	)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_artifacts_status ON artifacts(status)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_artifacts_platform ON artifacts(platform)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_artifacts_name ON artifacts(name)`)
	db.Exec(`CREATE TABLE IF NOT EXISTS extras(
		resource TEXT NOT NULL,
		record_id TEXT NOT NULL,
		data TEXT NOT NULL DEFAULT '{}',
		PRIMARY KEY(resource, record_id)
	)`)
	return &DB{db: db}, nil
}

func (d *DB) Close() error { return d.db.Close() }

func genID() string { return fmt.Sprintf("%d", time.Now().UnixNano()) }
func now() string   { return time.Now().UTC().Format(time.RFC3339) }

func (d *DB) Create(e *Artifact) error {
	e.ID = genID()
	e.CreatedAt = now()
	if e.Status == "" {
		e.Status = "available"
	}
	_, err := d.db.Exec(
		`INSERT INTO artifacts(id, name, build_id, version, platform, size_bytes, checksum, status, download_url, created_at)
		 VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		e.ID, e.Name, e.BuildID, e.Version, e.Platform, e.SizeBytes, e.Checksum, e.Status, e.DownloadURL, e.CreatedAt,
	)
	return err
}

func (d *DB) Get(id string) *Artifact {
	var e Artifact
	err := d.db.QueryRow(
		`SELECT id, name, build_id, version, platform, size_bytes, checksum, status, download_url, created_at
		 FROM artifacts WHERE id=?`,
		id,
	).Scan(&e.ID, &e.Name, &e.BuildID, &e.Version, &e.Platform, &e.SizeBytes, &e.Checksum, &e.Status, &e.DownloadURL, &e.CreatedAt)
	if err != nil {
		return nil
	}
	return &e
}

func (d *DB) List() []Artifact {
	rows, _ := d.db.Query(
		`SELECT id, name, build_id, version, platform, size_bytes, checksum, status, download_url, created_at
		 FROM artifacts ORDER BY created_at DESC`,
	)
	if rows == nil {
		return nil
	}
	defer rows.Close()
	var o []Artifact
	for rows.Next() {
		var e Artifact
		rows.Scan(&e.ID, &e.Name, &e.BuildID, &e.Version, &e.Platform, &e.SizeBytes, &e.Checksum, &e.Status, &e.DownloadURL, &e.CreatedAt)
		o = append(o, e)
	}
	return o
}

func (d *DB) Update(e *Artifact) error {
	_, err := d.db.Exec(
		`UPDATE artifacts SET name=?, build_id=?, version=?, platform=?, size_bytes=?, checksum=?, status=?, download_url=?
		 WHERE id=?`,
		e.Name, e.BuildID, e.Version, e.Platform, e.SizeBytes, e.Checksum, e.Status, e.DownloadURL, e.ID,
	)
	return err
}

func (d *DB) Delete(id string) error {
	_, err := d.db.Exec(`DELETE FROM artifacts WHERE id=?`, id)
	return err
}

func (d *DB) Count() int {
	var n int
	d.db.QueryRow(`SELECT COUNT(*) FROM artifacts`).Scan(&n)
	return n
}

func (d *DB) Search(q string, filters map[string]string) []Artifact {
	where := "1=1"
	args := []any{}
	if q != "" {
		where += " AND (name LIKE ? OR build_id LIKE ? OR version LIKE ?)"
		s := "%" + q + "%"
		args = append(args, s, s, s)
	}
	if v, ok := filters["status"]; ok && v != "" {
		where += " AND status=?"
		args = append(args, v)
	}
	if v, ok := filters["platform"]; ok && v != "" {
		where += " AND platform=?"
		args = append(args, v)
	}
	if v, ok := filters["name"]; ok && v != "" {
		where += " AND name=?"
		args = append(args, v)
	}
	rows, _ := d.db.Query(
		`SELECT id, name, build_id, version, platform, size_bytes, checksum, status, download_url, created_at
		 FROM artifacts WHERE `+where+`
		 ORDER BY created_at DESC`,
		args...,
	)
	if rows == nil {
		return nil
	}
	defer rows.Close()
	var o []Artifact
	for rows.Next() {
		var e Artifact
		rows.Scan(&e.ID, &e.Name, &e.BuildID, &e.Version, &e.Platform, &e.SizeBytes, &e.Checksum, &e.Status, &e.DownloadURL, &e.CreatedAt)
		o = append(o, e)
	}
	return o
}

// Stats returns total artifacts, total bytes stored across all artifacts,
// counts by status, counts by platform, and the count of distinct
// artifact names.
func (d *DB) Stats() map[string]any {
	m := map[string]any{
		"total":          d.Count(),
		"total_bytes":    0,
		"distinct_names": 0,
		"by_status":      map[string]int{},
		"by_platform":    map[string]int{},
	}

	var totalBytes int64
	d.db.QueryRow(`SELECT COALESCE(SUM(size_bytes), 0) FROM artifacts`).Scan(&totalBytes)
	m["total_bytes"] = totalBytes

	var distinct int
	d.db.QueryRow(`SELECT COUNT(DISTINCT name) FROM artifacts`).Scan(&distinct)
	m["distinct_names"] = distinct

	if rows, _ := d.db.Query(`SELECT status, COUNT(*) FROM artifacts GROUP BY status`); rows != nil {
		defer rows.Close()
		by := map[string]int{}
		for rows.Next() {
			var s string
			var c int
			rows.Scan(&s, &c)
			by[s] = c
		}
		m["by_status"] = by
	}

	if rows, _ := d.db.Query(`SELECT platform, COUNT(*) FROM artifacts WHERE platform != '' GROUP BY platform`); rows != nil {
		defer rows.Close()
		by := map[string]int{}
		for rows.Next() {
			var s string
			var c int
			rows.Scan(&s, &c)
			by[s] = c
		}
		m["by_platform"] = by
	}

	return m
}

// ─── Extras ───────────────────────────────────────────────────────

func (d *DB) GetExtras(resource, recordID string) string {
	var data string
	err := d.db.QueryRow(
		`SELECT data FROM extras WHERE resource=? AND record_id=?`,
		resource, recordID,
	).Scan(&data)
	if err != nil || data == "" {
		return "{}"
	}
	return data
}

func (d *DB) SetExtras(resource, recordID, data string) error {
	if data == "" {
		data = "{}"
	}
	_, err := d.db.Exec(
		`INSERT INTO extras(resource, record_id, data) VALUES(?, ?, ?)
		 ON CONFLICT(resource, record_id) DO UPDATE SET data=excluded.data`,
		resource, recordID, data,
	)
	return err
}

func (d *DB) DeleteExtras(resource, recordID string) error {
	_, err := d.db.Exec(
		`DELETE FROM extras WHERE resource=? AND record_id=?`,
		resource, recordID,
	)
	return err
}

func (d *DB) AllExtras(resource string) map[string]string {
	out := make(map[string]string)
	rows, _ := d.db.Query(
		`SELECT record_id, data FROM extras WHERE resource=?`,
		resource,
	)
	if rows == nil {
		return out
	}
	defer rows.Close()
	for rows.Next() {
		var id, data string
		rows.Scan(&id, &data)
		out[id] = data
	}
	return out
}
