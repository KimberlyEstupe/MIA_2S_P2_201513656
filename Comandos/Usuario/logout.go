package usuario

import (
	"MIA_2S_P2_201513656/Structs"
	"fmt"
)

func Logout() string{
	var respuesta string
	if Structs.UsuarioActual.Status {
		Structs.SalirUsuario()
		fmt.Println("Se ha cerrado la sesion")
		respuesta += "Se ha cerrado la sesion"
	}else{
		respuesta += "ERROR LOGUT: NO HAY SECION INICIADA"
	}

	return respuesta
}