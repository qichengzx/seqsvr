package service

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/satori/go.uuid"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
)

// Service is the basic information for a unique ID generator
type Service struct {
	startID uint64
	curID   uint64
	maxID   uint64
	uuid    uuid.UUID

	db *sql.DB

	config *Config

	mu sync.Mutex
}

// New returns a new service that can be used to generate
// IDs
func New(conf *Config) *Service {
	id, _ := uuid.NewV4()
	return &Service{
		uuid:   id,
		db:     conn(conf.MySQL),
		config: conf,
	}
}

type H map[string]interface{}

type response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data H      `json:"data"`
}

func (svr *Service) ServeHttp(w http.ResponseWriter, r *http.Request) {
	id, err := svr.nextID()
	if err != nil {
		resp := response{
			Code: -1,
			Msg:  "",
		}
		res, _ := json.Marshal(resp)
		fmt.Fprintf(w, string(res))
		return
	}

	resp := response{
		Code: 0,
		Msg:  "ok",
		Data: H{
			"id": id,
		},
	}
	res, _ := json.Marshal(resp)
	fmt.Fprintf(w, string(res))
	return
}

// NextID creates and returns a unique ID
func (svr *Service) NextID() (uint64, error) {
	return svr.nextID()
}

func (svr *Service) nextID() (uint64, error) {
	svr.mu.Lock()
	s := atomic.LoadUint64(&svr.startID)
	if s == 0 || svr.curID == svr.maxID {
		id, err := svr.newDBID()
		if err != nil {
			log.Fatal(err)
		}

		u64id := uint64(id)
		curid := u64id * svr.config.STEP
		n64id := u64id + 1

		atomic.StoreUint64(&svr.startID, curid)
		atomic.StoreUint64(&svr.curID, curid)
		atomic.StoreUint64(&svr.maxID, n64id*svr.config.STEP)
	}
	svr.mu.Unlock()

	if !atomic.CompareAndSwapUint64(&svr.curID, svr.curID, svr.curID+1) {
		//TODO got a bug here in Windows
		return 0, errors.New("CompareAndSwapUint64 faild")
	}

	log.Printf("start id : %d, current id : %d, max id : %d\n", svr.startID, svr.curID, svr.maxID)

	return svr.curID, nil
}
