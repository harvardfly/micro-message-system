create table gateway
(
    id          int auto_increment
        primary key,
    token       varchar(60) default ''                not null,
    im_address  varchar(60) default ''                not null,
    server_name varchar(60) default ''                not null,
    topic       varchar(60) default ''                not null,
    create_time timestamp   default CURRENT_TIMESTAMP not null,
    update_time timestamp   default CURRENT_TIMESTAMP not null,
    constraint token
        unique (token)
);

INSERT INTO gateway.gateway (id, token, im_address, server_name, topic, create_time, update_time) VALUES (1, '8e488ab4-7f1f-46d6-bd27-ece5f0673be8', '127.0.0.1:7273', 'im.server.2', 'im.server.2', '2020-04-10 05:19:42', '2019-04-10 05:19:42');
INSERT INTO gateway.gateway (id, token, im_address, server_name, topic, create_time, update_time) VALUES (6, 'b0be0b91-9719-4361-9b59-8f3ff4d35d55', '127.0.0.1:7272', 'im.server.1', 'im.server.1', '2019-04-10 05:36:31', '2019-04-10 05:36:31');