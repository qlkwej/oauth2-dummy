package server

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/models"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"
	"log"
	"net/http"
)


// Middleware
func ValidateToken(f http.HandlerFunc,srv *server.Server) http.HandlerFunc {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_, err := srv.ValidationBearerToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		f.ServeHTTP(writer, request)
	})



}



func Manager() {
	manager := manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	clientStore := store.NewClientStore()
	manager.MapClientStorage(clientStore)

	srv := server.NewDefaultServer(manager)
	srv.SetAllowGetAccessRequest(true)
	srv.SetClientInfoHandler(server.ClientFormHandler)
	manager.SetRefreshTokenCfg(manage.DefaultRefreshTokenCfg)




	// ROUTES

	// defined credentials path
	http.HandleFunc("/credentials", func(w http.ResponseWriter, r *http.Request) {
		clientId := uuid.New().String()[:8]
		clientSecret := uuid.New().String()[:8]
		fullClientId := uuid.New().String()
		fmt.Println("Full Client ID DEBUG: ", fullClientId)
		err := clientStore.Set(clientId,&models.Client{
			ID: clientId,
			Secret: clientSecret,
			Domain: "http://localhost:9094",
		})
		if err != nil {
			fmt.Println("Error set client store: ", err.Error())
		}

		w.Header().Set("Content-Type","application/json")
		json.NewEncoder(w).Encode(map[string]string{"CLIENT_ID":clientId,"CLIENT_SECRET":clientSecret})
	})

	// define token
	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		srv.HandleTokenRequest(w, r)
	})


	// defined protected path
	http.HandleFunc("/protected",ValidateToken (func (w http.ResponseWriter,r *http.Request){
		res, err := w.Write([]byte("Hello Im protected!"))
		if err != nil {
			log.Fatalln("Error: ", err)
		}
		fmt.Println("int val: ", res)
	}, srv))

	log.Fatal(http.ListenAndServe(":9090",nil))
	// END ROUTES

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error: ",err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error", re.Error.Error())
		
	})

}
