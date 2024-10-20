import { Routes, Route, HashRouter, Link } from 'react-router-dom'
import { useState } from 'react';

import hacker from '../iconos/data-entry.png';
import Consola from '../Pages/Consola';
import Discos from '../Pages/Discos';
import Partitions from '../Pages/Particiones';
import Login from '../Pages/Login';
import Explorer from '../Pages/SArchivos';
import "../Stylesheets/navbar.css"

export default function Navegador(){
    const [ ip, setIP ] = useState("localhost")
    
    //handleChange sirve para poner valor por cada cambio que detecte
    const handleChange = (e) => {
        console.log(e.target.value)
        setIP(e.target.value)
    }
    
    const logOut = (e) => {
        e.preventDefault()
        
        fetch(`http://${ip}:8080/logout`)
        .then(Response => Response.json())
        .then(rawData => {
            console.log(rawData);  
            if (rawData === 0){
                alert('sesion cerrada')
                window.location.href = '#/Discos';
            }else{
                alert('No hay sesion abierta')
            }
        }) 
        .catch(error => {
            console.error('Error en la solicitud Fetch:', error);
            // Maneja el error aquí, como mostrar un mensaje al usuario
            //alert('Error en la solicitud Fetch. Por favor, inténtalo de nuevo más tarde.');
        });
    };

    const limpiar = (e) => {
        e.preventDefault()
        console.log("limpiando")
        fetch(`http://${ip}:8080/limpiar`)
        .then(Response => Response.json())
        .then(rawData => {
            console.log(rawData); 
            if (rawData === 1){
                alert('Discos y reportes eliminados')
                window.location.href = '#/Comandos';
            }else{
                alert('Error al eliminar archiovs')
            }
        }) 
    }

    return(
        <HashRouter>
            <nav className="navbar navbar-expand-lg navbar-dark bg-dark">
                {/*COLUMNAS*/}
                
                <div className="conteiner-fluid"> 
                    <img id="imgIcon" src={hacker} alt="" width="64" height="64" className="d-inline-block align-text-top"></img>
                </div>

                <div className="conteiner"> 
                    {/*Fila 1 (titulo del proyecto, RESPALDO)*/}
                    <div className="container-fluid">
                        <a className="navbar-brand" type="submit" >
                            MIA PROYECTO 2            
                        </a>
                        {/*Cada bloque div aqui dentro es una nueva fila hacia abajo*/}
                        {/*Fila 2 (menus)*/}
                        <div className="collapse navbar-collapse" id="navbarColor02">
                            {/*ul es una lista no ordenada*/}
                            <ul className="navbar-nav me-auto">
                                {/*LISTA DE MENUS QUE ESTARAN EN LA BARRA DE NAVEGACION*/}
                                <li className="nav-item">
                                    {/* Utiliza Link en lugar de a para navegar entre rutas */}
                                    <Link className="nav-link active" to="/Consola">Comandos</Link>
                                </li>

                                <li className="nav-item">
                                    {/* Enlaza primero a discos porque el flujo es empezar por discos luego particiones y luego el sistema de archivos */}
                                    <Link className="nav-link" to="/Login">Explorador</Link>
                                </li>

                                <li className="nav-item">
                                    <button id="btnLogOut" onClick={logOut} className="nav-link">Logout</button>
                                </li>
                            </ul>{/*Fin de lista de menus*/}
                        </div>{/*Fila de menus en la barra de navegacion*/}
                    </div>{/*Fila Titulo*/}
                </div>{/*Cierro tercer columna (Menu)*/}
                <input id="InIP" className="form-control me-2 mx-auto" style={{ maxWidth: "180px" }} placeholder="IP" onChange={handleChange}/>
                
            </nav> 
            
            {/* 
            <Routes>
                <Route path="/" element ={<Consola newIp={ip}/>}/> 
                <Route path="/Consola" element ={<Consola newIp={ip}/>}/> 
                <Route path="/Discos" element ={<Discos newIp={ip}/>}/> 
                <Route path="/Disco/:id" element ={<Partitions newIp={ip}/>}/> 
                <Route path="/Login/:disk/:part" element ={<Login newIp={ip}/>}/>
                <Route path="/Explorador/:id" element ={<Explorer newIp={ip}/>}/>              
            </Routes>
            
            */}
            <Routes>
                <Route path="/" element ={<Consola newIp={ip}/>}/> {/*home*/}
                <Route path="/Consola" element ={<Consola newIp={ip}/>}/> 
                <Route path="/Login" element ={<Login newIp={ip}/>}/>
                <Route path="/Discos/:id" element ={<Discos newIp={ip}/>}/> 
                
                <Route path="/Explorador/:id" element ={<Explorer newIp={ip}/>}/>              
            </Routes>
        </HashRouter>
    );
}