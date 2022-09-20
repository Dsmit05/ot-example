CREATE TABLE msgs
(
    id         bigint PRIMARY KEY,
    msg text
);

create unique index msgs_id_index
    on msgs (id);