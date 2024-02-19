INSERT INTO public.cat_categories
(id, created_at, updated_at, deleted_at, "name", description, logical_delete)
VALUES(1, '2024-02-07 23:02:00.191', '2024-02-07 23:02:00.191', NULL, 'General', 'Categoria general', false);


INSERT INTO public.cat_statuses
(id, created_at, updated_at, deleted_at, "name", description, logical_delete)
VALUES(1, '2024-02-07 23:08:41.723', '2024-02-07 23:08:41.723', NULL, 'Pendiente', 'Tarea pendiente', false);
INSERT INTO public.cat_statuses
(id, created_at, updated_at, deleted_at, "name", description, logical_delete)
VALUES(2, '2024-02-07 23:08:41.723', '2024-02-07 23:08:41.723', NULL, 'Realizada', 'Tarea realizada', false);
INSERT INTO public.cat_statuses
(id, created_at, updated_at, deleted_at, "name", description, logical_delete)
VALUES(3, '2024-02-07 23:08:41.723', '2024-02-07 23:08:41.723', NULL, 'Activo', 'Usuario activo', false);
INSERT INTO public.cat_statuses
(id, created_at, updated_at, deleted_at, "name", description, logical_delete)
VALUES(4, '2024-02-07 23:08:41.723', '2024-02-07 23:08:41.723', NULL, 'Inactivo', 'Usuario inactivo', false);


INSERT INTO public.users
(id, created_at, updated_at, deleted_at, "name", email, "password", password_expiration_date, fk_id_cat_status, last_session, logical_delete)
VALUES(4, '2024-02-07 23:16:40.012', '2024-02-07 23:16:40.012', NULL, 'simon', 'simoncra@gmail.com', '$2a$14$I8eVrIZipomAaKhjVQfVOOa71r/kc/MA.XtFtxjLRHl7U5B3Zxddu', '2024-02-07 23:16:40.009', 4, '2024-02-07 23:16:40.009', false);
