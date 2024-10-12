package administradordiscos

import (
	"MIA_2S_P2_201513656/Herramientas"
	"MIA_2S_P2_201513656/Structs"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Fdisk(entrada []string) string{
	var respuesta string
	//size, path, name obligatorios
	//unit, type, fit
	unit:=1024 	//Valores B,K,M; por defauld es en k
	tipe:="P"	//Valores P(primaria) E(extendida) L(Logica); por defauld P
	fit :="W"	//Puede ser FF, BF, WF, por default es FF
	var size int			//Obligatorio	
	var pathE string		//Obligatorio
	var name string			//Obligatorio
	
	var add int           //opcional (para aumentar o reducir el tamaño de una particion)
	var delete int		  //1-> full; 2->fast
	var opcion int        // 0 -> crear; 1 -> add; 2 -> delete(por defecto es 0 = CREAR)
	
	Valido := true        //Para validar que los parametros cumplen con los requisitos
	var sizeValErr string //Para reportar el error si no se pudo convertir a entero el size

	for _,parametro :=range entrada[1:]{
		tmp := strings.TrimRight(parametro," ")
		valores := strings.Split(tmp,"=")

		if len(valores)!=2{
			fmt.Println("ERROR FDISK, valor desconocido de parametros ",valores[1])
			Valido = false
			//Si falta el valor del parametro actual lo reconoce como error e interrumpe el proceso
			return "ERROR FDISK, valor desconocido de parametros "+valores[1]
		}

		//********************  SIZE *****************
		if strings.ToLower(valores[0])=="size"{
			var err error
			//size, err = strconv.Atoi(tmp) //se convierte el valor en un entero
			//if err != nil || size <= 0 { //Se manejaria como un solo error
			size, err = strconv.Atoi(valores[1]) //se convierte el valor en un entero
			if err != nil {
				sizeValErr = valores[1] //guarda para el reporte del error si es necesario validar size
			}

		//*************** UNIT ***********************
		} else if strings.ToLower(valores[0]) == "unit" {
			//si la unidad es k
			if strings.ToLower(valores[1]) == "b" {
				unit = 1
				//si la unidad no es k ni m es error (si fuera m toma el valor con el que se inicializo unit al inicio del metodo)
			} else if strings.ToLower(valores[1]) == "m" {
				unit = 1048576 //1024*1024
			} else if strings.ToLower(valores[1]) != "k" {
				Valido = false
				fmt.Println("ERROR FDISK en -unit. Valores aceptados: b, k, m. ingreso: ", valores[1])
				return "ERROR FDISK en -unit. Valores aceptados: b, k, m. ingreso: "+ valores[1]				
			}

		//******************* PATH *************
		} else if strings.ToLower(valores[0]) == "path" {
			pathE = strings.ReplaceAll(valores[1],"\"","")

			_, err := os.Stat(pathE)
			if os.IsNotExist(err) {
				fmt.Println("ERROR FDISK: El disco no existe")
				return "ERROR FDISK: El disco no existe"// Terminar el bucle porque encontramos un nombre único
			}
		
		//******************* Type *************		
		} else if strings.ToLower(valores[0]) == "type" {
			//p esta predeterminado
			if strings.ToLower(valores[1]) == "e" {
				tipe = "E"
			} else if strings.ToLower(valores[1]) == "l" {
				tipe = "L"
			} else if strings.ToLower(valores[1]) != "p" {
				fmt.Println("ERROR FDISK en -type. Valores aceptados: e, l, p. ingreso: ", valores[1])
				return "ERROR FDISK en -type. Valores aceptados: e, l, p. ingreso: "+ valores[1]
			}

		//********************  Fit *****************
		}else if strings.ToLower(valores[0])=="fit"{
			if strings.ToLower(strings.TrimSpace(valores[1]))=="bf"{
				fit = "B"
			}else if strings.ToLower(valores[1])=="ff"{
				fit = "F"
			}else if strings.ToLower(valores[1])!="wf"{
				fmt.Println("EEROR: PARAMETRO FIT INCORRECTO. VALORES ACEPTADO: FF, BF,WF. SE INGRESO:",valores[1])
				return "EEROR: PARAMETRO FIT INCORRECTO. VALORES ACEPTADO: FF, BF,WF. SE INGRESO:"+valores[1]
			}
			
			
		//********************  NAME *****************
		} else if strings.ToLower(valores[0]) == "name" {
			// Eliminar comillas
			name = strings.ReplaceAll(valores[1], "\"", "")
			// Eliminar espacios en blanco al final
			name = strings.TrimSpace(name)		
		
		//******************** DELETE *****************	
		} else if strings.ToLower(valores[0]) == "delete" {	
			if strings.ToLower(valores[1]) == "full" {
				if opcion == 0 {
					opcion = 2 // 2 es delete
					delete = 1
				}
			} else if strings.ToLower(valores[1]) == "fast" {
				if opcion == 0 {
					delete = 2
					opcion = 2 // 2 es delete
				}
			} else {
				fmt.Println("ERROR FDISK. Valor de delete desconocido")
				Valido = false
				return "ERROR FDISK. Valor de delete desconocido"
			}
		//******************** ADD *****************
		} else if strings.ToLower(valores[0]) == "add" {
			var err error
			add, err = strconv.Atoi(valores[1]) //se convierte el valor en un entero
			if err != nil {
				fmt.Println("ERROR FDISK: El valor de \"add\" debe ser un valor numerico. se leyo ", valores[1])
				Valido = false
				return "ERROR FDISK: El valor de \"add\" debe ser un valor numerico. se leyo "+ valores[1]
			} else {
				if opcion == 0 {
					opcion = 1
				}
			}
		//******************* ERROR EN LOS PARAMETROS *************
		} else {
			fmt.Println("ERROR FDISK: Parametro desconocido: ", valores[0])
			return "ERROR FDISK: Parametro desconocido: "+ valores[0] //por si en el camino reconoce algo invalido de una vez se sale
		}
	}

	// ------------ VALIDAR PARAMETROS OBLIGATORIOS -----------
	if size != 0{
		if sizeValErr == "" { //Si es un numero (si es numero la variable sizeValErr sera una cadena vacia)
			if size <= 0 { //se valida que sea mayor a 0 (positivo)
				fmt.Println("ERROR FDISK: -size debe ser un valor positivo mayor a cero (0). se leyo ", size)
				Valido = false
				return "ERROR FDISK: -size debe ser un valor positivo mayor a cero (0). se leyo " + string(size)
			}
		} else { //Si sizeValErr es una cadena (por lo que no se pudo dar valor a size)
			fmt.Println("ERROR FDISK: -size debe ser un valor numerico. se leyo ", sizeValErr)
			Valido = false
			return "ERROR FDISK: -size debe ser un valor numerico. se leyo "+ sizeValErr
		}
	}else{
		fmt.Println("ERROR FDISK: FALTO PARAMETRO SIZE")
		Valido =false
		return "ERROR FDISK: FALTO PARAMETRO SIZE"
	}

	if pathE == ""{
		fmt.Println("ERROR FDISK: FALTA PARAMETRO PATH")
		Valido = false
		return "ERROR FDISK: FALTA PARAMETRO PATH"
	}
	if name == ""{
		fmt.Println("ERROR FDISK: FALTA PARAMETRO NAME")
		Valido = false
		return "ERROR FDISK: FALTA PARAMETRO NAME"
	}

	if Valido{
		//Parametros correctos, se puede comenzar a crear las particiones
		disco, err := Herramientas.OpenFile(pathE)
		if err != nil {
			fmt.Println("ERROR FDISK: No se pudo leer el disco")
			return "ERROR FDISK: No se pudo leer el disco"+ "\n"
		}

		//Se crea un mbr para cargar el mbr del disco
		var mbr Structs.MBR
		//Guardo el mbr leido
		if err := Herramientas.ReadObject(disco, &mbr, 0); err != nil {
			return "ERROR FDISK Read " + err.Error()+ "\n"
		}

		// ==================== CREACION: OPCION 1====================
		if opcion == 0{
			//Si la particion es tipo extendida validar que no exista alguna extendida
			isPartExtend := false //Indica si se puede usar la particion extendida
			isName := true        //Valida si el nombre no se repite (true no se repite)
			if tipe == "E" {
				for i := 0; i < 4; i++ {
					tipo := string(mbr.Partitions[i].Type[:])
					
					if tipo != "E" {
						isPartExtend = true
					} else {
						isPartExtend = false
						isName = false //Para que ya no evalue el nombre ni intente hacer nada mas
						fmt.Println("ERROR FDISK. Ya existe una particion extendida")
						fmt.Println("ERROR FDISK. No se puede crear la nueva particion con nombre: ", name)
						return "ERROR FDISK. Ya existe una particion extendida \nFDISK Error. No se puede crear la nueva particion con nombre:  " + name+ "\n"
					}
				}
			}

			//verificar si  el nombre existe en las particiones primarias o extendida
			if isName {
				for i := 0; i < 4; i++ {
					nombre := Structs.GetName(string(mbr.Partitions[i].Name[:]))
					if nombre == name {
						isName = false
						fmt.Println("ERROR FDISK. Ya existe la particion : ", name)
						fmt.Println("ERROR FDISK. No se puede crear la nueva particion con nombre: ", name)
						return "ERROR FDISK. Ya existe la particion : " + name + "\nFDISK Error. No se puede crear la nueva particion con nombre: " + name+ "\n"

					}
				}
			}

			if isName{
				//Buscar en las logicas si ya existe
				var partExtendida Structs.Partition
				//buscar en que particion esta la particion extendida y guardarla en partExtend
				if string(mbr.Partitions[0].Type[:]) == "E" {
					partExtendida = mbr.Partitions[0]
				} else if string(mbr.Partitions[1].Type[:]) == "E" {
					partExtendida = mbr.Partitions[1]
				} else if string(mbr.Partitions[2].Type[:]) == "E" {
					partExtendida = mbr.Partitions[2]
				} else if string(mbr.Partitions[3].Type[:]) == "E" {
					partExtendida = mbr.Partitions[3]
				}

				if partExtendida.Size != 0{
					var actual Structs.EBR
					if err := Herramientas.ReadObject(disco, &actual, int64(partExtendida.Start)); err != nil {
						return "ERROR FDISK Read " + err.Error()+ "\n"
					}

					//Evaluo la primer ebr
					if Structs.GetName(string(actual.Name[:])) == name {
						isName = false
					} else{
						for actual.Next != -1 {
							//actual = actual.next
							if err := Herramientas.ReadObject(disco, &actual, int64(actual.Next)); err != nil {
								return "ERROR FDISK Read " + err.Error()+ "\n"
							}
							if Structs.GetName(string(actual.Name[:])) == name {
								isName = false
								break
							}
						}
					}

					if !isName {
						fmt.Println("ERROR FDISK. Ya existe la particion : ", name)
						fmt.Println("ERROR FDISK. No se puede crear la nueva particion con nombre: ", name)
						respuesta += "ERROR FDISK. Ya existe la particion : " + name
						respuesta += "\nFDISK Error. No se puede crear la nueva particion con nombre: " + name+ "\n"
						return respuesta
						
					}
				}
			}

			//INGRESO DE PARTICIONES PRIMARIAS Y/O EXTENDIDA (SIN LOGICAS)
			sizeNewPart := size * unit //Tamaño de la nueva particion (tamaño * unidades)
			guardar := false           //Indica si se debe guardar la particion, es decir, escribir en el disco
			var newPart Structs.Partition
			if (tipe == "P" || isPartExtend) && isName{//para que  isPartExtend sea true, typee tendra que ser "E"
				sizeMBR := int32(binary.Size(mbr)) //obtener el tamaño del mbr (el que ocupa fisicamente: 165)
				//Para manejar los demas ajustes hacer un if del fit para llamar a la funcion adecuada
				//F = primer ajuste; B = mejor ajuste; else -> peor ajuste

				//INSERTAR PARTICION (Primer ajuste)
				var resTem string
				mbr, newPart, resTem = primerAjuste(mbr, tipe, sizeMBR, int32(sizeNewPart), name, fit) //int32(sizeNewPart) es para castear el int a int32 que es el tipo que tiene el atributo en el struct Partition
				respuesta += resTem
				guardar = newPart.Size != 0

				//escribimos el MBR en el archivo. Lo que no se llegue a escribir en el archivo (aqui) se pierde, es decir, los cambios no se guardan
				if guardar{
					//sobreescribir el mbr
					if err := Herramientas.WriteObject(disco, mbr, 0); err != nil {
						return "ERROR FDISK Write " +err.Error()+ "\n"
					}

					//Se agrega el ebr de la particion extendida en el disco
					if isPartExtend {
						var ebr Structs.EBR
						ebr.Start = newPart.Start
						ebr.Next = -1
						if err := Herramientas.WriteObject(disco, ebr, int64(ebr.Start)); err != nil {
							return "ERROR FDISK Write " +err.Error()+ "\n"
						}
					}
					//para verificar que lo guardo
					var TempMBR2 Structs.MBR
					// Read object from bin file
					if err := Herramientas.ReadObject(disco, &TempMBR2, 0); err != nil {
						return "ERROR FDISK Read " + err.Error()+ "\n"
					}
					Structs.PrintMBR(TempMBR2)
					fmt.Println("\nParticion con nombre " + name + " creada exitosamente")
					respuesta += "\nParticion con nombre " + name + " creada exitosamente"+ "\n"					
				}else {
					//Lo podría eliminar pero tendria que modificar en el metodo del ajuste todos los errores para que aparezca el nombre que se intento ingresar como nueva particion
					fmt.Println("ERROR FDISK. No se puede crear la nueva particion con nombre: ", name)
					return "ERROR FDISK. No se puede crear la nueva particion con nombre: "+ name
				}
			}else if tipe == "L" && isName{
				var partExtend Structs.Partition
				if string(mbr.Partitions[0].Type[:]) == "E" {
					partExtend = mbr.Partitions[0]
				} else if string(mbr.Partitions[1].Type[:]) == "E" {
					partExtend = mbr.Partitions[1]
				} else if string(mbr.Partitions[2].Type[:]) == "E" {
					partExtend = mbr.Partitions[2]
				} else if string(mbr.Partitions[3].Type[:]) == "E" {
					partExtend = mbr.Partitions[3]
				} else {
					fmt.Println("ERROR FDISK. No existe una particion extendida en la cual crear un particion logica")
					return"ERROR FDISK. No existe una particion extendida en la cual crear un particion logica"+ "\n"
				}

				//valido que la particion extendida si exista (podría haber entrado al error que no existe extendida)
				if partExtend.Size != 0 {
					//si tuviera los demas ajustes con un if del fit y uso el metodo segun ajuste
					respuesta += primerAjusteLogicas(disco, partExtend, int32(sizeNewPart), name, fit) + "\n"//int32(sizeNewPart) es para castear el int a int32 que es el tipo que tiene el atributo en el struct Partition
					//repLogicas(partExtend, disco)
				}
				return respuesta
			}
		// =====================================================================================
		// ======================================== ADD ========================================
		// =====================================================================================
		}else if opcion == 1 {
			add = add * unit
			//-------------------------si se quita espacio----------------------------------------------------------------------
			//Particiones extendida o primarias
			if add < 0 {
				fmt.Println("Reducir espacio")
				reducir := true //Si cambia a false es que redujo una de las primarias o la extendida
				for i := 0; i < 4; i++{
					nombre := Structs.GetName(string(mbr.Partitions[i].Name[:]))
					if nombre == name{
						reducir = false
						newSize := mbr.Partitions[i].Size + int32(add)
						if newSize > 0{
							mbr.Partitions[i].Size += int32(add)
							if err := Herramientas.WriteObject(disco, mbr, 0); err != nil { //Sobre escribir el mbr
								return "ERROR FDISK: " + err.Error()
							}
							fmt.Println("Particion con nombre ", name, " se redujo correctamente")
							return "Particion con nombre "+ name+" se redujo correctamente"
						}else {
							fmt.Println("ERROR FDISK. El tamaño que intenta eliminar es demasiado grande")
							return "ERROR FDISK. El tamaño que intenta eliminar es demasiado grande"
						}
					}
				}

				//particiones logicas
				if reducir{
					var partExtendida Structs.Partition
					//buscar en que particion esta la particion extendida y guardarla en partExtend
					if string(mbr.Partitions[0].Type[:]) == "E" {
						partExtendida = mbr.Partitions[0]
					} else if string(mbr.Partitions[1].Type[:]) == "E" {
						partExtendida = mbr.Partitions[1]
					} else if string(mbr.Partitions[2].Type[:]) == "E" {
						partExtendida = mbr.Partitions[2]
					} else if string(mbr.Partitions[3].Type[:]) == "E" {
						partExtendida = mbr.Partitions[3]
					}

					//Si existe la extendida
					if partExtendida.Size != 0{
						var actual Structs.EBR
						if err := Herramientas.ReadObject(disco, &actual, int64(partExtendida.Start)); err != nil {
							return "ERROR FDISK, READ "+ err.Error()
						}

						//Evaluar si es la primera
						if Structs.GetName(string(actual.Name[:])) == name {
							reducir = false
						} else {
							for actual.Next != -1 {
								//actual = actual.next
								if err := Herramientas.ReadObject(disco, &actual, int64(actual.Next)); err != nil {
									return "ERROR FDISK, READ "+ err.Error()
								}
								if Structs.GetName(string(actual.Name[:])) == name {
									reducir = false
									break
								}
							}
						}

						if !reducir {
							actual.Size += int32(add)
							if actual.Size > 0 {
								if err := Herramientas.WriteObject(disco, actual, int64(actual.Start)); err != nil { //Sobre escribir el ebr
									return "ERROR FDISK, write "+ err.Error()
								}
								fmt.Println("Particion con nombre ", name, " se redujo correctamente")
								return "Particion con nombre "+ name+ " se redujo correctamente"
							} else {
								fmt.Println("ERROR FDISK. El tamaño que intenta eliminar es demasiado grande")
								return "ERROR FDISK. El tamaño que intenta eliminar es demasiado grande"
							}
						}
					}
				}

				if reducir {
					fmt.Println("ERROR FDISK. No existe la particion a reducir")
					return "ERROR FDISK. No existe la particion a reducir"
				}
			//Fin reducir espacio
			}else if add > 0{
				fmt.Println("aumentar espacio")
				//Primarias y/o extendida
				evaluar := 0
				//---------------------  Si el aumento es en particion 1 ---------------------
				if Structs.GetName(string(mbr.Partitions[0].Name[:])) == name{
					if mbr.Partitions[1].Start == 0 {
						if mbr.Partitions[2].Start == 0 {
							if mbr.Partitions[3].Start == 0 {
								evaluar = int(mbr.MbrSize - mbr.Partitions[0].GetEnd())
							} else {
								evaluar = int(mbr.Partitions[3].Start - mbr.Partitions[0].GetEnd())
							}
						} else {
							evaluar = int(mbr.Partitions[2].Start - mbr.Partitions[0].GetEnd())
						}
					} else {
						evaluar = int(mbr.Partitions[1].Start - mbr.Partitions[0].GetEnd())
					}

					//evaluar > 0 -> si hay espacio para aumentar. add <= evaluar -> si lo que quiero aumentar cabe en el espacio disponible
					if evaluar > 0 && add <= evaluar {
						//aumenta el tamaño de 1
						mbr.Partitions[0].Size += int32(add)
						if err := Herramientas.WriteObject(disco, mbr, 0); err != nil { //Sobre escribir el mbr
							return "ERROR FDISK, write "+ err.Error()
						}
						fmt.Println("Particion con nombre ", name, " aumento el espacio correctamente")
						return "Particion con nombre "+ name+ " aumento el espacio correctamente"
					} else {
						fmt.Println("ERROR FDISK. El tamaño que intenta aumentar es demasiado grande para la particion ", name)
						return "ERROR FDISK. El tamaño que intenta aumentar es demasiado grande para la particion " + name
					}
				//---------------------  Si el aumento es en particion 2 ---------------------
				}else if Structs.GetName(string(mbr.Partitions[1].Name[:])) == name{
					if mbr.Partitions[2].Start == 0 {
						if mbr.Partitions[3].Start == 0 {
							evaluar = int(mbr.MbrSize - mbr.Partitions[1].GetEnd())
						} else {
							evaluar = int(mbr.Partitions[3].Start - mbr.Partitions[1].GetEnd())
						}
					} else {
						evaluar = int(mbr.Partitions[2].Start - mbr.Partitions[1].GetEnd())
					}
					//aumenta el tamaño de 2
					if evaluar > 0 && add <= evaluar {
						mbr.Partitions[1].Size += int32(add)
						if err := Herramientas.WriteObject(disco, mbr, 0); err != nil { //Sobre escribir el mbr
							return "ERROR FDISK WRITE "+err.Error()
						}
						fmt.Println("Particion con nombre ", name, " aumento el espacio correctamente")
						return "Particion con nombre "+ name+ " aumento el espacio correctamente"
					} else {
						fmt.Println("ERROR FDISK. El tamaño que intenta aumentar es demasiado grande para la particion ", name)
						return "ERROR FDISK. El tamaño que intenta aumentar es demasiado grande para la particion " + name
					}
				//---------------------  Si el aumento es en particion 3 ---------------------
				}else if Structs.GetName(string(mbr.Partitions[2].Name[:])) == name{
					if mbr.Partitions[3].Start == 0 {
						evaluar = int(mbr.MbrSize - mbr.Partitions[2].GetEnd())
					} else {
						evaluar = int(mbr.Partitions[3].Start - mbr.Partitions[2].GetEnd())
					}
					//aumenta el tamaño de 3
					if evaluar > 0 && add <= evaluar {
						mbr.Partitions[2].Size += int32(add)
						if err := Herramientas.WriteObject(disco, mbr, 0); err != nil { //Sobre escribir el mbr
							return "ERROR FDISK WRITE "+err.Error()
						}
						fmt.Println("Particion con nombre ", name, " aumento el espacio correctamente")
						return "Particion con nombre " + name+ " aumento el espacio correctamente"
					} else {
						fmt.Println("ERROR FDISK. El tamaño que intenta aumentar es demasiado grande para la particion ", name)
						return "ERROR FDISK. El tamaño que intenta aumentar es demasiado grande para la particion "+name
					}
				//---------------------  Si el aumento es en particion 4 ---------------------
				}else if Structs.GetName(string(mbr.Partitions[3].Name[:])) == name{
					evaluar = int(mbr.MbrSize - mbr.Partitions[3].GetEnd())
					//aumenta el tamaño de 4
					if evaluar > 0 && add <= evaluar {
						mbr.Partitions[3].Size += int32(add)
						if err := Herramientas.WriteObject(disco, mbr, 0); err != nil { //Sobre escribir el mbr
							return "ERROR FDISK, WRITE "+err.Error()
						}
						fmt.Println("Particion con nombre ", name, " aumento el espacio correctamente")
						return "Particion con nombre "+ name+ " aumento el espacio correctamente"
					} else {
						fmt.Println("ERROR FDISK. El tamaño que intenta aumentar es demasiado grande para la particion ", name)
						return "ERROR FDISK. El tamaño que intenta aumentar es demasiado grande para la particion "+ name
					}		
				//---------------------  Si el aumento logicas ---------------------			
				}else{
					//Aumentar logica
					var partExtendida Structs.Partition
					//buscar en que particion esta la particion extendida y guardarla en partExtend
					if string(mbr.Partitions[0].Type[:]) == "E" {
						partExtendida = mbr.Partitions[0]
					} else if string(mbr.Partitions[1].Type[:]) == "E" {
						partExtendida = mbr.Partitions[1]
					} else if string(mbr.Partitions[2].Type[:]) == "E" {
						partExtendida = mbr.Partitions[2]
					} else if string(mbr.Partitions[3].Type[:]) == "E" {
						partExtendida = mbr.Partitions[3]
					}

					//Si existe la extendida
					if partExtendida.Size != 0{
						aumentar := false
						var actual Structs.EBR
						if err := Herramientas.ReadObject(disco, &actual, int64(partExtendida.Start)); err != nil {
							return "ERROR FDISK " +err.Error()
						}

						//Reviso si es la primera
						if Structs.GetName(string(actual.Name[:])) == name {
							aumentar = true
						} else{
							for actual.Next != -1 {
								//actual = actual.next
								if err := Herramientas.ReadObject(disco, &actual, int64(actual.Next)); err != nil {
									return "ERROR FDISK " +err.Error()
								}
								if Structs.GetName(string(actual.Name[:])) == name {
									aumentar = true
									break
								}
							}
						}

						if aumentar {
							if actual.Next != -1 {
								if add <= int(actual.Next)-int(actual.GetEnd()) {
									actual.Size += int32(add)
									if err := Herramientas.WriteObject(disco, actual, int64(actual.Start)); err != nil { //Sobre escribir el ebr
										return "ERROR FDISK "+err.Error()
									}
									fmt.Println("Particion con nombre ", name, " aumento el espacio correctamente")
								} else {
									fmt.Println("ERROR FDISK. El tamaño que intenta aumentar es demasiado grande para la particion ", name)
								}
							}else{
								if add <= int(partExtendida.GetEnd())-int(actual.GetEnd()) {
									actual.Size += int32(add)
									if err := Herramientas.WriteObject(disco, actual, int64(actual.Start)); err != nil { //Sobre escribir el ebr
										return "ERROR FDISK "+err.Error()
									}
									fmt.Println("Particion con nombre ", name, " aumento el espacio correctamente")
									return "Particion con nombre "+ name+ " aumento el espacio correctamente"
								} else {
									fmt.Println("ERROR FDISK. El tamaño que intenta aumentar es demasiado grande para la particion ", name)
									return "ERROR FDISK. El tamaño que intenta aumentar es demasiado grande para la particion " + name
								}
							}
						}else {
							fmt.Println("ERROR FDISK. No existe la particion a aumentar")
							return "ERROR FDISK. No existe la particion a aumentar"
						}
					}else{
						fmt.Println("ERROR FDISK. No existe particion extendida")
						return "ERROR FDISK. No existe particion extendida"
					}
				}
			} else {
				fmt.Println("ERROR FDISK. 0 no es un valor valido para aumentar o disminuir particiones")
				return "ERROR FDISK. 0 no es un valor valido para aumentar o disminuir particiones"
			}


		// ==========================================================================================
		// ======================================== ELIMINAR ========================================
		// ==========================================================================================
		}else if opcion == 2 {
			//-------- primarias o extendida-----------------------------------------------------
			del := true //para saber si se elimino la particion (true es que no se elimino, esto para facilitar el if que valida esta varible)
			for i := 0; i < 4; i++ {
				nombre := Structs.GetName(string(mbr.Partitions[i].Name[:]))
				if nombre == name {
					if delete == 1{// Elimina full
						Herramientas.DeletePart(disco, int64(mbr.Partitions[i].Start), mbr.Partitions[i].Size)
					}
					var newPart Structs.Partition
					mbr.Partitions[i] = newPart
					if err := Herramientas.WriteObject(disco, mbr, 0); err != nil { //Sobre escribir el mbr
						return "ERROR FDISK WRITE " + err.Error()
					}
					del = false
					fmt.Println("particion con nombre ", name, " eliminada")
					return "particion con nombre "+ name+ " eliminada"
				}
			}

			// -------------------------- Particiones LOgicas--------------------------------
			if del{
				fmt.Println("Eliminar particiones logicas")
			}
		}else {
			fmt.Println("ERROR FDISK. Operación desconocida (operaciones aceptadas: crear, modificar o eliminar)")
			return "ERROR FDISK. Operación desconocida (operaciones aceptadas: crear, modificar o eliminar)"
		}

		// Cierro el disco
		defer disco.Close()
	}//fin valido

	return respuesta
}

// ============================================= PRIMER AJUSTE ==================================
func primerAjuste(mbr Structs.MBR, typee string, sizeMBR int32, sizeNewPart int32, name string, fit string) (Structs.MBR, Structs.Partition, string) {
	var respuesta string
	var newPart Structs.Partition
	var noPart Structs.Partition //para revertir el set info (simula volverla null)

	//PARTICION 1 (libre) - (size = 0 no se ha creado)
	if mbr.Partitions[0].Size == 0 {
		newPart.SetInfo(typee, fit, sizeMBR, sizeNewPart, name, 1)
		if mbr.Partitions[1].Size == 0 {
			if mbr.Partitions[2].Size == 0 {
				//caso particion 4 (no existe)
				if mbr.Partitions[3].Size == 0 {
					//859 <= 1024 - 165
					if sizeNewPart <= mbr.MbrSize-sizeMBR {
						mbr.Partitions[0] = newPart
					} else {
						newPart = noPart
						fmt.Println("ERROR FDISK. Espacio insuficiente")
						return mbr, newPart, "ERROR FDISK. Espacio insuficiente \n"
					}
				} else {
					//particion 4 existe
					// 600 < 765 - 165 (600 maximo aceptado)
					if sizeNewPart <= mbr.Partitions[3].Start-sizeMBR {
						mbr.Partitions[0] = newPart
					} else {
						//Si cabe despues de 4
						newPart.SetInfo(typee, fit, mbr.Partitions[3].GetEnd(), sizeNewPart, name, 4)
						if sizeNewPart <= mbr.MbrSize-newPart.Start {
							mbr.Partitions[2] = mbr.Partitions[3]
							mbr.Partitions[3] = newPart
							//Reordeno el correlativo para que coincida con el numero de particion en que se guardo
							mbr.Partitions[2].Correlative = 3
						} else {
							newPart = noPart
							fmt.Println("ERROR FDISK. Espacio insuficiente")
							return mbr, newPart, "ERROR FDISK. Espacio insuficiente"+ "\n"
						}
					}
				}
				//Fin no existe particion 4
			} else {
				// 3 existe
				//entre mbr y 3 -> 300 <= 465 -165
				if sizeNewPart <= mbr.Partitions[2].Start-sizeMBR {
					mbr.Partitions[0] = newPart
				} else {
					//si no cabe entre el mbr y 3 debe ser despues de 3, es decir, en 4
					newPart.SetInfo(typee, fit, mbr.Partitions[2].GetEnd(), sizeNewPart, name, 4)
					if mbr.Partitions[3].Size == 0 {
						if sizeNewPart <= mbr.MbrSize-newPart.Start {
							mbr.Partitions[3] = newPart
						} else {
							newPart = noPart
							fmt.Println("ERROR FDISK. Espacio insuficiente")
							return mbr, newPart, "ERROR FDISK. Espacio insuficiente"+ "\n"
						}
					} else {
						//4 existe
						//hay espacio entre 3 y 4
						if sizeNewPart <= mbr.Partitions[3].Start-newPart.Start {
							mbr.Partitions[1] = mbr.Partitions[2]
							mbr.Partitions[2] = newPart
							//Reordenando los correlativos
							mbr.Partitions[1].Correlative = 2
							mbr.Partitions[2].Correlative = 3 //new part traia 4 y quedo en la tercer particion por eso tambien se modifica aqui
						} else if sizeNewPart <= mbr.MbrSize-mbr.Partitions[3].GetEnd() {
							//Hay espacio despues de 4
							newPart.SetInfo(typee, fit, mbr.Partitions[3].GetEnd(), sizeNewPart, name, 4)
							mbr.Partitions[1] = mbr.Partitions[2]
							mbr.Partitions[2] = mbr.Partitions[3]
							mbr.Partitions[3] = newPart
							//reconfiguro los correlativos
							mbr.Partitions[1].Correlative = 2
							mbr.Partitions[2].Correlative = 3
						} else {
							newPart = noPart
							fmt.Println("ERROR FDISK. Espacio insuficiente")
							return mbr, newPart, "ERROR FDISK. Espacio insuficiente"+ "\n"							
						}
					} //fin si hay espacio entre 3 y 4
				} //fin si no cabe antes de 3
			} //fin 3 existe
		} else {
			//2 existe
			//Si la nueva particion se puede guardar antes de 2
			if sizeNewPart <= mbr.Partitions[1].Start-sizeMBR {
				mbr.Partitions[0] = newPart
			} else {
				//Si no cabe entre mbr y 2
				//Validar si existen 3 y 4
				newPart.SetInfo(typee, fit, mbr.Partitions[1].GetEnd(), sizeNewPart, name, 3)
				if mbr.Partitions[2].Size == 0 {
					if mbr.Partitions[3].Size == 0 {
						if sizeNewPart <= mbr.MbrSize-newPart.Start {
							mbr.Partitions[2] = newPart
						} else {
							newPart = noPart
							fmt.Println("ERROR FDISK. Espacio insuficiente")
							return mbr, newPart, "ERROR FDISK. Espacio insuficiente \n"
						}
					} else {
						//4 existe (estamos entre 2 y 4)
						//62 < 69-6 (62 maximo aceptado)
						if sizeNewPart <= mbr.Partitions[3].Start-newPart.Start {
							mbr.Partitions[2] = newPart
						} else {
							//Si no cabe entre 2 y 4, ver si cabe despues de 4
							newPart.SetInfo(typee, fit, mbr.Partitions[3].GetEnd(), sizeNewPart, name, 4)
							if sizeNewPart <= mbr.MbrSize-newPart.Start { //1 <= 100-99
								mbr.Partitions[2] = mbr.Partitions[3]
								mbr.Partitions[3] = newPart
								//reordeno correlativos
								mbr.Partitions[2].Correlative = 3
							} else {
								newPart = noPart
								fmt.Println("ERROR FDISK. Espacio insuficiente")
								return mbr, newPart,"ERROR FDISK. Espacio insuficiente \n"								
							}
						} //Fin si cabe antes o despues de 4
					} //fin de 4 existe o no existe
				} else {
					//3 existe
					//entre 2 y 3
					if sizeNewPart <= mbr.Partitions[2].Start-newPart.Start {
						mbr.Partitions[0] = mbr.Partitions[1]
						mbr.Partitions[1] = newPart
						//Reordeno correlativos
						mbr.Partitions[0].Correlative = 1
						mbr.Partitions[1].Correlative = 2
					} else if mbr.Partitions[3].Size == 0 {
						//entre 3 y el final
						//cambiamos el inicio de la nueva particion porque 3 existe y no cabe antes de 3
						newPart.SetInfo(typee, fit, mbr.Partitions[2].GetEnd(), sizeNewPart, name, 4)
						if sizeNewPart <= mbr.MbrSize-newPart.Start {
							mbr.Partitions[3] = newPart
						} else {
							newPart = noPart
							fmt.Println("ERROR FDISK. Espacio insuficiente")
							return mbr, newPart, "ERROR FDISK. Espacio insuficiente"
						}
					} else {
						//si 4 existe
						//hay espacio entre 3 y 4
						newPart.SetInfo(typee, fit, mbr.Partitions[2].GetEnd(), sizeNewPart, name, 3)
						if sizeNewPart <= mbr.Partitions[3].Start-newPart.Start {
							mbr.Partitions[0] = mbr.Partitions[1]
							mbr.Partitions[1] = mbr.Partitions[2]
							mbr.Partitions[2] = newPart
							//Reordeno correlativos
							mbr.Partitions[0].Correlative = 1
							mbr.Partitions[1].Correlative = 2
						} else if sizeNewPart <= mbr.MbrSize-mbr.Partitions[3].GetEnd() {
							//entre 4 y el final
							newPart.SetInfo(typee, fit, mbr.Partitions[3].GetEnd(), sizeNewPart, name, 4)
							mbr.Partitions[0] = mbr.Partitions[1]
							mbr.Partitions[1] = mbr.Partitions[2]
							mbr.Partitions[2] = mbr.Partitions[3]
							mbr.Partitions[3] = newPart
							//Reordeno correlativos
							mbr.Partitions[0].Correlative = 1
							mbr.Partitions[1].Correlative = 2
							mbr.Partitions[2].Correlative = 3
						} else {
							newPart = noPart
							fmt.Println("ERROR FDISK. Espacio insuficiente")
							return mbr, newPart, "ERROR FDISK. Espacio insuficiente"
						}
					} //Fin si 4 existe o no (3 activa)
				} //Fin 3 existe o no existe
			} //Fin entre 2 y final (antes de 2 o depues de 2)
		} //Fin 2 existe o no existe
		//Fin de 1 no existe

		//PARTICION 2 (no existe)
	} else if mbr.Partitions[1].Size == 0 {
		//Si hay espacio entre el mbr y particion 1
		newPart.SetInfo(typee, fit, sizeMBR, sizeNewPart, name, 1)
		if sizeNewPart <= mbr.Partitions[0].Start-newPart.Start { //particion 1 ya existe (debe existir para entrar a este bloque)
			mbr.Partitions[1] = mbr.Partitions[0]
			mbr.Partitions[0] = newPart
			//Reordeno correlativo
			mbr.Partitions[1].Correlative = 2
		} else {
			//Si no hay espacio antes de particion 1
			newPart.SetInfo(typee, fit, mbr.Partitions[0].GetEnd(), sizeNewPart, name, 2) //el nuevo inicio es donde termina 1
			if mbr.Partitions[2].Size == 0 {
				if mbr.Partitions[3].Size == 0 {
					if sizeNewPart <= mbr.MbrSize-newPart.Start {
						mbr.Partitions[1] = newPart
					} else {
						newPart = noPart
						fmt.Println("ERROR FDISK. Espacio insuficiente")
						return mbr, newPart,"ERROR FDISK. Espacio insuficiente"
					}
				} else {
					//4 existe
					//entre 1 y 4
					if sizeNewPart <= mbr.Partitions[3].Start-newPart.Start {
						mbr.Partitions[1] = newPart
					} else if sizeNewPart <= mbr.MbrSize-mbr.Partitions[3].GetEnd() {
						//despues de 4
						newPart.SetInfo(typee, fit, mbr.Partitions[3].GetEnd(), sizeNewPart, name, 4)
						mbr.Partitions[2] = mbr.Partitions[3]
						mbr.Partitions[3] = newPart
						//Reordeno correlativo
						mbr.Partitions[2].Correlative = 3
					} else {
						newPart = noPart
						fmt.Println("ERROR FDISK. Espacio insuficiente")
						return mbr, newPart, "ERROR FDISK. Espacio insuficiente"
					}
				} //Fin 4 existe o no existe
			} else {
				//3 Activa
				//entre 1 y 3
				if sizeNewPart <= mbr.Partitions[2].Start-newPart.Start {
					mbr.Partitions[1] = newPart
				} else {
					//despues de 3
					newPart.SetInfo(typee, fit, mbr.Partitions[2].GetEnd(), sizeNewPart, name, 3)
					if mbr.Partitions[3].Size == 0 {
						if sizeNewPart <= mbr.MbrSize-newPart.Start {
							mbr.Partitions[3] = newPart
							//corrijo el correlativo
							mbr.Partitions[3].Correlative = 4
						} else {
							newPart = noPart
							fmt.Println("ERROR FDISK. Espacio insuficiente")
							return mbr, newPart,"ERROR FDISK. Espacio insuficiente"
						}
					} else {
						//4 existe
						//entre 3 y 4
						if sizeNewPart <= mbr.Partitions[3].Start-newPart.Start {
							mbr.Partitions[1] = mbr.Partitions[2]
							mbr.Partitions[2] = newPart
							//Corrijo el correlativo
							mbr.Partitions[1].Correlative = 2
						} else if sizeNewPart <= mbr.MbrSize-mbr.Partitions[3].GetEnd() {
							//Despues de 4
							newPart.SetInfo(typee, fit, mbr.Partitions[3].GetEnd(), sizeNewPart, name, 4)
							mbr.Partitions[1] = mbr.Partitions[2]
							mbr.Partitions[2] = mbr.Partitions[3]
							mbr.Partitions[3] = newPart
							//Corrijo los correlativos
							mbr.Partitions[1].Correlative = 2
							mbr.Partitions[2].Correlative = 3
						} else {
							newPart = noPart
							fmt.Println("ERROR FDISK. Espacio insuficiente")
							return mbr, newPart, "ERROR FDISK. Espacio insuficiente"
						}
					} //fin 4 existe o no existe
				} //Fin para entre 1 y 3, y despues de 3
			} //Fin 3 existe o no existe
		} //Fin antes o despues de particion 1
		//Fin particion 2 no existe

		//PARTICION 3
	} else if mbr.Partitions[2].Size == 0 {
		//antes de 1
		newPart.SetInfo(typee, fit, sizeMBR, sizeNewPart, name, 1)
		if sizeNewPart <= mbr.Partitions[0].Start-newPart.Start {
			mbr.Partitions[2] = mbr.Partitions[1]
			mbr.Partitions[1] = mbr.Partitions[0]
			mbr.Partitions[0] = newPart
			//Reordeno los correlativos
			mbr.Partitions[2].Correlative = 3
			mbr.Partitions[1].Correlative = 2
		} else {
			//entre 1 y 2
			newPart.SetInfo(typee, fit, mbr.Partitions[0].GetEnd(), sizeNewPart, name, 2)
			if sizeNewPart <= mbr.Partitions[1].Start-newPart.Start {
				mbr.Partitions[2] = mbr.Partitions[1]
				mbr.Partitions[1] = newPart
				//Reordeno correlativo
				mbr.Partitions[2].Correlative = 3
			} else {
				//despues de 2
				newPart.SetInfo(typee, fit, mbr.Partitions[1].GetEnd(), sizeNewPart, name, 3)
				if mbr.Partitions[3].Size == 0 {
					if sizeNewPart <= mbr.MbrSize-newPart.Start {
						mbr.Partitions[2] = newPart
					} else {
						newPart = noPart
						fmt.Println("ERROR FDISK. Espacio insuficiente")
						return mbr, newPart, "ERROR FDISK. Espacio insuficiente"
					}
				} else {
					//4 existe
					//entre 2 y 4
					if sizeNewPart <= mbr.Partitions[3].Start-newPart.Start {
						mbr.Partitions[2] = newPart
					} else if sizeNewPart <= mbr.MbrSize-mbr.Partitions[3].GetEnd() {
						//despues de 4
						newPart.SetInfo(typee, fit, mbr.Partitions[3].GetEnd(), sizeNewPart, name, 4)
						mbr.Partitions[2] = mbr.Partitions[3]
						mbr.Partitions[3] = newPart
						//Reordeno correlativo
						mbr.Partitions[2].Correlative = 3
					} else {
						newPart = noPart
						fmt.Println("ERROR FDISK. Espacio insuficiente")
						return mbr, newPart, "ERROR FDISK. Espacio insuficiente"
					}
				} //Fin de 4 existe o no existe
			} //Fin espacio entre 1 y 2 o despues de 2
		} //Fin espacio antes de 1
		//Fin particion 3

		//PARTICION 4
	} else if mbr.Partitions[3].Size == 0 {
		//antes de 1
		newPart.SetInfo(typee, fit, sizeMBR, sizeNewPart, name, 1)
		if sizeNewPart <= mbr.Partitions[0].Start-newPart.Start {
			mbr.Partitions[3] = mbr.Partitions[2]
			mbr.Partitions[2] = mbr.Partitions[1]
			mbr.Partitions[1] = mbr.Partitions[0]
			mbr.Partitions[0] = newPart
			//Reordeno los correlativos
			mbr.Partitions[3].Correlative = 4
			mbr.Partitions[2].Correlative = 3
			mbr.Partitions[1].Correlative = 2
		} else {
			//si no cabe antes de 1
			//entre 1 y 2
			newPart.SetInfo(typee, fit, mbr.Partitions[0].GetEnd(), sizeNewPart, name, 2)
			if sizeNewPart <= mbr.Partitions[1].Start-newPart.Start {
				mbr.Partitions[3] = mbr.Partitions[2]
				mbr.Partitions[2] = mbr.Partitions[1]
				mbr.Partitions[1] = newPart
				//Reordeno correlativos
				mbr.Partitions[3].Correlative = 4
				mbr.Partitions[2].Correlative = 3
			} else if sizeNewPart <= mbr.Partitions[2].Start-mbr.Partitions[1].GetEnd() {
				//entre 2 y 3
				newPart.SetInfo(typee, fit, mbr.Partitions[1].GetEnd(), sizeNewPart, name, 3)
				mbr.Partitions[3] = mbr.Partitions[2]
				mbr.Partitions[2] = newPart
				//Reordeno correlativo
				mbr.Partitions[3].Correlative = 4
			} else if sizeNewPart <= mbr.MbrSize-mbr.Partitions[2].GetEnd() {
				//despues de 3
				newPart.SetInfo(typee, fit, mbr.Partitions[2].GetEnd(), sizeNewPart, name, 4)
				mbr.Partitions[3] = newPart
			} else {
				newPart = noPart
				fmt.Println("ERROR FDISK. Espacio insuficiente")
				return mbr, newPart,"ERROR FDISK. Espacio insuficiente"
			}
		} //Fin antes y despues de 1
		//Fin particion 4
	} else {
		newPart = noPart
		fmt.Println("ERROR FDISK. Particiones primarias y/o extendidas ya no disponibles")
		return mbr, newPart,"ERROR FDISK. Particiones primarias y/o extendidas ya no disponibles"
	}

	return mbr, newPart, respuesta
}

func primerAjusteLogicas(disco *os.File, partExtend Structs.Partition, sizeNewPart int32, name string, fit string) string{
	var respuesta string
	//Se crea un ebr para cargar el ebr desde el disco y la particion extendida
	save := true //false indica que guardo en el primer ebr, true significa que debe seguir buscando
	var actual Structs.EBR
	sizeEBR := int32(binary.Size(actual)) //obtener el tamaño del ebr (el que ocupa fisicamente: 31)
	//fmt.Println("Tamaño fisico del ebr ", sizeEBR)

	//Guardo el ebr leido
	if err := Herramientas.ReadObject(disco, &actual, int64(partExtend.Start)); err != nil {
		respuesta += "ERROR FDISK Read " + err.Error()+ "\n"
		return respuesta
	}

	//NOTA: debe caber la particion con el tamaño establecido MAS su EBR
	//NOTA2: Recordar que a la hora de escribir (usar) la particion se inicia donde termina fisicamente la estructura del ebr
	//ej: si el ebr ocupa 5 bytes y la particion es de 10 bytes. los primeros 5 son del ebr entonces uso de 5-15 para escribir en el archivo el contenido de la particion

	//si el primer ebr esta vacio o no existe
	if actual.Size == 0 {
		if actual.Next == -1 {
			//validar si el tamaño de la nueva particion junto al ebr es menor al tamaño de la particion extendida
			if sizeNewPart+sizeEBR <= partExtend.Size {
				actual.SetInfo(fit, partExtend.Start, sizeNewPart, name, -1)
				if err := Herramientas.WriteObject(disco, actual, int64(actual.Start)); err != nil {
					return"ERROR FDISK Write " +err.Error()+ "\n"
				}
				save = false //ya guardo la nueva particion
				fmt.Println("Particion con nombre ", name, " creada correctamente")
				respuesta += "Particion con nombre "+ name+ " creada correctamente"+ "\n"
			} else {
				fmt.Println("ERROR FDISK. Espacio insuficiente logicas")
				return "ERROR FDISK. Espacio insuficiente logicas"+ "\n"
			}
		} else {
			//Para insertar si se elimino la primera particion (primer EBR)
			//Si actual.Next no es -1 significa que hay otra particion despues de la actual y actual.next tiene el inicio de esa particion
			disponible := actual.Next - partExtend.Start //del inicio hasta donde inicia la siguiente
			if sizeNewPart+sizeEBR <= disponible {
				actual.SetInfo(fit, partExtend.Start, sizeNewPart, name, actual.Next)
				if err := Herramientas.WriteObject(disco, actual, int64(actual.Start)); err != nil {
					return "ERROR FDISK Write " +err.Error()+ "\n"
				}
				save = false //ya guardo la nueva particion
				fmt.Println("Particion con nombre ", name, " creada correctamente")
				respuesta += "Particion con nombre " + name+ " creada correctamente"+ "\n"
			} else {
				fmt.Println("ERROR FDISK. Espacio insuficiente logicas 2")
				return "ERROR FDISK. Espacio insuficiente logicas"+ "\n"
			}
		}
		//Si esta despues del primer ebr
	}

	if save {
		//siguiente = actual.next //el valor del siguiente es el inicio de la siguiente particion
		for actual.Next != -1 {
			//si el ebr y la particion caben
			if sizeNewPart+sizeEBR <= actual.Next-actual.GetEnd() {
				break
			}
			//paso al siguiente ebr (simula un actual = actual.next)
			if err := Herramientas.ReadObject(disco, &actual, int64(actual.Next)); err != nil {
				respuesta += "ERROR FDISK Read " + err.Error()+ "\n"
				return respuesta
			}

		}

		//Despues de la ultima particion
		if actual.Next == -1 {
			//ya no es el tamaño porque ya hay espacio ocupado por lo que tomo donde termina la extendida y se resta donde termina la ultima
			if sizeNewPart+sizeEBR <= (partExtend.GetEnd() - actual.GetEnd()) {
				//guardar cambios en el ebr actual (cambio el Next)
				actual.Next = actual.GetEnd()
				if err := Herramientas.WriteObject(disco, actual, int64(actual.Start)); err != nil {
					return"ERROR FDISK Write " +err.Error()+ "\n"
				}

				//crea y guarda la nueva particion logica
				newStart := actual.GetEnd()                          //la nueva ebr inicia donde termina la ultima ebr
				actual.SetInfo(fit, newStart, sizeNewPart, name, -1) //cambia actual con los nuevos valores
				if err := Herramientas.WriteObject(disco, actual, int64(actual.Start)); err != nil {
					return "ERROR FDISK Write " +err.Error()+ "\n"
				}
				fmt.Println("Particion con nombre ", name, " creada correctamente")
				respuesta += "Particion con nombre "+ name+" creada correctamente"+ "\n"
			} else {
				fmt.Println("ERROR FDISK. Espacio insuficiente logicas 3")
				return "ERROR FDISK. Espacio insuficiente logicas"+ "\n"
			}
		} else {
			//Entre dos particiones
			if sizeNewPart+sizeEBR <= (actual.Next - actual.GetEnd()) {
				siguiente := actual.Next //guardo el siguiente de la actual para ponerlo en el siguiente de la nueva particion
				//guardar cambio de siguiente en la actual
				actual.Next = actual.GetEnd()
				if err := Herramientas.WriteObject(disco, actual, int64(actual.Start)); err != nil {
					return "ERROR FDISK Write " +err.Error()+ "\n"
				}

				//agrego la nueva particion apuntando a la siguiente de la actual
				newStart := actual.GetEnd()                                 //la nueva ebr inicia donde termina la ultima ebr
				actual.SetInfo(fit, newStart, sizeNewPart, name, siguiente) //cambia actual con los nuevos valores
				if err := Herramientas.WriteObject(disco, actual, int64(actual.Start)); err != nil {
					return "ERROR FDISK Write " +err.Error()+ "\n"
				}
				fmt.Println("Particion con nombre ", name, " creada correctamente")
				respuesta += "Particion con nombre "+ name +" creada correctamente"+ "\n"
			} else {
				fmt.Println("ERROR FDISK. Espacio insuficiente logicas 4")
				return "ERROR FDISK. Espacio insuficiente logicas "+ "\n"
			}
		}
	}
	return respuesta
}