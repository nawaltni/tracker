-- migrate:up
alter table user_positions add column backend_user_id varchar(255) not null;
create index user_positions_backend_user_id_index on user_positions(backend_user_id);

-- migrate:down
alter table user_positions drop column backend_user_id;
drop index user_positions_backend_user_id_index;
