package store
import ("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type Artifact struct {
	ID string `json:"id"`
	Name string `json:"name"`
	BuildID string `json:"build_id"`
	Version string `json:"version"`
	Platform string `json:"platform"`
	SizeBytes int `json:"size_bytes"`
	Checksum string `json:"checksum"`
	Status string `json:"status"`
	DownloadURL string `json:"download_url"`
	CreatedAt string `json:"created_at"`
}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"ironworks.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS artifacts(id TEXT PRIMARY KEY,name TEXT NOT NULL,build_id TEXT DEFAULT '',version TEXT DEFAULT '',platform TEXT DEFAULT '',size_bytes INTEGER DEFAULT 0,checksum TEXT DEFAULT '',status TEXT DEFAULT 'available',download_url TEXT DEFAULT '',created_at TEXT DEFAULT(datetime('now')))`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)Create(e *Artifact)error{e.ID=genID();e.CreatedAt=now();_,err:=d.db.Exec(`INSERT INTO artifacts(id,name,build_id,version,platform,size_bytes,checksum,status,download_url,created_at)VALUES(?,?,?,?,?,?,?,?,?,?)`,e.ID,e.Name,e.BuildID,e.Version,e.Platform,e.SizeBytes,e.Checksum,e.Status,e.DownloadURL,e.CreatedAt);return err}
func(d *DB)Get(id string)*Artifact{var e Artifact;if d.db.QueryRow(`SELECT id,name,build_id,version,platform,size_bytes,checksum,status,download_url,created_at FROM artifacts WHERE id=?`,id).Scan(&e.ID,&e.Name,&e.BuildID,&e.Version,&e.Platform,&e.SizeBytes,&e.Checksum,&e.Status,&e.DownloadURL,&e.CreatedAt)!=nil{return nil};return &e}
func(d *DB)List()[]Artifact{rows,_:=d.db.Query(`SELECT id,name,build_id,version,platform,size_bytes,checksum,status,download_url,created_at FROM artifacts ORDER BY created_at DESC`);if rows==nil{return nil};defer rows.Close();var o []Artifact;for rows.Next(){var e Artifact;rows.Scan(&e.ID,&e.Name,&e.BuildID,&e.Version,&e.Platform,&e.SizeBytes,&e.Checksum,&e.Status,&e.DownloadURL,&e.CreatedAt);o=append(o,e)};return o}
func(d *DB)Update(e *Artifact)error{_,err:=d.db.Exec(`UPDATE artifacts SET name=?,build_id=?,version=?,platform=?,size_bytes=?,checksum=?,status=?,download_url=? WHERE id=?`,e.Name,e.BuildID,e.Version,e.Platform,e.SizeBytes,e.Checksum,e.Status,e.DownloadURL,e.ID);return err}
func(d *DB)Delete(id string)error{_,err:=d.db.Exec(`DELETE FROM artifacts WHERE id=?`,id);return err}
func(d *DB)Count()int{var n int;d.db.QueryRow(`SELECT COUNT(*) FROM artifacts`).Scan(&n);return n}

func(d *DB)Search(q string, filters map[string]string)[]Artifact{
    where:="1=1"
    args:=[]any{}
    if q!=""{
        where+=" AND (name LIKE ?)"
        args=append(args,"%"+q+"%");
    }
    if v,ok:=filters["status"];ok&&v!=""{where+=" AND status=?";args=append(args,v)}
    rows,_:=d.db.Query(`SELECT id,name,build_id,version,platform,size_bytes,checksum,status,download_url,created_at FROM artifacts WHERE `+where+` ORDER BY created_at DESC`,args...)
    if rows==nil{return nil};defer rows.Close()
    var o []Artifact;for rows.Next(){var e Artifact;rows.Scan(&e.ID,&e.Name,&e.BuildID,&e.Version,&e.Platform,&e.SizeBytes,&e.Checksum,&e.Status,&e.DownloadURL,&e.CreatedAt);o=append(o,e)};return o
}

func(d *DB)Stats()map[string]any{
    m:=map[string]any{"total":d.Count()}
    rows,_:=d.db.Query(`SELECT status,COUNT(*) FROM artifacts GROUP BY status`)
    if rows!=nil{defer rows.Close();by:=map[string]int{};for rows.Next(){var s string;var c int;rows.Scan(&s,&c);by[s]=c};m["by_status"]=by}
    return m
}
