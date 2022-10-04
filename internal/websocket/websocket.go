package websoc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"practice/internal/comment"
	"practice/internal/user"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Servicewebsoc struct {
	DB user.DataRepo
}

type UserMessage struct {
	Connec  *websocket.Conn
	Massage *chan string
}

type ContollerWebsoc struct {
	X          *mux.Router
	S          Servicewebsoc
	Register   chan *UserMessage
	UnRegister chan *UserMessage
	Msk        *chan string
}

func (Cw ContollerWebsoc) WebSoctest(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	var msgchan = make(chan string)
	var errchan = make(chan error)
	var temp UserMessage = UserMessage{
		Connec:  conn,
		Massage: &msgchan,
	}
	Cw.Register <- &temp
	go temp.ConnecCheker(&errchan)
	for {
		select {
		case msgs := <-msgchan:
			if err = conn.WriteMessage(1, []byte(msgs)); err != nil {
				fmt.Println(err)
			}
		case err = <-errchan:
			Cw.UnRegister <- &temp
			fmt.Println(err)
			conn.Close()
			return
		}

	}
}

func (Us UserMessage) ConnecCheker(errch *chan error) {
	for {
		if _, _, err := Us.Connec.ReadMessage(); err != nil {
			*errch <- err
			break
		}
	}
}

func (Cw ContollerWebsoc) HubMassage() {
	var Userada map[*UserMessage]bool = make(map[*UserMessage]bool)
	var hold = *(Cw.Msk)
	for {
		select {
		case user := <-Cw.Register:
			Userada[user] = true
		case unser := <-Cw.UnRegister:
			delete(Userada, unser)
		case mskan := <-hold:
			for i := range Userada {
				new := *(i.Massage)
				new <- mskan
				fmt.Println("called")
			}
		}
	}
}

func (Cw ContollerWebsoc) PostMassageUser(w http.ResponseWriter, r *http.Request) {
	var commentuser comment.AUserCommentJson
	var hold = *(Cw.Msk)
	err := json.NewDecoder(r.Body).Decode(&commentuser)
	if err != nil {
		return
	}
	hold <- commentuser.Comment
	fmt.Println("succeful post")
}
func (Cw ContollerWebsoc) InitializeWebController() {
	go Cw.HubMassage()
	Cw.X.HandleFunc("/Websoc", Cw.WebSoctest).Methods(http.MethodGet)
	Cw.X.HandleFunc("/PostMassege", Cw.PostMassageUser).Methods(http.MethodPost)
}

func ProvideControlerWebsoc(X *mux.Router, S Servicewebsoc, Msk *chan string, Register chan *UserMessage, UnRegister chan *UserMessage) ContollerWebsoc {
	return ContollerWebsoc{X: X, S: S, Msk: Msk, Register: Register, UnRegister: UnRegister}
}
