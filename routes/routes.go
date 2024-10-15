package routes

import (
	"cognito-lambda-handler/internal/cognito"
	"cognito-lambda-handler/internal/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

// RegisterRoutes 関数はすべてのAPIルートを登録します
func RegisterRoutes(cognitoService *cognito.Service) *mux.Router {
	r := mux.NewRouter()

	// ルートの設定: handlersで定義したハンドラーを直接使用
	r.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) { handlers.SignUpHandler(w, r, cognitoService) }).Methods("POST")
	r.HandleFunc("/signin", func(w http.ResponseWriter, r *http.Request) { handlers.SignInHandler(w, r, cognitoService) }).Methods("POST")
	r.HandleFunc("/confirm", func(w http.ResponseWriter, r *http.Request) { handlers.ConfirmSignUpHandler(w, r, cognitoService) }).Methods("POST")
	r.HandleFunc("/forgot-password", func(w http.ResponseWriter, r *http.Request) { handlers.ForgotPasswordHandler(w, r, cognitoService) }).Methods("POST")
	r.HandleFunc("/reset-password", func(w http.ResponseWriter, r *http.Request) { handlers.ResetPasswordHandler(w, r, cognitoService) }).Methods("POST")

	return r
}
