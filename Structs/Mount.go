package Structs

/*Almacena la informacion de los Discos montados:
Se asigna una letra a cada disco montado y 
va sumando 1 cada vez que se monta otra particion en dicho disco
*/
var Pmontaje []DMontado
type DMontado struct{
	MPath  string	//Path del Disco
	Letter byte		//Letra que se le asigna
	Cont   int 		//COntador del numero de particion montada
}

//Ingresa la informacion al Struct
func AddPathM (path string, L byte, cont int){
	Pmontaje = append(Pmontaje, DMontado{MPath: path ,Letter: L,Cont: cont})
}

// ==============================================================================

//Almacena la informacion de cada Id junto a su Path

var Montadas []mountAlready
type mountAlready struct{
	Id		 string	//Id de la particion
	PathM	 string	//Path del disco al que pertenece la particion
}

//Ingresa la informacion al Struct
func AddMontadas(id string, path string){
	Montadas = append(Montadas, mountAlready{Id: id, PathM: path})
}