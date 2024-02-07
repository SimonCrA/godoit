 ## App de tareas para procastinadores excepcionales

### Pasos para correr el proyecto

1. Clonar el repositorio

2. ir al directorio del proyecto
```
cd godoit

```

3. ejecutar el comando 

```
go mod init
go mod tidy

```
4. crear un contenedor con la base de datos


```
sudo docker run -d --name my_postgres -e "POSTGRES_PASSWORD=12345678" -e  "POSTGRES_USER=docker" -e "POSTGRES_DB=docker-db" -p 5433:5432 postgres

```

5. crear un archivo .env y copiar el contenido de .env.template
