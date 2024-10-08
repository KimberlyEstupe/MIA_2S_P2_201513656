package administradorpermisos

import (
	"MIA_2S_P2_201513656/Herramientas"
	"MIA_2S_P2_201513656/Structs"
	TI "MIA_2S_P2_201513656/ToolsInodos"
	"encoding/binary"
	"fmt"
	"path/filepath"
	"strings"
)

func Rename(entrada []string) string {
	respuesta := ""
	var path string	//ruta del archivo que cambiara el nombre
	var name string	//nombre que recibira el archivo/carpeta

	UsuarioA := Structs.UsuarioActual

	if !UsuarioA.Status {
		respuesta += "ERROR RENAME: NO HAY SECION INICIADA" + "\n"
		respuesta += "POR FAVOR INICIAR SESION PARA CONTINUAR" + "\n"
		return respuesta
	}

	for _, parametro := range entrada[1:] {
		tmp := strings.TrimRight(parametro," ")
		valores := strings.Split(tmp,"=")

		if len(valores)!=2{
			fmt.Println("ERROR RENAME, valor desconocido de parametros ",valores[1])
			respuesta += "ERROR RENAME, valor desconocido de parametros " + valores[1]
			//Si falta el valor del parametro actual lo reconoce como error e interrumpe el proceso
			return respuesta
		}

		//******************* PATH *************
		if strings.ToLower(valores[0]) == "path" {
			path = strings.ReplaceAll(valores[1],"\"","")	
		//********************  NAME *****************
		}else if strings.ToLower(valores[0]) == "name" {
			// Eliminar comillas
			name = strings.ReplaceAll(valores[1], "\"", "")
			// Eliminar espacios en blanco al final
			name = strings.TrimSpace(name)
		
		//******************* ERROR EN LOS PARAMETROS *************
		} else {
			fmt.Println("ERROR RENAME: Parametro desconocido: ", valores[0])
			respuesta += "ERROR RENAME: Parametro desconocido: "+ valores[0]
			return respuesta //por si en el camino reconoce algo invalido de una vez se sale
		}
	}

	if path==""{
		fmt.Println("ERROR RENAME FALTA PAREMETRO PATH")
		return "ERROR RENAME FALTA PAREMETRO PATH"
	}

	if name==""{
		fmt.Println("ERROR RENAME FALTA PAREMETRO NAME")
		return "ERROR RENAME FALTA PAREMETRO NAME"
	}

	//Abrimos el disco
	Disco, err := Herramientas.OpenFile(UsuarioA.PathD)
	if err != nil {
		return "CAR ERROR OPEN FILE " + err.Error() + "\n"
	}

	var mbr Structs.MBR
	// Read object from bin file
	if err := Herramientas.ReadObject(Disco, &mbr, 0); err != nil {
		return "CAR ERROR READ FILE " + err.Error() + "\n"
	}

	

	//Encontrar la particion correcta
	buscar := false
	part := -1 //particion a utilizar y modificar
	for i := 0; i < 4; i++ {
		identificador := Structs.GetId(string(mbr.Partitions[i].Id[:]))
		if identificador == UsuarioA.IdPart {
			part = i
			buscar = true
			break //para que ya no siga recorriendo si ya encontro la particion independientemente si se pudo o no reducir
		}
	}

	if buscar{
		//var contenido string
		//var fileBlock Structs.Fileblock
		var superBloque Structs.Superblock

		errREAD := Herramientas.ReadObject(Disco, &superBloque, int64(mbr.Partitions[part].Start))
		if errREAD != nil {
			fmt.Println("CAT ERROR. Particion sin formato")
			return "CAT ERROR. Particion sin formato" + "\n"
		}

		var inodo Structs.Inode
		var folderBlock Structs.Folderblock
		ArchivoEncontrado := true
		//divido el path entre el directorio y el nombre del archivo/carpeta que busco 
		//"carpeta/nombre"
		carpeta := filepath.Dir(path)
		tmp := strings.Split(path, "/")
		nombre := tmp[len(tmp)-1]

		fmt.Println("path ",carpeta," / ",nombre)
		idInodo := int32(0)
		//Si la carpeta es "/", significa que es la raiz por lo que idINodo es 0	
		if carpeta != "/"{
			idInodo = TI.BuscarInodo(0, carpeta, superBloque, Disco)
		}

		//Estoy en el inodo anterior al archivo/carpeta a cambiar el nombre
		//Si el inodo es diferente a -1, buscara el cada idBLoque que no este vacio
		Herramientas.ReadObject(Disco, &inodo, int64(superBloque.S_inode_start+(idInodo*int32(binary.Size(Structs.Inode{})))))
		//Verifico que el nuevo nombre del earchivo no xista
		rename := true
		for _, idBlock := range inodo.I_block {
			if idBlock != -1{
				Herramientas.ReadObject(Disco, &folderBlock, int64(superBloque.S_block_start+(idBlock*int32(binary.Size(Structs.Folderblock{})))))				
				for k := 2; k < 4; k++ {
					apuntador := folderBlock.B_content[k].B_inodo
					if apuntador != -1 {
						pathActual := Structs.GetB_name(string(folderBlock.B_content[k].B_name[:]))
						if pathActual == name{
							rename = false
							return "ERROR RENAME EL NOMBRE "+name+" YA EXISTE"
						}						
					}
				}
			}
		}
		
		//Verifico que el path exista
		for _, idBlock := range inodo.I_block {
			if idBlock != -1{
				Herramientas.ReadObject(Disco, &folderBlock, int64(superBloque.S_block_start+(idBlock*int32(binary.Size(Structs.Folderblock{})))))				
				//verifica que el archivo a cambiar si exista
				for k := 2; k < 4; k++ {
					apuntador := folderBlock.B_content[k].B_inodo
					if apuntador != -1 {
						pathActual := Structs.GetB_name(string(folderBlock.B_content[k].B_name[:]))
						if pathActual == nombre && rename{
							copy(folderBlock.B_content[k].B_name[:], name)
							ArchivoEncontrado = false
							//Escribir en el archivo los cambios del superBloque
							Herramientas.WriteObject(Disco, folderBlock, int64(superBloque.S_block_start+(idBlock*int32(binary.Size(Structs.Folderblock{})))))
						}						
					}
				}
			}
		}
		
		if ArchivoEncontrado {
			fmt.Println("Archivo/carpeta no existe")
			// Close bin file
			defer Disco.Close()
			return "ERROR RENAME: EL ARCHIVO O CARPETA EN PATH NO EXISTE"
		}else{
			fmt.Println("Archivo/carpeta no existe")
			// Close bin file
			defer Disco.Close()
			return "El nombre del archivo "+name+" fue modificado con exito "			
		}
	}

	
	return respuesta
}