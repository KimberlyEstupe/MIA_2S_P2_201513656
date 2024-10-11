package administradorpermisos

import (
	"MIA_2S_P2_201513656/Herramientas"
	"MIA_2S_P2_201513656/Structs"
	ToolsInodos "MIA_2S_P2_201513656/ToolsInodos"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

func Edit(entrada []string) string {
	respuesta := ""
	var path string			//ruta del archivo
	var contenido string	//nuevo contenido del archivo

	UsuarioA := Structs.UsuarioActual

	if !UsuarioA.Status {
		respuesta += "ERROR EDIT: NO HAY SECION INICIADA" + "\n"
		respuesta += "POR FAVOR INICIAR SESION PARA CONTINUAR" + "\n"
		return respuesta
	}

	for _, parametro := range entrada[1:]{
		tmp := strings.TrimRight(parametro," ")
		valores := strings.Split(tmp,"=")

		if len(valores)!=2{
			fmt.Println("ERROR EDIT, valor desconocido de parametros ",valores[1])
			respuesta += "ERROR EDIT, valor desconocido de parametros " + valores[1]
			//Si falta el valor del parametro actual lo reconoce como error e interrumpe el proceso
			return respuesta
		}

		//******************* PATH *************
		if strings.ToLower(valores[0]) == "path" {
			path = strings.ReplaceAll(valores[1],"\"","")	
		//********************  CONTENIDO *****************
		}else if strings.ToLower(valores[0]) == "contenido" {
			// Eliminar comillas
			contenido = strings.ReplaceAll(valores[1], "\"", "")
			_, err := os.Stat(contenido)
				if os.IsNotExist(err) {
					fmt.Println("MKFILE Error: El archivo cont no existe")
					respuesta +=  "MKFILE Error: El archivo cont no existe"+ "\n"
					return respuesta // Terminar el bucle porque encontramos un nombre único
				}
		
		//******************* ERROR EN LOS PARAMETROS *************
		} else {
			fmt.Println("ERROR EDIT: Parametro desconocido: ", valores[0])
			respuesta += "ERROR EDIT: Parametro desconocido: "+ valores[0]
			return respuesta //por si en el camino reconoce algo invalido de una vez se sale
		}
	}

	if path==""{
		fmt.Println("ERROR EDIT FALTA PAREMETRO PATH")
		return "ERROR EDIT FALTA PAREMETRO PATH"
	}

	if contenido==""{
		fmt.Println("ERROR EDIT FALTA PAREMETRO CONTENIDO")
		return "ERROR EDIT FALTA PAREMETRO CONTENIDO"
	}

	//Abrimos el disco
	Disco, err := Herramientas.OpenFile(UsuarioA.PathD)
	if err != nil {
		return "EDIT ERROR OPEN FILE " + err.Error() + "\n"
	}

	var mbr Structs.MBR
	// Read object from bin file
	if err := Herramientas.ReadObject(Disco, &mbr, 0); err != nil {
		return "EDIT ERROR READ FILE " + err.Error() + "\n"
	}
	
	//Encontrar la particion correcta
	editar := false
	part := -1 //particion a utilizar y modificar
	for i := 0; i < 4; i++ {
		identificador := Structs.GetId(string(mbr.Partitions[i].Id[:]))
		if identificador == UsuarioA.IdPart {
			part = i
			editar = true
			break //para que ya no siga recorriendo si ya encontro la particion independientemente si se pudo o no reducir
		}
	}

	if editar{
		var fileBlock Structs.Fileblock
		var superBloque Structs.Superblock

		errREAD := Herramientas.ReadObject(Disco, &superBloque, int64(mbr.Partitions[part].Start))
		if errREAD != nil {
			fmt.Println("EDIT ERROR. Particion sin formato")
			return "EDIT ERROR. Particion sin formato" + "\n"
		}

		//buscar el inodo que contiene el archivo buscado
		idInodo := ToolsInodos.BuscarInodo(0, path, superBloque, Disco)
		var inodo Structs.Inode
		Herramientas.ReadObject(Disco, &inodo, int64(superBloque.S_inode_start+(idInodo*int32(binary.Size(Structs.Inode{})))))

		//Verifica que el usuario logiado sea root(root tiene todos los permisos) o que sea el propietario del archivo
		if inodo.I_uid == UsuarioA.IdUsr || UsuarioA.Nombre=="root"{
			var oldContenido string					
			//recorrer los fileblocks del inodo para obtener toda su informacion
			for _, idBlock := range inodo.I_block {
				if idBlock != -1 {
					Herramientas.ReadObject(Disco, &fileBlock, int64(superBloque.S_block_start+(idBlock*int32(binary.Size(Structs.Fileblock{})))))
					oldContenido = string(fileBlock.B_content[:])										
				}
			}
			Editar(contenido, oldContenido, len(contenido), len(oldContenido), idInodo, int64(mbr.Partitions[part].Start), Disco)
			respuesta += "\n"
		}else{
			respuesta += "ERROR EDIT: No tiene permisos para visualizar el archivo \n"
		}
		
	}
	return respuesta
}

func Editar(NewCont string, OldCont string, NewSize int, OldSize int, idInodo int32, initSuperBloque int64, disco *os.File) string{
	var respuesta string
	return respuesta
}