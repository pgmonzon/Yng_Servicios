package core

// NOTA: Hay objetos de bson guardados en string. Tal vez es preferible ser consistente y guardar todos como objetos de bson
import (
    "net/http"
    "log"
    "errors"

    "github.com/pgmonzon/Yng_Servicios/models"
    "github.com/pgmonzon/Yng_Servicios/cfg"
    //"gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

func ChequearPermisos(r *http.Request, permisoBuscado string) (bool) {
  // esta funcion se encarga de responder SI o NO a la pregunta "¿tiene este usuario permisos para ejecutar lo que me esta pidiendo?"
  id := extraerClaim(r, "id")
  permiso, err := extraerInfoPermiso(permisoBuscado)
  if (err == nil) { return false }
  if (!permiso.Activo || permiso.Borrado){
    return false
  }
  if id == ""{
    return false  //Esto sería un error más que falta de permisos (no existe el campo id en el token o es un token invalido). Hay que buscar la forma de manejar estos errores
  }
  user, err := extraerInfoUsuario(id.(string)) // La tengo que convertir a string porque me devolvieron una interface{}
  if (err != nil) { return false }
  if user.Rol == cfg.GuestRol{
    return false //Es guest
  }
  a := extraerPermisosDelRol(user.Rol)
  for _, v := range a.IDPermisos {
    v_bson := bson.ObjectIdHex(v)
    if (v_bson == permiso.ID) {
      log.Println("Acceso permitido de:",user.User,"a:",permisoBuscado)
      return true
    }
  }
  //log.Println("Rol del usuario: ",extraerInfoRol(user.Rol, session).Nombre)
  return false
}

func extraerInfoRol(id string) (models.Roles, error) {
  var modelRol []models.Roles
  var modelError models.Roles
  if bson.IsObjectIdHex(id) != true { // Un poco de sanity.
    log.Println("ERROR: Id rol invalida.", id)
    return modelError, errors.New("La id del Rol no es un objeto bson")
  }
  id_bson := bson.ObjectIdHex(id)
  session := Session.Copy()
  defer session.Close()
  collection := session.DB(Dbname).C("roles")
  collection.Find(bson.M{"_id": id_bson}).All(&modelRol)
  if (len(modelRol) == 0) {
    log.Println("ERROR: Id invalida. El usuario", id, "tiene un Rol que no existe")
    return modelError, errors.New("El usuario tiene un rol que no existe")
  }
  return modelRol[0], nil
}

func extraerInfoUsuario(id string) (models.Usuario, error) {
  //que rol tiene la id que nos pasan???
  var usuario []models.Usuario
  var modelError models.Usuario
	if bson.IsObjectIdHex(id) != true { // Un poco de sanity.
    log.Printf("FATAL ERROR: Id usuario invalida.")
		return modelError, errors.New("Id usuario invalida") //Podria devolver la ID de un usuario especial (una especie de muñeco sin permisos)
	}
  id_bson := bson.ObjectIdHex(id)
	session := Session.Copy()
	defer session.Close()
	collection := session.DB(Dbname).C("usuarios")
  collection.Find(bson.M{"_id": id_bson}).All(&usuario)
	if (len(usuario) == 0) {
		log.Printf("ERROR: Id invalida. Lo cual significa que /login esta creando tokens con IDs rotas o la db esta corrupta")
		return modelError, errors.New("Id usuario invalida")
	}
	return usuario[0], nil
}

func extraerPermisosDelRol(id string) (models.RP){
  //le das una ID de rol a esta funcion, y te devuelve los permisos que tiene ese Rol (los devuelve en un array)
  var rp []models.RP
  if bson.IsObjectIdHex(id) != true { // Un poco de sanity.
    log.Printf("FATAL ERROR: Id rol invalida.")
    return rp[0]
  }
  session := Session.Copy()
  defer session.Close()
  collection := session.DB(Dbname).C("rp")
  err := collection.Find(bson.M{"idrol": id}).All(&rp)
  if err != nil {
    log.Printf("FATAL ERROR: Id invalida. Lo cual significa que /login esta creando tokens con IDs rotas")
    return rp[0]
  }
  return rp[0] //esto no es ideal, es temporal
}

func extraerInfoPermiso(permiso string) (models.Permisos, error) {
  //Nota: En caso que sea necesario, se puede hacer un case switch si "permiso string" es una ID o el nombre del permiso
  var modelPermisos []models.Permisos
  var modelError models.Permisos
  session := Session.Copy()
  defer session.Close()
  collection := session.DB(Dbname).C("permisos")
  collection.Find(bson.M{"nombre": permiso}).All(&modelPermisos)
  if (len(modelPermisos) == 0) {
    log.Printf("ERROR: Permiso invalido. El permiso buscado no existe.")
    return modelError, errors.New("El permiso buscado no existe")
  }
  return modelPermisos[0], nil
}

/*func EstaActivo(permiso string) (bool){
  if (extraerInfoPermiso(permiso).Activo == true) { return true }
  return false
}*/
