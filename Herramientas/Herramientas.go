package Herramientas

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

func CrearDisco(path string) error {
	//asegurar que exista la ruta (el directorio) creando la ruta
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		fmt.Println("Error al crear el disco, path: ", err)
		return err
	}

	//crear el archivo si aun no existe
	if _, err := os.Stat(path); os.IsNotExist(err) {
		newFile, err := os.Create(path)
		if err != nil {
			fmt.Println("Error al crear el disco: ", err)
			return err
		}
		defer newFile.Close()
	}
	return nil
}

func OpenFile(name string) (*os.File, error) {
	file, err := os.OpenFile(name, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Error OpenFile ==", err)
		return nil, err
	}
	return file, nil
}

// Function to Write an object in a bin file
func WriteObject(file *os.File, data interface{}, position int64) error {
	file.Seek(position, 0) //(posicion, desde donde) -> (5,0) significa a la posicion 5 desde el inicio del archivo
	err := binary.Write(file, binary.LittleEndian, data)
	if err != nil {
		fmt.Println("Err WriteObject==", err)
		return err
	}
	return nil
}

// Function to Read an object from a bin file
func ReadObject(file *os.File, data interface{}, position int64) error {
	file.Seek(position, 0)
	err := binary.Read(file, binary.LittleEndian, data)
	if err != nil {
		fmt.Println("Err ReadObject==", err)
		return err
	}
	return nil
}

// para eliminar en el archivo una particion logica
func DelPartL(size int32) []byte {
	datos := make([]byte, size)
	return datos
}

//Elimina caracteres Ilegibles de una cadena de entrada
func EliminartIlegibles(entrada string) string{
	// Función de transformación que elimina caracteres no legibles
	transformFunc := func(r rune) rune {
		//unicode.IsPrint indica si es legible o no.
		//si el caracter se puede leer, lo regresa, de lo contrario devuekve -1
		if unicode.IsPrint(r) {
			return r
		}
		return -1
	}

	// Aplicar la función de transformación a la cadena de entrada
	salida := strings.Map(transformFunc, entrada)
	return salida	
}

// probar la escritura de la particion logica
func EscribirPartL(size int32) string {
	cad := strings.Repeat("L", int(size))
	return cad
}