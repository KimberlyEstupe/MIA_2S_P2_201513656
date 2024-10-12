package main

import (
	AD "MIA_2S_P2_201513656/Comandos/AdministradorDiscos"
	SA 	"MIA_2S_P2_201513656/Comandos/SistemaArchivos"
	USR  "MIA_2S_P2_201513656/Comandos/Usuario"
	REP "MIA_2S_P2_201513656/Comandos/Rep"
	AP	"MIA_2S_P2_201513656/Comandos/AdministradorPermisos"
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/rs/cors"
)

type Entrada struct {
	Text string `json:"text"`
}

type StatusResponse struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

func main()  {
	//EndPoint 
	http.HandleFunc("/analizar", getCadenaAnalizar)

	// Configurar CORS con opciones predeterminadas
	//Permisos para enviar y recir informacion
	c := cors.Default()

	// Configurar el manejador HTTP con CORS
	handler := c.Handler(http.DefaultServeMux)

	// Iniciar el servidor en el puerto 8080
	fmt.Println("Servidor escuchando en http://localhost:8080")
	http.ListenAndServe(":8080", handler)
}

func getCadenaAnalizar(w http.ResponseWriter, r *http.Request) {
	var respuesta string
	// Configurar la cabecera de respuesta
	w.Header().Set("Content-Type", "application/json")
	
	var status StatusResponse
	//verificar que sea un post
	if r.Method == http.MethodPost {
		var entrada Entrada
		if err := json.NewDecoder(r.Body).Decode(&entrada); err != nil {
			http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
			status = StatusResponse{Message: "Error al decodificar JSON", Type: "unsucces"}
			json.NewEncoder(w).Encode(status)
			return
		}
		
		//creo un lector de bufer para el archivo
		lector := bufio.NewScanner(strings.NewReader(entrada.Text))
		//leer el archivo linea por linea
		for lector.Scan() {
			//Elimina los saltos de linea
			if lector.Text() != ""{
				//Divido por # para ignorar todo lo que este a la derecha del mismo
				linea := strings.Split(lector.Text(), "#") //lector.Text() retorna la linea leida
				if len(linea[0]) != 0 {
					fmt.Println("\n*********************************************************************************************")
					fmt.Println("Comando en ejecucion: ", linea[0])
					respuesta += "***************************************************************************************************************************\n"
					respuesta += "Comando en ejecucion: " + linea[0] + "\n"
					respuesta += Analizar(linea[0])  + "\n"
				}	
				//Comentarios			
				if len(linea) > 1 && linea[1] != "" {
					fmt.Println("#"+linea[1] +"\n")
					respuesta += "#"+linea[1] +"\n"
				}
			}
			
		}

		//fmt.Println("Cadena recibida ", entrada.Text)
		w.WriteHeader(http.StatusOK)

		status = StatusResponse{Message: respuesta, Type: "succes"}
		json.NewEncoder(w).Encode(status)

	} else {
		//http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		status = StatusResponse{Message: "Metodo no permitido", Type: "unsucces"}
		json.NewEncoder(w).Encode(status)
	}
}

func Analizar(entrada string)string{
	tmp := strings.TrimRight(entrada," ")
	//Recibe una linea y la descompone entre el comando y sus parametros
	parametros:= strings.Split(tmp, " -")

	// *============================* ADMINISTRACION DE DISCOS *============================*
	//mkdisk -size=5 -unit=M -path=Calificacion_MIA/Discos/Disco1.mia
	//mkdisk -size=5 -unit=M -path="Calificacion_MIA/Discos/Disco_1.mia"
	if strings.ToLower(parametros[0])=="mkdisk"{
		if len(parametros)>1{	
			AD.Mkdisk(parametros)				
			return AD.Mkdisk(parametros)
		}else{
			fmt.Println("ERROR EN MKDISK, FALTAN PARAMETROS EN MKDISK")
			return "ERROR EN MKDISK, FALTAN PARAMETROS EN MKDISK"
		}

	//rmdisk -path="/home/mis discos/Disco4.mia"
	}else if strings.ToLower(parametros[0])=="rmdisk"{
		if len(parametros)>1{	
			return AD.Rmdisk(parametros)		
		}else{
			fmt.Println("ERROR EN RMDISK, FALTAN PARAMETROS EN MKDISK")
			return  "ERROR EN RMDISK, FALTAN PARAMETROS EN MKDISK"
		}

	//fdisk -type=P -unit=b -name=Part1 -size=10485760 -path=Calificacion_MIA/Discos/Disco1.mia
	//fdisk -add=-23760 -path=Calificacion_MIA/Discos/Disco1.mia -name=Part1 -size=10485760
	}else if strings.ToLower(parametros[0])=="fdisk"{
		if len(parametros)>1{	
			return AD.Fdisk(parametros)		
		}else{
			fmt.Println("ERROR EN FDISK, FALTAN PARAMETROS EN MKDISK")
			return  "ERROR EN FDISK, FALTAN PARAMETROS EN MKDISK"
		}

	//mount -path=/home/Disco3.mia -name=Part1
	}else if strings.ToLower(parametros[0])=="mount"{
		if len(parametros)>1{	
			return AD.Mount(parametros)			
		}else{
			fmt.Println("ERROR EN MKDISK, FALTAN PARAMETROS EN MKDISK")
			return "ERROR EN MKDISK, FALTAN PARAMETROS EN MKDISK"
		}

	//unmount -id=561A
	}else if strings.ToLower(parametros[0])=="unmount"{
		if len(parametros)>1{	
			return AD.Unmoun(parametros)		
		}else{
			fmt.Println("ERROR EN MKDISK, FALTAN PARAMETROS EN MKDISK")
			return  "ERROR EN MKDISK, FALTAN PARAMETROS EN MKDISK"
		}
	// *===================* ADMINISTRACION DE SISTEMA DE ARCHIVOS *======================*
	//ejm: mkfs -type=full -id=341A	-fs=3fs
	}else if strings.ToLower(parametros[0])=="mkfs"{		
		if len(parametros)>1{			
			return SA.MKfs(parametros)
		}else{
			fmt.Println("ERROR EN MKFS, FALTAN PARAMETROS")
			return "ERROR EN MKFS, FALTAN PARAMETROS"
		}
	
	// *===================* ADMINISTRACION DE USUARIOS Y CARPETAS *======================*
	//login -user=root -pass=123 -id=561A
	}else if strings.ToLower(parametros[0])=="login"{		
		if len(parametros)>1{			
			return USR.Login(parametros)
		}else{
			fmt.Println("ERROR EN LOGIN, FALTAN PARAMETROS")
			return "ERROR EN LOGIN, FALTAN PARAMETROS"
		}

	//logout
	}else if strings.ToLower(parametros[0])=="logout"{		
		return USR.Logout()
	
	//mkgrp -name=usuarios
	}else if strings.ToLower(parametros[0])=="mkgrp"{		
		if len(parametros)>1{			
			return USR.Mkgrp(parametros)
		}else{
			fmt.Println("ERROR EN MKGRP, FALTAN PARAMETROS")
			return "ERROR EN MKGRP, FALTAN PARAMETROS"
		}
		
	//rmgrp -name=usuarios
	}else if strings.ToLower(parametros[0])=="rmgrp"{		
		if len(parametros)>1{			
			return USR.Rmgrp(parametros)
		}else{
			fmt.Println("ERROR EN RMGRP, FALTAN PARAMETROS")
			return "ERROR EN RMGRP, FALTAN PARAMETROS"
		}
	
	}else if strings.ToLower(parametros[0])=="mkusr"{		
		if len(parametros)>1{			
			return USR.Mkusr(parametros)
		}else{
			fmt.Println("ERROR EN RMGRP, FALTAN PARAMETROS")
			return "ERROR EN RMGRP, FALTAN PARAMETROS"
		}	

	}else if strings.ToLower(parametros[0])=="rmusr"{		
		if len(parametros)>1{			
			return USR.Rmusr(parametros)
		}else{
			fmt.Println("ERROR EN RMUSR, FALTAN PARAMETROS")
			return "ERROR EN RMUSR, FALTAN PARAMETROS"
		}

	}else if strings.ToLower(parametros[0])=="chgrp"{		
		if len(parametros)>1{			
			return USR.Chgrp(parametros)
		}else{
			fmt.Println("ERROR EN CHGRP, FALTAN PARAMETROS")
			return "ERROR EN CHGRP, FALTAN PARAMETROS"
		}
	
	// *=======================* PERMISOS DE CARPETAS Y ARHICVOS *============================*
	//mkfile -path=/home/archivos/docs/Tarea2.txt -size=75 -r
	}else if strings.ToLower(parametros[0])=="mkfile"{		
		if len(parametros)>1{			
			return AP.MKfile(parametros)
		}else{
			fmt.Println("ERROR EN MKFILE, FALTAN PARAMETROS")
			return "ERROR EN MKFILE, FALTAN PARAMETROS"
		}
	
	//EJ: cat -file1=/home/user/docs/a.txt -file12=/home/user/docs/b.txt
	}else if strings.ToLower(parametros[0])=="cat"{		
		if len(parametros)>1{			
			return AP.Cat(parametros)
		}else{
			fmt.Println("ERROR EN CAT, FALTAN PARAMETROS")
			return "ERROR EN CAT, FALTAN PARAMETROS"
		}

	//mkdir -r -path=/home/archivos/Fotos
	}else if strings.ToLower(parametros[0])=="mkdir"{		
		if len(parametros)>1{			
			return AP.Mkdir(parametros)
		}else{
			fmt.Println("ERROR EN MKDIR, FALTAN PARAMETROS")
			return "ERROR EN MKDIR, FALTAN PARAMETROS"
		}
	
	}else if strings.ToLower(parametros[0])=="rename"{		
		if len(parametros)>1{			
			return AP.Rename(parametros)
		}else{
			fmt.Println("ERROR EN RENAME, FALTAN PARAMETROS")
			return "ERROR EN RENAME, FALTAN PARAMETROS"
		}
	
	}else if strings.ToLower(parametros[0])=="edit"{		
		if len(parametros)>1{			
			return AP.Edit(parametros)
		}else{
			fmt.Println("ERROR EN EDIT, FALTAN PARAMETROS")
			return "ERROR EN EDIT, FALTAN PARAMETROS"
		}
	
	}else if strings.ToLower(parametros[0])=="copy"{		
		if len(parametros)>1{			
			return AP.Copy(parametros)
		}else{
			fmt.Println("ERROR EN COPY, FALTAN PARAMETROS")
			return "ERROR EN COPY, FALTAN PARAMETROS"
		}
	// *============================* OTROS *============================*
	} else if strings.ToLower(parametros[0]) == "rep" {
		//REP
		if len(parametros) > 1 {
			return REP.Rep(parametros)
		} else {
			fmt.Println("REP ERROR: parametros no encontrados")
			return "REP ERROR: parametros no encontrados"
		}
	} else if strings.ToLower(parametros[0]) == "" {
		//para agregar lineas con cada enter sin tomarlo como error
		return ""
	} else {
		fmt.Println("ERROR: COMANDO "+parametros[0]+" NO RECONOCIBLE")
		return "ERROR: COMANDO "+parametros[0]+" NO RECONOCIBLE"
	}
}