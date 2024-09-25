package main

import (
	AD "MIA_2S_P2_201513656/Comandos/AdministradorDiscos"
	"fmt"
	"strings"
)

func main()  {
	fmt.Println("Hola")
}

func Analizar(entrada string){
	tmp := strings.TrimRight(entrada," ")
	//Recibe una linea y la descompone entre el comando y sus parametros
	parametros:= strings.Split(tmp, " -")

	// *============================* ADMINISTRACION DE DISCOS *============================*
	if strings.ToLower(parametros[0])=="mkdisk"{
		if len(parametros)>1{	
			AD.Mkdisk(parametros)				
			//respuesta = AD.Mkdisk(parametros)
		}else{
			fmt.Println("ERROR EN MKDISK, FALTAN PARAMETROS EN MKDISK")
			//respuesta = "ERROR EN MKDISK, FALTAN PARAMETROS EN MKDISK"
		}
	}else if strings.ToLower(parametros[0])=="rmdisk"{
		if len(parametros)>1{	
			AD.Rmdisk(parametros)				
			//respuesta = AD.Mkdisk(parametros)
		}else{
			fmt.Println("ERROR EN MKDISK, FALTAN PARAMETROS EN MKDISK")
			//respuesta = "ERROR EN MKDISK, FALTAN PARAMETROS EN MKDISK"
		}
	}else if strings.ToLower(parametros[0])=="fdisk"{
		if len(parametros)>1{	
			AD.Rmdisk(parametros)				
			//respuesta = AD.Mkdisk(parametros)
		}else{
			fmt.Println("ERROR EN MKDISK, FALTAN PARAMETROS EN MKDISK")
			//respuesta = "ERROR EN MKDISK, FALTAN PARAMETROS EN MKDISK"
		}
	// *============================* OTROS *============================*
	} else if strings.ToLower(parametros[0]) == "" {
		//para agregar lineas con cada enter sin tomarlo como error
		return 
	} else {
		fmt.Println("Comando no reconocible")
		//return "ERROR: COMANDO NO RECONOCIBLE"
	}
}