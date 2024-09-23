package administradordiscos

import (
	"fmt"
	"strings"
)

func Mkdisk(entrada []string)  {
	/*
	Se recorren todos los parametros
	_ seria el indice, pero se omite. 
	El [1:] indica que se inicializa en el primer parametro de mkdisk
	*/
	for _,parametro :=range entrada[1:]{
		tmp := strings.TrimRight(parametro," ")
		//Dividir los parametros entre parametro y valor
		valores := strings.Split(tmp,"=")

		if len(valores)!=2{
			fmt.Println("ERROR MKDIS, valor desconocido de parametros ",valores[1])
			//return "ERROR MKDIS, valor desconocido de parametros "+valores[1]
		}

	}
}