package rep

import (
	toolsinodos "MIA_2S_P2_201513656/ToolsInodos"
	"MIA_2S_P2_201513656/Herramientas"
	"MIA_2S_P2_201513656/Structs"	
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Rep(entrada []string) string{
	var respuesta string
	var name string //obligatorio Nombre del reporte a generar
	var path string //obligatorio Nombre que tendrá el reporte
	var id string   //obligatorio sera el del disco o el de la particion
	var rutaFile string	//nombre del archivo o carpeta reporte file/IS
	Valido := true 

	for _, parametro := range entrada[1:]{
		tmp := strings.TrimRight(parametro," ")
		valores := strings.Split(tmp,"=")

		if len(valores)!=2{
			fmt.Println("ERROR REP, valor desconocido de parametros ",valores[1])
			//Si falta el valor del parametro actual lo reconoce como error e interrumpe el proceso
			return "ERROR REP, valor desconocido de parametros "+valores[1]
		}

		if strings.ToLower(valores[0]) == "name" {
			name = strings.ToLower(valores[1])
		} else if strings.ToLower(valores[0]) == "path" {
			path = strings.ReplaceAll(valores[1], "\"", "")
		} else if strings.ToLower(valores[0]) == "id" {
			id = strings.ToUpper(valores[1])
		} else if strings.ToLower(valores[0]) == "path_file_ls" {
			rutaFile = strings.ReplaceAll(valores[1], "\"", "")
		} else {
			fmt.Println("REP Error: Parametro desconocido: ", valores[0])
			respuesta+="REP Error: Parametro desconocido: " + valores[0]
			Valido = false
			break //por si en el camino reconoce algo invalido de una vez se sale
		}
	}

	if Valido{
		if name != "" && id != "" && path != "" {			
			switch name{
			case "mbr"://rep -id=561A -path=Calificacion_MIA/reports/reporte1.jpg -name=mbr
				fmt.Println("reporte mbr")
				respuesta+= Rmbr(path, id)
			case "disk"://rep -id=561A -path=Calificacion_MIA/reports/report2.pdf -name=disk
				fmt.Println("reporte disk")
				respuesta+= disk(path, id)
			case "bm_inode"://rep -id=561A -path=Calificacion_MIA/reports/report5.txt -name=bm_inode
				fmt.Println("reporte bm_inode")
				respuesta += BM_inode(path, id)
			case "bm_block"://rep -id=561A -path=Calificacion_MIA/reports/report6.txt -name=bm_block
				fmt.Println("reporte bm_block")
				respuesta += BM_Bloque(path, id)
			case "sb"://rep -id=561A -path=Calificacion_MIA/reports/report8.jpg -name=sb
				fmt.Println("reporte sb")
				respuesta += superBloque(path, id)
			case "file"://rep -id=561A -path=Calificacion_MIA/reports/report9.txt -path_file_ls=/users.txt -name=file
				fmt.Println("reporte file")
				respuesta += FILE(path, id, rutaFile)
			case "ls":// rep -id=561A -path=Calificacion_MIA/reports/report10.jpg -path_file_ls=/ -name=ls
				respuesta += LS(path, id, rutaFile)
				fmt.Println("reporte ls")
			case "journal"://rep -id=561A -path=Calificacion_MIA/reports/reportJournal.pdf -name=journal
				journal(path, id)
				fmt.Println("reporte journal")
			case "tree"://rep -id=561A -path=Calificacion_MIA/reports/tree.pdf -name=tree
				tree(path,id)
				fmt.Println("reporte tree")
			default:
				fmt.Println("REP Error: Reporte ", name, " desconocido")
				respuesta+="REP Error: Reporte "+ name+" desconocido"
			}
		}else{
			fmt.Println("REP Error: Faltan parametros")
			respuesta+= "REP Error: Faltan parametros"
		}
	}
	return respuesta
}

// =============================== MBR ===============================
func Rmbr (path string, id string) string{
	var Respuesta string
	var pathDico string
	Valido := false

	//BUsca en struck de particiones montadas el id ingresado
	for _,montado := range Structs.Montadas{
		if montado.Id == id{
			pathDico = montado.PathM
			Valido = true
		}
	}

	if Valido{
		tmp := strings.Split(path, "/")
		nombre := strings.Split(tmp[len(tmp)-1], ".")[0]
		
		tmp = strings.Split(pathDico, "/")
		NOmbreDis := strings.Split(tmp[len(tmp)-1], ".")[0]
		
		file, err := Herramientas.OpenFile(pathDico)
		if err != nil {
			Respuesta += "ERROR REP MBR Open "+ err.Error()		
		}

		var mbr Structs.MBR
		// Read object from bin file
		if err := Herramientas.ReadObject(file, &mbr, 0); err != nil {
			Respuesta += "ERROR REP MBR Read "+ err.Error()		
		}

		// Close bin file
		defer file.Close()

		//Crea reporte
		cad := "digraph { \nnode [ shape=none ] \nTablaReportNodo [ label = < <table border=\"1\"> \n"
		cad += " <tr>\n  <td bgcolor='SlateBlue' COLSPAN=\"2\"> Reporte MBR </td> \n </tr> \n"
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> mbr_tamano </td> \n  <td bgcolor='Azure'> %d </td> \n </tr> \n", mbr.MbrSize)
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='#AFA1D1'> mbr_fecha_creacion </td> \n  <td bgcolor='#AFA1D1'> %s </td> \n </tr> \n", string(mbr.FechaC[:]))
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> mbr_disk_signature </td> \n  <td bgcolor='Azure'> %d </td> \n </tr>  \n", mbr.Id)
		cad += Structs.RepGraphviz(mbr, file)
		cad += "</table> > ]\n}"

		carpeta := filepath.Dir(path)
		rutaReporte := carpeta + "/" + nombre + ".dot"

		Herramientas.RepGraphizMBR(rutaReporte, cad, nombre)
		Respuesta += "Reporte de MBR del disco "+NOmbreDis+" creado con el nombre "+nombre+".png"
	}else{
		Respuesta += "ERROR: EL ID INGRESADO NO EXISTE"
	}

	
	return Respuesta
}

//=============================== DISK ===============================
func disk(path string, id string)string{
	var Respuesta string
	var pathDico string
	Valido := false

	//BUsca en struck de particiones montadas el id ingresado
	for _,montado := range Structs.Montadas{
		if montado.Id == id{
			pathDico = montado.PathM
			Valido = true
		}
	}

	if Valido{
		tmp := strings.Split(path, "/")
		nombre := strings.Split(tmp[len(tmp)-1], ".")[0]
		
		tmp = strings.Split(pathDico, "/")
		NOmbreDis := strings.Split(tmp[len(tmp)-1], ".")[0]
		
		file, err := Herramientas.OpenFile(pathDico)
		if err != nil {
			Respuesta += "ERROR REP DISK Open "+ err.Error()	
			return Respuesta	
		}

		var TempMBR Structs.MBR
		// Read object from bin file
		if err := Herramientas.ReadObject(file, &TempMBR, 0); err != nil {
			Respuesta += "ERROR REP READ Open "+ err.Error()
			return Respuesta	
		}

		defer file.Close()

		//inicia contenido del reporte graphviz del disco
		cad := "digraph { \nnode [ shape=none ] \nTablaReportNodo [ label = < <table border=\"1\"> \n<tr> \n"
		cad += " <td bgcolor='SlateBlue'  ROWSPAN='3'> MBR </td>\n"
		cad += Structs.RepDiskGraphviz(TempMBR, file)
		cad += "\n</table> > ]\n}"

		carpeta := filepath.Dir(path)
		rutaReporte := carpeta + "/" + nombre + ".dot"

		fmt.Println("RP ", rutaReporte," name ",nombre)

		Herramientas.RepGraphizMBR(rutaReporte, cad, nombre)
		Respuesta += "Reporte Disk del disco "+NOmbreDis+" creado con el nombre "+nombre+".png"
	}else{
		Respuesta += "ERROR: EL ID INGRESADO NO EXISTE"
	}
	
	return Respuesta
}

// =============================== SB ===============================
func superBloque (path string, id string) string{
	var respuesta string
	var pathDico string
	reportar := false

	//BUsca en struck de particiones montadas el id ingresado
	for _,montado := range Structs.Montadas{
		if montado.Id == id{
			pathDico = montado.PathM
			reportar = true
		}
	}

	//Verifica que se encontro el ID y la Path del disco
	if pathDico == ""{
		reportar = false
		return "ERROR REP: ID NO ENCONTRADO"
	}

	if reportar{
		tmp := strings.Split(path, "/")
		nombre := strings.Split(tmp[len(tmp)-1], ".")[0]

		tmp2 := strings.Split(pathDico, "/")
		nombreDisco := strings.Split(tmp2[len(tmp2)-1], ".")[0]

		file, err := Herramientas.OpenFile(pathDico)
		if err != nil {
			return "ERROR REP SB OPEN FILE "+err.Error()
		}

		var mbr Structs.MBR
		// Read object from bin file
		if err := Herramientas.ReadObject(file, &mbr, 0); err != nil {
			return "ERROR REP SB READ FILE "+err.Error()
		}

		// Close bin file
		defer file.Close()

		//Encontrar la particion correcta
		part := -1
		for i := 0; i < 4; i++ {
			identificador := Structs.GetId(string(mbr.Partitions[i].Id[:]))
			if identificador == id {
				reportar = true
				part = i
				break //para que ya no siga recorriendo si ya encontro la particion independientemente si se pudo o no reducir
			}
		}

		cad := "digraph { \nnode [ shape=none ] \nTablaReportNodo [ label = < <table border=\"1\"> \n"
		cad += " <tr>\n  <td bgcolor='darkgreen' COLSPAN=\"2\"> <font color='white'> Reporte SUPERBLOQUE </font> </td> \n </tr> \n"
		cad += Structs.RepSB(mbr.Partitions[part], file)
		cad += "</table> > ]\n}"

		//reporte requerido
		carpeta := filepath.Dir(path)
		rutaReporte := carpeta + "/" + nombre + ".dot"
		respuesta += "Reporte BM_Bloque " + nombre +" creado \n"
		respuesta += " Pertenece al disco: " + nombreDisco

		Herramientas.RepGraphizMBR(rutaReporte, cad, nombre)
	}

	return respuesta
}

// =============================== BM INODE ===============================
func BM_inode(path string, id string) string{
	var respuesta string
	var pathDico string
	reportar := false

	//Busca en struck de particiones montadas el id ingresado
	for _,montado := range Structs.Montadas{
		if montado.Id == id{
			pathDico = montado.PathM
			reportar = true
		}
	}

	//Verifica que se encontro el ID y la Path del disco
	if pathDico == ""{
		reportar = false
		return "ERROR REP: ID NO ENCONTRADO"
	}

	if reportar{
		//Obtenermos el nombre del reporte que vamos a crear
		tmp := strings.Split(path, "/")
		nombre := strings.Split(tmp[len(tmp)-1], ".")[0]

		tmp2 := strings.Split(pathDico, "/")
		nombreDisco := strings.Split(tmp2[len(tmp2)-1], ".")[0]

		file, err := Herramientas.OpenFile(pathDico)
		if err != nil {
			return "ERROR REP SB OPEN FILE "+err.Error()
		}

		//Obtener mbr
		var mbr Structs.MBR
		// Read object from bin file
		if err := Herramientas.ReadObject(file, &mbr, 0); err != nil {
			return "ERROR REP SB READ FILE "+err.Error()
		}

		// Close bin file
		defer file.Close()

		//Encontrar la particion correcta
		part := -1
		for i := 0; i < 4; i++ {
			identificador := Structs.GetId(string(mbr.Partitions[i].Id[:]))
			if identificador == id {
				reportar = true
				part = i
				break //para que ya no siga recorriendo si ya encontro la particion independientemente si se pudo o no reducir
			}
		}

		var superBloque Structs.Superblock
		errREAD := Herramientas.ReadObject(file, &superBloque, int64(mbr.Partitions[part].Start))
		if errREAD != nil {
			fmt.Println("REP Error. Particion sin formato")
			return "REP Error. Particion sin formato"
		}

		cad := ""
		inicio := superBloque.S_bm_inode_start
		fin := superBloque.S_bm_block_start
		count := 1 //para contar el numero de caracteres por linea (maximo 20)

		//objeto para leer un byte decodificado
		var bm Structs.Bite

		for i := inicio; i < fin; i++ {
			//cargo el byte (struct de [1]byte) decodificado como las demas estructuras
			Herramientas.ReadObject(file, &bm, int64(i))

			if bm.Val[0] == 0 {
				cad += string("0 ")
			} else {
				cad += "1 "
			}

			if count == 20 {
				cad += "\n"
				count = 0
			}

			count++
		}

		//reporte requerido
		carpeta := filepath.Dir(path)//DIr es para obtener el directorio
		rutaReporte := carpeta + "/" + nombre + ".txt"
		Herramientas.Reporte(rutaReporte, cad)
		respuesta += "Reporte BM Inode " + nombre +" creado \n"
		respuesta += " Pertenece al disco: " + nombreDisco
	}

	return respuesta
}

// =============================== BM BLOQUE ===============================
func BM_Bloque(path string, id string) string{
	var respuesta string
	var pathDico string
	reportar := false

	//BUsca en struck de particiones montadas el id ingresado
	for _,montado := range Structs.Montadas{
		if montado.Id == id{
			pathDico = montado.PathM
			reportar = true
		}
	}

	//Verifica que se encontro el ID y la Path del disco
	if pathDico == ""{
		reportar = false
		return "ERROR REP: ID NO ENCONTRADO"
	}

	if reportar{
		//Obtenermos el nombre del reporte que vamos a crear
		tmp := strings.Split(path, "/")
		nombre := strings.Split(tmp[len(tmp)-1], ".")[0]

		tmp2 := strings.Split(pathDico, "/")
		nombreDisco := strings.Split(tmp2[len(tmp2)-1], ".")[0]

		file, err := Herramientas.OpenFile(pathDico)
		if err != nil {
			return "ERROR REP SB OPEN FILE "+err.Error()
		}

		var mbr Structs.MBR
		// Read object from bin file
		if err := Herramientas.ReadObject(file, &mbr, 0); err != nil {
			return "ERROR REP SB READ FILE "+err.Error()
		}

		// Close bin file
		defer file.Close()

		//Encontrar la particion correcta
		part := -1
		for i := 0; i < 4; i++ {		
			identificador := Structs.GetId(string(mbr.Partitions[i].Id[:]))
			if identificador == id {
				part = i
				break //para que ya no siga recorriendo si ya encontro la particion independientemente si se pudo o no reducir
			}
		}

		var superBloque Structs.Superblock
		errREAD := Herramientas.ReadObject(file, &superBloque, int64(mbr.Partitions[part].Start))
		if errREAD != nil {
			fmt.Println("REP Error. Particion sin formato")
			return "REP Error. Particion sin formato"
		}

		cad := ""
		inicio := superBloque.S_bm_block_start
		fin := superBloque.S_inode_start
		count := 1 //para contar el numero de caracteres por linea (maximo 20)

		//objeto para leer un byte decodificado
		var bm Structs.Bite

		for i := inicio; i < fin; i++ {
			//cargo el byte (struct de [1]byte) decodificado como las demas estructuras
			Herramientas.ReadObject(file, &bm, int64(i))

			if bm.Val[0] == 0 {
				cad += string("0 ")
			} else {
				cad += "1 "
			}

			if count == 20 {
				cad += "\n"
				count = 0
			}

			count++
		}


		//reporte requerido
		carpeta := filepath.Dir(path)//DIr es para obtener el directorio
		rutaReporte := carpeta + "/" + nombre + ".txt"
		Herramientas.Reporte(rutaReporte, cad)		
		respuesta += "Reporte BM_Bloque " + nombre +" creado \n"
		respuesta += " Pertenece al disco: " + nombreDisco
	}
	return respuesta
}

// =============================== FILE ===============================
func FILE(path string, id string, rutaFile string)string{
	var respuesta string
	var pathDico string
	var contenido string
	reportar := false

	//BUsca en struck de particiones montadas el id ingresado
	for _,montado := range Structs.Montadas{
		if montado.Id == id{
			pathDico = montado.PathM
			reportar = true
		}
	}

	//Verifica que se encontro el ID y la Path del disco
	if pathDico == ""{
		reportar = false
		return "ERROR REP: ID NO ENCONTRADO"
	}

	if reportar{
		//Obtenermos el nombre del reporte que vamos a crear
		tmp := strings.Split(path, "/")
		nombre := strings.Split(tmp[len(tmp)-1], ".")[0]

		tmp2 := strings.Split(pathDico, "/")
		nombreDisco := strings.Split(tmp2[len(tmp2)-1], ".")[0]

		Disco, err := Herramientas.OpenFile(pathDico)
		if err != nil {
			return "ERROR REP OPEN FILE "+err.Error()
		}

		var mbr Structs.MBR
		// Read object from bin file
		if err := Herramientas.ReadObject(Disco, &mbr, 0); err != nil {
			return "ERROR REP READ FILE "+err.Error()
		}

		// Close bin file
		defer Disco.Close()

		//Encontrar la particion correcta
		part := -1
		for i := 0; i < 4; i++ {		
			identificador := Structs.GetId(string(mbr.Partitions[i].Id[:]))
			if identificador == id {
				part = i
				break //para que ya no siga recorriendo si ya encontro la particion independientemente si se pudo o no reducir
			}
		}
		
		var superBloque Structs.Superblock
		var fileBlock Structs.Fileblock
		errREAD := Herramientas.ReadObject(Disco, &superBloque, int64(mbr.Partitions[part].Start))
		if errREAD != nil {
			fmt.Println("REP Error. Particion sin formato")
			return "REP Error. Particion sin formato"
		}

		//buscar el inodo que contiene el archivo buscado
		idInodo := toolsinodos.BuscarInodo(0, rutaFile, superBloque, Disco)
		var inodo Structs.Inode

		//idInodo: solo puede existir archivos desde el inodo 1 en adelante (-1 no existe, 0 es raiz)
		if idInodo > 0 {
			contenido += "Contenido del archivo: '"+rutaFile+"'\n"
			Herramientas.ReadObject(Disco, &inodo, int64(superBloque.S_inode_start+(idInodo*int32(binary.Size(Structs.Inode{})))))
			//recorrer los fileblocks del inodo para obtener toda su informacion
			for _, idBlock := range inodo.I_block {
				if idBlock != -1 {
					Herramientas.ReadObject(Disco, &fileBlock, int64(superBloque.S_block_start+(idBlock*int32(binary.Size(Structs.Fileblock{})))))
					tmpConvertir := Herramientas.EliminartIlegibles(string(fileBlock.B_content[:]))
					contenido += tmpConvertir				
				}
			}

			contenido += "\n"
			
		} else {
			fmt.Println("REP ERROR: No se encontro el archivo ", rutaFile)
			return "REP ERROR: No se encontro el archivo " + rutaFile
		}

		//reporte requerido
		carpeta := filepath.Dir(path)//DIr es para obtener el directorio
		rutaReporte := carpeta + "/" + nombre + ".txt"
		Herramientas.Reporte(rutaReporte, contenido)
		respuesta += "Reporte BM_Bloque " + nombre +" creado \n"
		respuesta += "Pertenece al disco: " + nombreDisco
	}
	return respuesta
}

// =============================== LS ===============================
func LS(path string, id string, rutaFile string)string{
	var respuesta string
	var contenido string
	var pathDico string
	reportar := false

	//BUsca en struck de particiones montadas el id ingresado
	for _,montado := range Structs.Montadas{
		if montado.Id == id{
			pathDico = montado.PathM
			reportar = true
		}
	}

	//Verifica que se encontro el ID y la Path del disco
	if pathDico == ""{
		reportar = false
		return "ERROR REP: ID NO ENCONTRADO"
	}

	if reportar{
		Color := "BlueViolet"	
		contenido = "digraph {\nnode [ shape=none ] \nTablaReportNodo [ label = < <table border=\"1\"> \n<tr>\n\t<td bgcolor='"+Color+"'>PERMISOS</td>\n\t<td bgcolor='"+Color+"'> USUARIO </td>\n\t<td bgcolor='"+Color+"'> GRUPO </td>\n\t<td bgcolor='"+Color+"'> SIZE </td>\n\t<td bgcolor='"+Color+"'> FECHA/HORA </td> \n\t<td bgcolor='"+Color+"'> NOMBRE </td>\n\t<td bgcolor='"+Color+"'> TIPO </td>\n </tr>"
		
		//Obtenermos el nombre del reporte que vamos a crear
		tmp := strings.Split(path, "/")
		nombre := strings.Split(tmp[len(tmp)-1], ".")[0]

		tmp2 := strings.Split(pathDico, "/")
		nombreDisco := strings.Split(tmp2[len(tmp2)-1], ".")[0]

		Disco, err := Herramientas.OpenFile(pathDico)
		if err != nil {
			return "ERROR REP OPEN FILE "+err.Error()
		}

		var mbr Structs.MBR
		// Read object from bin file
		if err := Herramientas.ReadObject(Disco, &mbr, 0); err != nil {
			return "ERROR REP READ FILE "+err.Error()
		}

		// Close bin file
		defer Disco.Close()

		//Encontrar la particion correcta
		part := -1
		for i := 0; i < 4; i++ {		
			identificador := Structs.GetId(string(mbr.Partitions[i].Id[:]))
			if identificador == id {
				part = i
				break //para que ya no siga recorriendo si ya encontro la particion independientemente si se pudo o no reducir
			}
		}
		
		//var fileBlock Structs.Fileblock
		var superBloque Structs.Superblock
		errREAD := Herramientas.ReadObject(Disco, &superBloque, int64(mbr.Partitions[part].Start))
		if errREAD != nil {
			fmt.Println("CAT ERROR. Particion sin formato")
			return "CAT ERROR. Particion sin formato" + "\n"
		}


		var FstInodo Structs.Inode		
		//Le agrego una structura de inodo para ver el user.txt que esta en el primer inodo del sb
		Herramientas.ReadObject(Disco, &FstInodo, int64(superBloque.S_inode_start + int32(binary.Size(Structs.Inode{}))))
			

		var contUs string
		var FistfileBlock Structs.Fileblock
		for _, item := range FstInodo.I_block {
			if item != -1 {
				Herramientas.ReadObject(Disco, &FistfileBlock, int64(superBloque.S_block_start+(item*int32(binary.Size(Structs.Fileblock{})))))
				contUs += string(FistfileBlock.B_content[:])
			}
		}
		lineaID := strings.Split(contUs, "\n")
		

		idInodo := toolsinodos.BuscarInodo(0, rutaFile, superBloque, Disco)
		var inodo Structs.Inode

		if idInodo > 0 {
			Herramientas.ReadObject(Disco, &inodo, int64(superBloque.S_inode_start+(idInodo*int32(binary.Size(Structs.Inode{})))))
			var folderBlock Structs.Folderblock
			for _, idBlock := range inodo.I_block {
				if idBlock != -1 {
					Herramientas.ReadObject(Disco, &folderBlock, int64(superBloque.S_block_start+(idBlock*int32(binary.Size(Structs.Folderblock{})))))					
					for k := 2; k < 4; k++ {
						apuntador := folderBlock.B_content[k].B_inodo
						if apuntador != -1 {
							pathActual := Structs.GetB_name(string(folderBlock.B_content[k].B_name[:]))
							
							contenido += InodoLs(pathActual, lineaID, apuntador , superBloque, Disco)
						}
					}					
				}
			}
			
			
		}else{
			respuesta = "REP ERROR NO SE ENCONTRO LA PATH INGRESADA"
		}

		contenido += "\n</table> > ]\n}"
		cad := Herramientas.EliminartIlegibles(contenido)

		//reporte requerido
		carpeta := filepath.Dir(path)//DIr es para obtener el directorio
		rutaReporte := carpeta + "/" + nombre + ".dot"
		Herramientas.Reporte(rutaReporte, contenido)
		respuesta += "Reporte BM_Bloque " + nombre +" creado \n"
		respuesta += "Pertenece al disco: " + nombreDisco
		Herramientas.RepGraphizMBR(rutaReporte, cad, nombre)
	}	
	return respuesta	
}

			// Nombre,   contenia users.txt		no. bloque		superbloque					DIsco
func InodoLs(name string,lineaID []string,  idInodo int32, superBloque Structs.Superblock, file *os.File)string{
	var contenido string

	//cargar el inodo a reportar
	var inodo Structs.Inode
	Herramientas.ReadObject(file, &inodo, int64(superBloque.S_inode_start+(idInodo*int32(binary.Size(Structs.Inode{})))))
	
	//Busco el grupo y el usuario 							
	usuario:= ""
	grupo:=""							
	for m:=0; m<len(lineaID); m++{
		datos := strings.Split(lineaID[m], ",")
		if len(datos) == 5 {	
			us := fmt.Sprintf("%d",inodo.I_uid)													
			if us== datos[0]{
				usuario = datos[3]
			}		
		}
		if len(datos) == 3 {	
			gr := fmt.Sprintf("%d",inodo.I_gid)									
			if gr== (datos[0]){
				grupo = datos[2]
			}		
		}

	}
	
	Color := "Pink"
	tipoArchivo := "Archivo"
	var permisos string	
	
	//Los permisos son 3 numeros porque son aplicados a: propierarios   grupos  y  otros
	//Cada numero representa los permisos de lectura, escritura y ejecucion: r w x
	// r lectura
	// w escritura
	// x ejecucion 
	//Si el numero de permisos es: 764, significa que:
	//el propierario(7) tiene permisos de lectura escritura ejecucion
	//el grupo(6) tiene permisos de lectura escritura
	//otros(4) tienen permisos de lectura
	for i:=0; i<3; i++{	
		if string(inodo.I_perm[i])=="0"{//ninun permiso
			permisos+="---"
		}else if string(inodo.I_perm[i])=="1"{// ejecucion
			permisos+="--x"
		}else if string(inodo.I_perm[i])=="2"{//	escritura
			permisos+="-w-"
		}else if string(inodo.I_perm[i])=="3"{// 	ecritura ejecucion
			permisos+="-wx"
		}else if string(inodo.I_perm[i])=="4"{//lectura
			permisos+="r--"
		}else if string(inodo.I_perm[i])=="5"{//lectura  	ejecucion
			permisos+="r-x"
		}else if string(inodo.I_perm[i])=="6"{// lectura escritura
			permisos+="rw-"
		}else if string(inodo.I_perm[i])=="7"{//lectura escritura ejecucion
			permisos+="rwx"
		}
	}

	if string(inodo.I_type[:]) == "0"{
		Color = "Violet"
		tipoArchivo = "Carpeta"	
	}
	
	permisos = "rw-rw-r--"	
	contenido += "\n  <tr>"
	contenido += "\n\t <td bgcolor='"+Color+"'> "+ permisos +"</td>"
	contenido += fmt.Sprintf("\n\t <td bgcolor='%s'> %s</td>",Color,usuario)
	contenido += fmt.Sprintf("\n\t <td bgcolor='%s'> %s</td>",Color,grupo)
	contenido += fmt.Sprintf("\n\t <td bgcolor='%s'> %d</td>", Color, inodo.I_size)
	contenido += fmt.Sprintf("\n\t <td bgcolor='%s'> %s </td> ", Color, string(inodo.I_ctime[:]))
	contenido += "\n\t <td bgcolor='"+Color+"'> "+ name +"</td>"
	contenido += "\n\t <td bgcolor='"+Color+"'> "+ tipoArchivo +"</td>"
	contenido += "\n  </tr>"
	//reportar el inodo
	return contenido
}

// =============================== Journal ===============================
func journal(path string, id string) string{
	var respuesta string
	var pathDico string
	reportar := false

	//Busca en struck de particiones montadas el id ingresado
	for _,montado := range Structs.Montadas{
		if montado.Id == id{
			pathDico = montado.PathM
			reportar = true
		}
	}

	//Verifica que se encontro el ID y la Path del disco
	if pathDico == ""{
		reportar = false
		return "ERROR REP: ID NO ENCONTRADO"
	}

	//if true { //para probar los reporte hayan o no particiones montadas
	if reportar {
		//Obtenermos el nombre del reporte que vamos a crear
		tmp := strings.Split(path, "/")
		nombre := strings.Split(tmp[len(tmp)-1], ".")[0]

		tmp2 := strings.Split(pathDico, "/")
		nombreDisco := strings.Split(tmp2[len(tmp2)-1], ".")[0]

		file, err := Herramientas.OpenFile(pathDico)
		if err != nil {
			return "ERROR REP SB OPEN FILE "+err.Error()
		}

		//Obtener mbr
		var mbr Structs.MBR
		// Read object from bin file
		if err := Herramientas.ReadObject(file, &mbr, 0); err != nil {
			return "ERROR REP SB READ FILE "+err.Error()
		}

		// Close bin file
		defer file.Close()

		//Encontrar la particion correcta
		part := -1
		for i := 0; i < 4; i++ {
			identificador := Structs.GetId(string(mbr.Partitions[i].Id[:]))
			if identificador == id {
				reportar = true
				part = i
				break //para que ya no siga recorriendo si ya encontro la particion independientemente si se pudo o no reducir
			}
		}

		cad := "digraph { \nnode [ shape=none ] \nTablaReportNodo [ label = < <table border=\"1\"> \n"
		cad += Structs.RepJournal(mbr.Partitions[part], file)
		cad += "</table> > ]\n}"

		//reporte requerido
		carpeta := filepath.Dir(path)//DIr es para obtener el directorio
		rutaReporte := carpeta + "/" + nombre + ".dot"
		Herramientas.Reporte(rutaReporte, cad)
		respuesta += "Reporte BM Inode " + nombre +" creado \n"
		respuesta += " Pertenece al disco: " + nombreDisco

		Herramientas.RepGraphizMBR(rutaReporte, cad, nombre)
	} else {
		fmt.Println("REP Error: Id no existe")
	}
	return respuesta
}

// =============================== TREE ===============================
func tree(path string, id string) string{
	var respuesta string
	var pathDico string
	reportar := false

	//Busca en struck de particiones montadas el id ingresado
	for _,montado := range Structs.Montadas{
		if montado.Id == id{
			pathDico = montado.PathM
			reportar = true
		}
	}

	//Verifica que se encontro el ID y la Path del disco
	if pathDico == ""{
		reportar = false
		return "ERROR REP: ID NO ENCONTRADO"
	}

	if reportar {
		//Obtenermos el nombre del reporte que vamos a crear
		tmp := strings.Split(path, "/")
		nombre := strings.Split(tmp[len(tmp)-1], ".")[0]

		tmp2 := strings.Split(pathDico, "/")
		nombreDisco := strings.Split(tmp2[len(tmp2)-1], ".")[0]

		file, err := Herramientas.OpenFile(pathDico)
		if err != nil {
			return "ERROR REP SB OPEN FILE "+err.Error()
		}

		//Obtener mbr
		var mbr Structs.MBR
		// Read object from bin file
		if err := Herramientas.ReadObject(file, &mbr, 0); err != nil {
			return "ERROR REP SB READ FILE "+err.Error()
		}

		// Close bin file
		defer file.Close()

		//Encontrar la particion correcta
		part := -1
		for i := 0; i < 4; i++ {
			identificador := Structs.GetId(string(mbr.Partitions[i].Id[:]))
			if identificador == id {
				reportar = true
				part = i
				break //para que ya no siga recorriendo si ya encontro la particion independientemente si se pudo o no reducir
			}
		}

		var superBloque Structs.Superblock
		err = Herramientas.ReadObject(file, &superBloque, int64(mbr.Partitions[part].Start))
		if err != nil {
			fmt.Println("REP Error. Particion sin formato")
			return "REP Error. Particion sin formato"
		}

		var Inode0 Structs.Inode
		Herramientas.ReadObject(file, &Inode0, int64(superBloque.S_inode_start))

		cad := "digraph { \n graph [pad=0.5, nodesep=0.5, ranksep=1] \n node [ shape=plaintext ] \n rankdir=LR \n"

		//reportar el inodo
		cad += "\n Inodo0 [ \n  label = < \n   <table border=\"0\" cellborder=\"1\" cellspacing=\"0\"> \n"
		cad += "    <tr> <td bgcolor='skyblue' colspan=\"2\" port='P0'> Inodo 0 </td> </tr> \n"

		for i := 0; i < 12; i++ {
			cad += fmt.Sprintf("    <tr> <td> AD%d </td> <td port='P%d'> %d </td> </tr> \n", i+1, i+1, Inode0.I_block[i])
		}
		//Separo los ultimos 3 para marcarlos con color diferente por ser indirectos
		for i := 12; i < 15; i++ {
			cad += fmt.Sprintf("    <tr> <td bgcolor='pink'> AD%d </td> <td port='P%d'> %d </td> </tr> \n", i+1, i+1, Inode0.I_block[i])
		}
		cad += "   </table> \n  > \n ]; \n"
		//fin primer inodo

		//llamar bloques
		for i := 0; i < 15; i++ {
			bloque := Inode0.I_block[i]
			if bloque != -1 {
				// No. bloque, tipo Inodo (carpeta/archivo), inodo padre, No port, superbloque, disco
				cad += treeBlock(bloque, string(Inode0.I_type[:]), 0, i+1, superBloque, file)
			}
		}
		//Inode0.I_block[12] -> trae un bloque indirecto antes de un bloque normal
		//Inode0.I_block[13] -> trae dos bloque indirecto antes de un bloque normal
		//Inode0.I_block[14] -> trae tres bloque indirecto antes de un bloque normal
		cad += "\n}"

		//reporte requerido
		carpeta := filepath.Dir(path)//DIr es para obtener el directorio
		rutaReporte := carpeta + "/" + nombre + ".dot"
		Herramientas.Reporte(rutaReporte, cad)
		respuesta += "Reporte tree " + nombre +" creado \n"
		respuesta += " Pertenece al disco: " + nombreDisco

		Herramientas.RepGraphizMBR(rutaReporte, cad, nombre)
	} else {
		fmt.Println("REP Error: Id no existe")
	}
	return respuesta
}

// Metodo recursivo del tree para buscar bloques
// .              No bloque,   tipo bloque,  inodo padre, No port,          superbloque,            disco
func treeBlock(idBloque int32, tipo string, idPadre int32, p int, superBloque Structs.Superblock, file *os.File) string {
	cad := fmt.Sprintf("\n Bloque%d [ \n  label = < \n   <table border=\"0\" cellborder=\"1\" cellspacing=\"0\"> \n", idBloque)

	if tipo == "0" {
		//FolderBlock
		var folderBlock Structs.Folderblock
		Herramientas.ReadObject(file, &folderBlock, int64(superBloque.S_block_start+(idBloque*int32(binary.Size(Structs.Folderblock{})))))

		//Reporte del bloque actual
		cad += fmt.Sprintf("    <tr> <td bgcolor='orchid' colspan=\"2\" port='P0'> Bloque %d </td> </tr> \n", idBloque)
		cad += fmt.Sprintf("    <tr> <td> . </td> <td port='P1'> %d </td> </tr> \n", folderBlock.B_content[0].B_inodo)
		cad += fmt.Sprintf("    <tr> <td> .. </td> <td port='P2'> %d </td> </tr> \n", folderBlock.B_content[1].B_inodo)
		cad += fmt.Sprintf("    <tr> <td> %s </td> <td port='P3'> %d </td> </tr> \n", Structs.GetB_name(string(folderBlock.B_content[2].B_name[:])), folderBlock.B_content[2].B_inodo)
		cad += fmt.Sprintf("    <tr> <td> %s </td> <td port='P4'> %d </td> </tr> \n", Structs.GetB_name(string(folderBlock.B_content[3].B_name[:])), folderBlock.B_content[3].B_inodo)
		cad += "   </table> \n  > \n ]; \n"
		//Enlazar inodo padre con bloque actual
		cad += fmt.Sprintf("\n Inodo%d:P%d -> Bloque%d:P0; \n", idPadre, p, idBloque) //p es el port del inodo que apunta al bloque actual
		//recorrero el folderblock para ver si apunta a otros inodos
		for i := 2; i < 4; i++ {
			inodo := folderBlock.B_content[i].B_inodo
			if inodo != -1 {
				//.       inodo hijo, bloque actual, No. port, superbloque, disco
				cad += treeInodo(inodo, idBloque, i+1, superBloque, file)
			}
		}
	} else {
		//Fileblock
		var fileBlock Structs.Fileblock
		Herramientas.ReadObject(file, &fileBlock, int64(superBloque.S_block_start+(idBloque*int32(binary.Size(Structs.Fileblock{})))))
		//Reporte del bloque actual
		cad += fmt.Sprintf("    <tr> <td bgcolor='#ffff99' port='P0'> Bloque %d </td> </tr> \n", idBloque)
		cad += fmt.Sprintf("    <tr> <td> %s </td> </tr> \n", Structs.GetB_content(string(fileBlock.B_content[:])))
		cad += "   </table> \n  > \n ]; \n"
		//Enlazar inodo padre con bloque actual
		cad += fmt.Sprintf("\n Inodo%d:P%d -> Bloque%d:P0; \n", idPadre, p, idBloque) //p es el port del inodo que apunta al bloque actual
	}

	return cad
}

// Metodo recursivo del tree para buscar inodos
// .            No. inode,     No. bloque,  No. port,          superbloque,            disco
func treeInodo(idInodo int32, idPadre int32, p int, superBloque Structs.Superblock, file *os.File) string {
	//cargar el inodo a reportar
	var Inode Structs.Inode
	Herramientas.ReadObject(file, &Inode, int64(superBloque.S_inode_start+(idInodo*int32(binary.Size(Structs.Inode{})))))

	//reportar el inodo
	cad := fmt.Sprintf("\n Inodo%d [ \n  label = < \n   <table border=\"0\" cellborder=\"1\" cellspacing=\"0\"> \n", idInodo)
	//color segun tipo de inodo
	if string(Inode.I_type[:]) == "0" {
		cad += fmt.Sprintf("    <tr> <td bgcolor='skyblue' colspan=\"2\" port='P0'> Inodo %d </td> </tr> \n", idInodo)
	} else {
		cad += fmt.Sprintf("    <tr> <td bgcolor='#7FC97F' colspan=\"2\" port='P0'> Inodo %d </td> </tr> \n", idInodo)
	}

	//recorrer los apuntadores
	for i := 0; i < 12; i++ {
		cad += fmt.Sprintf("    <tr> <td> AD%d </td> <td port='P%d'> %d </td> </tr> \n", i+1, i+1, Inode.I_block[i])
	}
	//Separo los ultimos 3 para marcarlos con color diferente por ser indirectos
	for i := 12; i < 15; i++ {
		cad += fmt.Sprintf("    <tr> <td bgcolor='pink'> AD%d </td> <td port='P%d'> %d </td> </tr> \n", i+1, i+1, Inode.I_block[i])
	}
	cad += "   </table> \n  > \n ]; \n"
	//fin inodo

	//Enlazar inodo padre con bloque actual
	cad += fmt.Sprintf("\n Bloque%d:P%d -> Inodo%d:P0; \n", idPadre, p, idInodo) //p es el port del inodo que apunta al bloque actual

	//llamar bloques
	for i := 0; i < 15; i++ {
		bloque := Inode.I_block[i]
		if bloque != -1 {
			//.          No. bloque, tipo Inodo (carpeta/archivo), inodo padre, port, superbloque, disco
			cad += treeBlock(bloque, string(Inode.I_type[:]), idInodo, i+1, superBloque, file)
		}
	}

	return cad
}