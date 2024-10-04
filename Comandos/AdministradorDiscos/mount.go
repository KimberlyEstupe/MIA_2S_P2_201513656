package administradordiscos

import (
	"MIA_2S_P2_201513656/Herramientas"
	"MIA_2S_P2_201513656/Structs"
	"fmt"
	"os"
	"strconv"
	"strings"
)

//var Pmontaje []Structs.Mount//GUarda en Ram las particones montadas
func Mount(entrada []string) (string){
	var respuesta string
	var name string	//Nobre de la particion a montar
	var pathE string	//Path del Disco
	Valido := true

	for _, parametro := range entrada[1:] {
		tmp := strings.TrimRight(parametro," ")
		valores := strings.Split(tmp,"=")

		if len(valores)!=2{
			fmt.Println("ERROR MOUNT, valor desconocido de parametros ",valores[1])
			respuesta += "ERROR MOUNT, valor desconocido de parametros " + valores[1]
			//Si falta el valor del parametro actual lo reconoce como error e interrumpe el proceso
			return respuesta
		}

		//******************* PATH *************
		if strings.ToLower(valores[0]) == "path" {
			pathE = strings.ReplaceAll(valores[1],"\"","")			
			_, err := os.Stat(pathE)
			if os.IsNotExist(err) {
				fmt.Println("ERROR MOUNT: El disco no existe")
				respuesta += "ERROR MOUNT: El disco no existe"
				Valido = false
				return respuesta // Terminar el bucle porque encontramos un nombre único
			}
		//********************  NAME *****************
		} else if strings.ToLower(valores[0]) == "name" {
			// Eliminar comillas
			name = strings.ReplaceAll(valores[1], "\"", "")
			// Eliminar espacios en blanco al final
			name = strings.TrimSpace(name)
		
		//******************* ERROR EN LOS PARAMETROS *************
		} else {
			fmt.Println("ERROR MOUNT: Parametro desconocido: ", valores[0])
			respuesta += "ERROR MOUNT: Parametro desconocido: "+ valores[0]
			return respuesta //por si en el camino reconoce algo invalido de una vez se sale
		}
	}

	if Valido{
		if pathE != ""{
			if name != ""{
				// Abrir y cargar el disco
				disco, err := Herramientas.OpenFile(pathE)
				if err != nil {
					respuesta += "ERROR NO SE PUEDE LEER EL DISCO " + err.Error()+ "\n"
					return  respuesta
				}	

				//Se crea un mbr para cargar el mbr del disco
				var mbr Structs.MBR
				//Guardo el mbr leido
				if err := Herramientas.ReadObject(disco, &mbr, 0); err != nil {
					respuesta += "ERROR Read " + err.Error()+ "\n"
					return  respuesta
				}
				
				// cerrar el archivo del disco
				defer disco.Close()

				montar := true //usar si se van a montar logicas
				reportar := false
				for i := 0; i < 4; i++ {
					nombre := Structs.GetName(string(mbr.Partitions[i].Name[:]))
					if nombre == name{
						montar = false
						if string(mbr.Partitions[i].Type[:]) != "E" {
							if string(mbr.Partitions[i].Status[:]) != "A" {
								var id string 							
								var nuevaLetra byte = 'A'// A
								contador := 1
								modificada := false															

								//Verifica si el path existe dentro de las particiones montadas
								for k:=0; k < len(Structs.Pmontaje); k++{
									if Structs.Pmontaje[k].MPath == pathE{
										//MOdifica el struct 
										Structs.Pmontaje[k].Cont = Structs.Pmontaje[k].Cont + 1
										contador = int(Structs.Pmontaje[k].Cont)										
										nuevaLetra = Structs.Pmontaje[k].Letter
										modificada = true	
										break 
									}
								}

								if !modificada{
									if len(Structs.Pmontaje) > 0{
										nuevaLetra = Structs.Pmontaje[len(Structs.Pmontaje)-1].Letter +1
									}
									Structs.AddPathM(pathE, nuevaLetra, 1)
								}

								id = "56"+strconv.Itoa(contador)+string(nuevaLetra) //Id de particion
								fmt.Println("ID:  Letra ", string(nuevaLetra), " cont ", contador)
								//Agregar al struct de montadas
								Structs.AddMontadas(id, pathE)

								//TODO modificar la particion que se va a montar								
								//copy(mbr.Partitions[i].Status[:], "A")
								copy(mbr.Partitions[i].Id[:], id)
								mbr.Partitions[i].Correlative = int32(contador)

								//sobreescribir el mbr para guardar los cambios
								if err := Herramientas.WriteObject(disco, mbr, 0); err != nil { //Sobre escribir el mbr
									respuesta += "Error "
									return "Error "+err.Error()
								}
								reportar = true

								respuesta+="Particion con nombre "+ name+ " montada correctamente. ID: "+id
								fmt.Println("Particion con nombre ", name, " montada correctamente. ID: ",id)
							}else{
								fmt.Println("ERROR MOUNT. ESTA PARTICION YA FUE MONTADA PREVIAMENTE")
								respuesta += "ERROR MOUNT. ESTA PARTICION YA FUE MONTADA PREVIAMENTE"
								return respuesta
							}
						}else{
							fmt.Println("ERROR MOUNT. No se puede montar una particion extendida")
							respuesta += "ERROR MOUNT. No se puede montar una particion extendida"
							return respuesta	
						}
					}
				}

				if montar {
					fmt.Println("ERROR MOUNT. No se pudo montar la particion ", name)
					fmt.Println("ERROR MOUNT. No se encontro la particion")
					respuesta += "ERROR MOUNT. NO SE ENCONTRO LA PARTICION " + name
					respuesta += "\nNO SE PUDO MONTAR LA PARICION \n"
					return respuesta
				}

				if reportar {
					partMontadas :="\n\nLISTA DE PARTICIONES MONTADAS EN EL DISCO\n"
					for i := 0; i < 4; i++ {
						estado := string(mbr.Partitions[i].Status[:])
						if estado == "A" {
							tmpMontadas:= "Particion: " + strconv.Itoa(i) + ", name: " +string(mbr.Partitions[i].Name[:]) + ", status: "+string(mbr.Partitions[i].Status[:])+", id: "+string(mbr.Partitions[i].Id[:])+", tipo: "+string(mbr.Partitions[i].Type[:])+", correlativo: "+ strconv.Itoa(int(mbr.Partitions[i].Correlative)) + ", fit: "+string(mbr.Partitions[i].Fit[:])+ ", start: "+strconv.Itoa(int(mbr.Partitions[i].Start))+ ", size: "+strconv.Itoa(int(mbr.Partitions[i].Size))
							partMontadas += Herramientas.EliminartIlegibles(tmpMontadas)+"\n"
						}
					}

					partMontadas +="\n\n\tPARTICIONES MONTADAS\n"
					for _,montada := range Structs.Montadas{
						partMontadas += "Id "+ string(montada.Id)+ ", Disco: "+ montada.PathM+"\n"
					}					
					respuesta += partMontadas
					fmt.Println(partMontadas)
				}
			}else{
				fmt.Println("ERROR: FALTA NAME  EN MOUNT")	
				respuesta += "ERROR: FALTA NAME  EN MOUNT"			
			}
		}else{
			fmt.Println("ERROR: FALTA PARAMETRO PATH EN MOUNT")
			respuesta += "ERROR: FALTA PATH EN MOUNT"	
		}
	}

	return respuesta
	
}