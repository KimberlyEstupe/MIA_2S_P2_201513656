# Manejo e implementacion de arhivos
**Proyecto 1**  
__2do semestre 2024__

Creacion de un sistema de Archivos Linux con EXT3.  
Para la creacion del proyesto se utilizo Go, react y graphviz
## Para correr el proyecto
```sh
//Backend
go run main.go  

//Fronted
cd front
npm start
```
## Instalaciones

### Graphviz en Ubuntu
```sh
sudo apt-get install graphviz graphviz-dev pkg-config
sudo apt-get install python3-pip
pip install pygraphviz
sudo apt update
```
Si no funciona, reiniciar Visual Code

### Utilizar otros archivos de go
```sh
go mod init
```
Revisar go.mod y verificar el module, este sera la carpeta raiz para realizar las importaciones.  
Las funciones que se utilizaran en otros archivos deben iniciar por una letra mayuscula

### Cors
```sh
go get -u github.com/rs/cors
