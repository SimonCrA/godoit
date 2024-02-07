--usuario
INSERT INTO profiles."user"
(id_user, "name", email, "password", fk_id_cat_status, create_date, pwd_expiration_date, last_session, logical_delete)
VALUES
(1, 'jsalge', 'josesalge@gmail.com', '1234', 1, '2024-01-26 12:27:20.493', '2024-01-26 12:27:20.493', '2024-01-26 12:27:20.493', false);

--clave
INSERT INTO transactions.task
(id_task, fk_id_user, fk_id_cat_category, description, fk_id_cat_status, create_date, current_task_date, last_status_change, logical_delete)
VALUES(1, 1, 1, 'prueba', 1, '2024-01-28 20:56:03.516', '2024-01-28 20:56:03.516', '2024-01-28 20:56:03.516', false);
INSERT INTO transactions.task
(id_task, fk_id_user, fk_id_cat_category, description, fk_id_cat_status, create_date, current_task_date, last_status_change, logical_delete)
VALUES(2, 1, 1, 'prueba', 2, '2024-01-28 20:56:03.516', '2024-01-28 20:56:03.516', '2024-01-28 20:56:03.516', false);
INSERT INTO transactions.task
(id_task, fk_id_user, fk_id_cat_category, description, fk_id_cat_status, create_date, current_task_date, last_status_change, logical_delete)
VALUES(3, 1, 1, 'prueba', 1, '2024-01-28 20:56:03.516', '2024-01-28 20:56:03.516', '2024-01-28 20:56:03.516', true);
INSERT INTO transactions.task
(id_task, fk_id_user, fk_id_cat_category, description, fk_id_cat_status, create_date, current_task_date, last_status_change, logical_delete)
VALUES(4, 1, 1, 'prueba', 1, '2024-01-28 20:56:03.516', '2024-01-29 20:56:03.516', '2024-01-28 20:56:03.516', false);
INSERT INTO transactions.task
(id_task, fk_id_user, fk_id_cat_category, description, fk_id_cat_status, create_date, current_task_date, last_status_change, logical_delete)
VALUES(5, 1, 1, 'prueba', 2, '2024-01-28 20:56:03.516', '2024-01-29 20:56:03.516', '2024-01-28 20:56:03.516', false);

SELECT setval('transactions.task_seq',(select coalesce(max(id_task),0) from transactions.task), true);