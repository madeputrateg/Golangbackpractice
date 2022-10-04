package pkg

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"practice/internal/auth"
	"practice/internal/user"
	websoc "practice/internal/websocket"

	"github.com/gorilla/mux"
)

func InitializeGlobalController(db *sql.DB, x *mux.Router) {
	var tempatmasuk = make(chan string)
	var Register = make(chan *websoc.UserMessage)
	var UnRegister = make(chan *websoc.UserMessage)
	userdb := user.ProvideDB(db)
	servicedb := auth.MakeUserRepoInterface(userdb)
	controlerauth := auth.SetupaAuthContoller(x, servicedb)
	controlerauth.InitializeController()
	testwebsoc := websoc.ProvideControlerWebsoc(x, websoc.Servicewebsoc{}, &tempatmasuk, Register, UnRegister)
	testwebsoc.InitializeWebController()
	controlerauth.InitializeControllerAuth()
	l, err := net.Listen("tcp4", ":12044")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	log.Fatal(http.Serve(l, x))
}
