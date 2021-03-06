package models

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

//Todo struct to todo
type Usuario struct {
	ID        bson.ObjectId `bson:"_id" json:"id"`
	User      string        `json:"user"`
	Email     string        `json:"email"`
	PassMD    int32         `json:"md5"`
	Facebook  UsuarioFacebook `json:"facebook"`
	Activado  bool		`json:"activado"`
	Codigo    string	`json:"codigo"`
	Creacion  time.Time	`json:"creacion"`
	Rol	  bson.ObjectId `json:"rol"`
}

type UsuarioFacebook struct {
	ID	string	`json:"id"`
	Name	string	`json:"name"`
	Link	string	`json:"link"`
	Email	string	`json:"email"`
}

type UsuarioCrudo struct {
	Nombre		string		`json:"user"`
	Pwd		string		`json:"password"`
	Email		string		`json:"email"`
}

type UsuarioCodigo struct {
	Codigo	string	`json:"codigo"`
}

type UsuarioRecuperar struct {
	Email		string		`json:"email"`
	Codigo		string		`json:"codigo"`
	Pwd				string		`json:"password"`
}
