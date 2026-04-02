package main
import ("fmt";"log";"net/http";"os";"github.com/stockyard-dev/stockyard-ironworks/internal/server";"github.com/stockyard-dev/stockyard-ironworks/internal/store")
func main(){port:=os.Getenv("PORT");if port==""{port="10000"};dataDir:=os.Getenv("DATA_DIR");if dataDir==""{dataDir="./ironworks-data"}
db,err:=store.Open(dataDir);if err!=nil{log.Fatalf("ironworks: %v",err)};defer db.Close();srv:=server.New(db)
fmt.Printf("\n  Ironworks — build artifact storage\n  Dashboard:  http://localhost:%s/ui\n  API:        http://localhost:%s/api\n\n",port,port)
log.Printf("ironworks: listening on :%s",port);log.Fatal(http.ListenAndServe(":"+port,srv))}
