package pkg

import (
	"database/sql"
	"fmt"
	"net/http"
	"practice/internal/auth"
	"practice/internal/user"

	"github.com/gorilla/mux"
)

func InitializeGlobalController(db *sql.DB, x *mux.Router) {
	userdb := user.ProvideDB(db)
	servicedb := auth.MakeUserRepoInterface(userdb)
	controlerauth := auth.SetupaAuthContoller(x, servicedb)
	controlerauth.InitializeController()

	err := http.ListenAndServe(":8000", x)
	if err != nil {
		fmt.Println("server start error")
	}
}
