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

func Copy(entrada []string) string {
	respuesta := ""
	var path string			//ruta del archivo
	var destino string	//nuevo contenido del archivo

	UsuarioA := Structs.UsuarioActual

	if !UsuarioA.Status {
		respuesta += "ERROR COPY: NO HAY SECION INICIADA" + "\n"
		respuesta += "POR FAVOR INICIAR SESION PARA CONTINUAR" + "\n"
		return respuesta
	}

	for _, parametro := range entrada[1:]{
		tmp := strings.TrimRight(parametro," ")
		valores := strings.Split(tmp,"=")

		if len(valores)!=2{
			fmt.Println("ERROR COPY, valor desconocido de parametros ",valores[1])
			respuesta += "ERROR COPY, valor desconocido de parametros " + valores[1]
			//Si falta el valor del parametro actual lo reconoce como error e interrumpe el proceso
			return respuesta
		}

		//******************* PATH *************
		if strings.ToLower(valores[0]) == "path" {
			path = strings.ReplaceAll(valores[1],"\"","")	
		//********************  CONTENIDO *****************
		}else if strings.ToLower(valores[0]) == "destino" {
			// Eliminar comillas
			destino = strings.ReplaceAll(valores[1], "\"", "")
		
		//******************* ERROR EN LOS PARAMETROS *************
		} else {
			fmt.Println("ERROR COPY: Parametro desconocido: ", valores[0])
			respuesta += "ERROR COPY: Parametro desconocido: "+ valores[0]
			return respuesta //por si en el camino reconoce algo invalido de una vez se sale
		}
	}

	if path==""{
		fmt.Println("ERROR COPY FALTA PAREMETRO PATH")
		return "ERROR COPY FALTA PAREMETRO PATH"
	}

	if destino==""{
		fmt.Println("ERROR COPY FALTA PAREMETRO NAME")
		return "ERROR COPY FALTA PAREMETRO NAME"
	}

	//Abrimos el disco
	Disco, err := Herramientas.OpenFile(UsuarioA.PathD)
	if err != nil {
		return "COPY ERROR OPEN FILE " + err.Error() + "\n"
	}

	var mbr Structs.MBR
	// Read object from bin file
	if err := Herramientas.ReadObject(Disco, &mbr, 0); err != nil {
		return "COPY ERROR READ FILE " + err.Error() + "\n"
	}
	
	//Encontrar la particion correcta
	copy := false
	part := -1 //particion a utilizar y modificar
	for i := 0; i < 4; i++ {
		identificador := Structs.GetId(string(mbr.Partitions[i].Id[:]))
		if identificador == UsuarioA.IdPart {
			part = i
			copy = true
			break //para que ya no siga recorriendo si ya encontro la particion independientemente si se pudo o no reducir
		}
	}

	if copy{
		var superBloque Structs.Superblock

		errREAD := Herramientas.ReadObject(Disco, &superBloque, int64(mbr.Partitions[part].Start))
		if errREAD != nil {
			fmt.Println("COPY ERROR. Particion sin formato")
			return "COPY ERROR. Particion sin formato" + "\n"
		}

		//buscar el inodo de la carpeta destino
		idInodoDestino := ToolsInodos.BuscarInodo(0, destino, superBloque, Disco)
		var inodoDestino Structs.Inode
		Herramientas.ReadObject(Disco, &inodoDestino, int64(superBloque.S_inode_start+(idInodoDestino*int32(binary.Size(Structs.Inode{})))))

		//Verifica que el usuario logiado sea root(root tiene todos los permisos) o que sea el propietario del archivo
		if inodoDestino.I_uid == UsuarioA.IdUsr || UsuarioA.Nombre=="root" || inodoDestino.I_gid == UsuarioA.IdGrp{					
			//buscar el inodo de la carpeta a copiar
			idNewInodo := ToolsInodos.BuscarInodo(0, path, superBloque, Disco)
			var NewInodo Structs.Inode
			Herramientas.ReadObject(Disco, &NewInodo, int64(superBloque.S_inode_start+(idNewInodo*int32(binary.Size(Structs.Inode{})))))

			var fileBlock Structs.Fileblock
			for _, idBlock := range NewInodo.I_block {
				if idBlock != -1 {
					Herramientas.ReadObject(Disco, &fileBlock, int64(superBloque.S_block_start+(idBlock*int32(binary.Size(Structs.Fileblock{})))))
					tmpConvertir := Herramientas.EliminartIlegibles(string(fileBlock.B_content[:]))
					fmt.Println(idBlock," ",tmpConvertir)				
				}
			}

		}else{
			respuesta += "ERROR COPY: No tiene permisos para copiar archivos a esta carpeta \n"
		}
		
	}
	return respuesta
}

func Copiar(idNodoDestino int32, idNodoCopiar int32, initSuperBloque int64, disco *os.File){}