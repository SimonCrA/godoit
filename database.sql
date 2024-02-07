--*--CREAR ROL

--jsalge
/*
DROP ROLE IF EXISTS jsalge;

CREATE ROLE jsalge WITH 
	SUPERUSER
	NOCREATEDB
	NOCREATEROLE
	INHERIT
	LOGIN
	PASSWORD '1234'
	NOREPLICATION
	NOBYPASSRLS;
*/

--g_developer

DROP ROLE IF EXISTS g_developer;

CREATE ROLE g_developer WITH 
	NOSUPERUSER
	NOCREATEDB
	NOCREATEROLE
	INHERIT
	NOLOGIN
	NOREPLICATION
	NOBYPASSRLS;

--sclemente

DROP ROLE IF EXISTS sclemente;

CREATE ROLE sclemente WITH 
	NOSUPERUSER
	NOCREATEDB
	NOCREATEROLE
	INHERIT
	LOGIN
	PASSWORD '4321'
	NOREPLICATION
	NOBYPASSRLS;

GRANT g_developer TO sclemente;

--c_op

DROP ROLE IF EXISTS c_op;

CREATE ROLE c_op WITH 
	NOSUPERUSER
	NOCREATEDB
	NOCREATEROLE
	INHERIT
	LOGIN
	PASSWORD '4321'
	NOREPLICATION
	NOBYPASSRLS;

GRANT g_developer TO c_op;

--*--CREAR SCHEMAS

CREATE SCHEMA profiles;
CREATE SCHEMA transactions;

GRANT ALL ON SCHEMA public TO g_developer;
GRANT ALL ON SCHEMA profiles TO g_developer;
GRANT ALL ON SCHEMA transactions TO g_developer;

--privilegios por default

--public
ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA public
GRANT INSERT, SELECT, UPDATE, DELETE, REFERENCES, TRIGGER ON TABLES TO g_developer;

ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA public
GRANT ALL ON SEQUENCES TO g_developer;

ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA public
GRANT EXECUTE ON FUNCTIONS TO g_developer;

ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA public
GRANT USAGE ON TYPES TO g_developer;

--profiles
ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA profiles
GRANT INSERT, SELECT, UPDATE, DELETE, REFERENCES, TRIGGER ON TABLES TO g_developer;

ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA profiles
GRANT ALL ON SEQUENCES TO g_developer;

ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA profiles
GRANT EXECUTE ON FUNCTIONS TO g_developer;

ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA profiles
GRANT USAGE ON TYPES TO g_developer;

--transactions
ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA transactions
GRANT INSERT, SELECT, UPDATE, DELETE, REFERENCES, TRIGGER ON TABLES TO g_developer;

ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA transactions
GRANT ALL ON SEQUENCES TO g_developer;

ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA transactions
GRANT EXECUTE ON FUNCTIONS TO g_developer;

ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA transactions
GRANT USAGE ON TYPES TO g_developer;


--*--CREAR TABLA STATUS

CREATE TABLE IF NOT EXISTS public.cat_status(

	id_cat_status BIGINT NOT NULL,
    name VARCHAR NOT NULL,
	description VARCHAR NOT NULL,
	create_date timestamptz NOT NULL DEFAULT NOW(),
    logical_delete boolean NOT NULL DEFAULT false,
	    CONSTRAINT pk_cat_status PRIMARY KEY (id_cat_status),
	    CONSTRAINT uk_status_name UNIQUE (name)
);

ALTER TABLE IF EXISTS profiles.cat_status
    OWNER to postgres;
        
-- Comentarios para la tabla

COMMENT ON TABLE public.cat_status IS 'Tabla que registra los estatus';

COMMENT ON COLUMN public.cat_status.id_cat_status IS 'Clave Primaria';

COMMENT ON COLUMN public.cat_status.name IS 'Nombre del estatus';

COMMENT ON COLUMN public.cat_status.description IS 'Descripción del estatus';

COMMENT ON COLUMN public.cat_status.create_date IS 'Fecha de creación';

COMMENT ON COLUMN public.cat_status.logical_delete IS 'Borrado logico';

--Insert
INSERT INTO public.cat_status 
(id_cat_status, "name", description) VALUES
(1,'Pendiente','Tarea pendiente'),
(2,'Realizada','Tarea realizada');

--*--CREAR TABLA DE CATEGORIA DE TAREAS

CREATE TABLE IF NOT EXISTS public.cat_category(

	id_cat_category BIGINT NOT NULL,
    name VARCHAR NOT NULL,
	description VARCHAR NOT NULL,
	create_date timestamptz NOT NULL DEFAULT NOW(),
    logical_delete boolean NOT NULL DEFAULT false,
	    CONSTRAINT pk_cat_category PRIMARY KEY (id_cat_category),
	    CONSTRAINT uk_category_name UNIQUE (name)
);

ALTER TABLE IF EXISTS profiles.cat_category
    OWNER to postgres;
        
-- Comentarios para la tabla

COMMENT ON TABLE public.cat_category IS 'Tabla que registra las categorias';

COMMENT ON COLUMN public.cat_category.id_cat_category IS 'Clave Primaria';

COMMENT ON COLUMN public.cat_category.name IS 'Nombre de la categoria';

COMMENT ON COLUMN public.cat_category.description IS 'Descripción de la categoria';

COMMENT ON COLUMN public.cat_category.create_date IS 'Fecha de creación';

COMMENT ON COLUMN public.cat_category.logical_delete IS 'Borrado logico';

--Insert
INSERT INTO public.cat_category 
(id_cat_category, "name", description) VALUES
(1,'General','Categoria general');

/*
INSERT INTO public.cat_category
(id_cat_category, "name", description)
VALUES
(1, 'Hoy', 'Tarea de hoy');
(2, 'Mañána', 'Tarea de mañana');
*/


--*--CREAR TABLA DE USERS

--- Crear secuencia

CREATE SEQUENCE IF NOT EXISTS profiles.user_seq
    INCREMENT 1
    START 1
    MINVALUE 0
    MAXVALUE 9223372036854775807
    CACHE 1;

-- Comentarios

    COMMENT ON SEQUENCE profiles.user_seq IS 'Secuencia para la tabla users';

ALTER SEQUENCE profiles.user_seq
    OWNER TO postgres;


-- Crear tabla

CREATE TABLE IF NOT EXISTS profiles.user(

	id_user BIGINT NOT NULL DEFAULT nextval('profiles.user_seq'),
    name VARCHAR NOT NULL,
	email VARCHAR NOT NULL,
	password VARCHAR NOT NULL,
	fk_id_cat_status bigint NOT NULL,
	create_date timestamptz NOT NULL DEFAULT NOW(),
	pwd_expiration_date timestamptz NOT NULL DEFAULT NOW(),
	last_session timestamptz NOT NULL DEFAULT NOW(),
    logical_delete boolean NOT NULL DEFAULT false,
	    CONSTRAINT pk_user PRIMARY KEY (id_user),
	    CONSTRAINT fk_cat_status FOREIGN KEY (fk_id_cat_status) REFERENCES public.cat_status(id_cat_status)
);

ALTER TABLE IF EXISTS profiles.user
    OWNER to postgres;
        
-- Comentarios para la tabla

COMMENT ON TABLE profiles.user IS 'Tabla que registra los usuarios';

COMMENT ON COLUMN profiles.user.id_user IS 'Clave Primaria';

COMMENT ON COLUMN profiles.user.name IS 'Nombre del usuario';

COMMENT ON COLUMN profiles.user.email IS 'Correo electronico del usuario';

COMMENT ON COLUMN profiles.user.password IS 'Contraseña del usuario';

COMMENT ON COLUMN profiles.user.fk_id_cat_status IS 'Estatus del usuario';

COMMENT ON COLUMN profiles.user.create_date IS 'Fecha de creación';

COMMENT ON COLUMN profiles.user.pwd_expiration_date IS 'Fecha de expiración de contraseña';

COMMENT ON COLUMN profiles.user.last_session IS 'Fecha de ultima sesion del usuario';

COMMENT ON COLUMN profiles.user.logical_delete IS 'Borrado logico';
   
--*-- Indice para miebro y pregunta de seguridad 
CREATE UNIQUE INDEX idx_uk_user ON profiles.user USING btree (id_user, logical_delete) WHERE (logical_delete = false);


--*--CREAR TABLA DE TAREAS

-- Crear secuencia

CREATE SEQUENCE IF NOT EXISTS transactions.task_seq
    INCREMENT 1
    START 1
    MINVALUE 0
    MAXVALUE 9223372036854775807
    CACHE 1;

-- Comentarios

    COMMENT ON SEQUENCE transactions.task_seq IS 'Secuencia para la tabla task';

ALTER SEQUENCE transactions.task_seq
    OWNER TO postgres;


-- Crear tabla

CREATE TABLE IF NOT EXISTS transactions.task(

	id_task BIGINT NOT NULL DEFAULT nextval('transactions.task_seq'),
	fk_id_user BIGINT NOT NULL,
    fk_id_cat_category BIGINT NOT NULL,
	description VARCHAR NOT NULL,
	fk_id_cat_status bigint NOT NULL DEFAULT 1,
	current_task_date timestamptz NOT NULL DEFAULT NOW(),
	last_status_change timestamptz NOT NULL DEFAULT NOW(),
	create_date timestamptz NOT NULL DEFAULT NOW(),
    logical_delete boolean NOT NULL DEFAULT false,
	    CONSTRAINT pk_task PRIMARY KEY (id_task),
		CONSTRAINT fk_user FOREIGN KEY (fk_id_user) REFERENCES profiles.user(id_user),
		CONSTRAINT fk_cat_category FOREIGN KEY (fk_id_cat_category) REFERENCES public.cat_category(id_cat_category),
	    CONSTRAINT fk_cat_status FOREIGN KEY (fk_id_cat_status) REFERENCES public.cat_status(id_cat_status)
);

ALTER TABLE IF EXISTS transactions.user
    OWNER to postgres;
        
-- Comentarios para la tabla transactions.user

COMMENT ON TABLE transactions.task IS 'Tabla que registra las tareas';

COMMENT ON COLUMN transactions.task.id_task IS 'Clave Primaria';

COMMENT ON COLUMN transactions.task.fk_id_user IS 'Usuario al que esta asociada la tarea';

COMMENT ON COLUMN transactions.task.fk_id_cat_category IS 'Categoria de la tarea';

COMMENT ON COLUMN transactions.task.description IS 'Descripción de la tarea';

COMMENT ON COLUMN transactions.task.fk_id_cat_status IS 'Estatus de la tarea';

COMMENT ON COLUMN transactions.task.last_status_change IS 'Fecha actual de la tarea';

COMMENT ON COLUMN transactions.task.last_status_change IS 'Fecha de ultimo cambio de estatus';

COMMENT ON COLUMN transactions.task.create_date IS 'Fecha de creación';

COMMENT ON COLUMN transactions.task.logical_delete IS 'Borrado logico';
