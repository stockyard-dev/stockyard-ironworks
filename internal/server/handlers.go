package server
import("encoding/json";"net/http";"strconv";"github.com/stockyard-dev/stockyard-ironworks/internal/store")
func(s *Server)handleList(w http.ResponseWriter,r *http.Request){name:=r.URL.Query().Get("name");list,_:=s.db.List(name);if list==nil{list=[]store.Artifact{}};writeJSON(w,200,list)}
func(s *Server)handleCreate(w http.ResponseWriter,r *http.Request){var a store.Artifact;json.NewDecoder(r.Body).Decode(&a);if a.Name==""||a.Version==""{writeError(w,400,"name and version required");return};s.db.Create(&a);writeJSON(w,201,a)}
func(s *Server)handleDownload(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);s.db.Download(id);writeJSON(w,200,map[string]string{"status":"downloaded"})}
func(s *Server)handleDelete(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);s.db.Delete(id);writeJSON(w,200,map[string]string{"status":"deleted"})}
func(s *Server)handleOverview(w http.ResponseWriter,r *http.Request){m,_:=s.db.Stats();writeJSON(w,200,m)}
