package handlers

import (
	"encoding/json"
	"net/http"
	"time"
	//"log"

	"github.com/pgmonzon/Yng_Servicios/core"
)


func PingHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	//var todos []models.Todo
	session := core.Session.Copy()
	session.Ping()
	defer session.Close()
	//respuesta, _ := json.Marshal("{Esta capa no tiene seguridad}")
	w.Header().Add("Access-Control-Allow-Origin", "*") //Porfavor no olvidarse de borrar esta mierda
	core.JSONError(w, r, start, "sin seguridad", http.StatusOK)
}

func SecuredPingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*") //Porfavor no olvidarse de borrar esta porqueria
	start := time.Now()
	if (!core.ChequearPermisos(r, "SecuredPing")){
		core.JSONError(w, r, start, "Este usuario no tiene permisos o hubo un error procesando tu request. Se ha contactado a un administrador.", http.StatusInternalServerError)
		return
	}
	session := core.Session.Copy()
	defer session.Close()
	respuesta, _ := json.Marshal("[respuesta: Estas autenticado.]")
	//log.Println(respuesta)
	core.JSONResponse(w, r, start, respuesta, http.StatusOK)
}
