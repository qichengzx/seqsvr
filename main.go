package main

import (
	"encoding/json"
	"fmt"
	"github.com/qichengzx/seqsvr/service"
	"github.com/satori/go.uuid"
	"log"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
)

var (
	startID int64
	curID   int64
	maxID   int64
	uid     = uuid.Must(uuid.NewV4())
	host, _ = os.Hostname()
	conf    *service.Config
	mu      sync.Mutex
)

type H map[string]interface{}

type response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data H      `json:"data"`
}

func init() {
	conf = service.NewConfig()
	service.InitDB(conf.MySQL)
}

func main() {
	http.HandleFunc("/new", New)

	log.Println("ID server is listening on", conf.PORT)
	log.Fatal(http.ListenAndServe(conf.PORT, nil))
}

func New(w http.ResponseWriter, r *http.Request) {
	log.Println("request start")
	s := atomic.LoadInt64(&startID)

	if s == 0 || curID == maxID {
		mu.Lock()
		id, err := service.New(uid)
		mu.Unlock()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("get new id : %d for ip : %s，uuid : %s\n", id, host, uid)

		atomic.StoreInt64(&startID, id*conf.STEP)
		atomic.StoreInt64(&curID, id*conf.STEP)
		atomic.StoreInt64(&maxID, (id+1)*conf.STEP)

		log.Println("start id get =====>", startID)
	}

	if !atomic.CompareAndSwapInt64(&curID, curID, curID+1) {
		resp := response{
			Code: -1,
			Msg:  "",
		}
		res, _ := json.Marshal(resp)
		fmt.Fprintf(w, string(res))
		return
	}

	log.Printf("start id : %d, current id : %d, max id : %d\n", startID, curID, maxID)
	log.Println("request end")

	w.Header().Set("Content-Type", "application/json")
	resp := response{
		Code: 0,
		Msg:  "ok",
		Data: H{
			"id": curID,
		},
	}
	res, _ := json.Marshal(resp)
	fmt.Fprintf(w, string(res))
	return
}
