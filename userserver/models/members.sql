create table members.members
(
    id          int auto_increment
        primary key,
    token       varchar(60) default ''                not null,
    username    varchar(60)                           not null,
    password    varchar(60)                           not null,
    create_time timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time timestamp(0) NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP(0),
    constraint token
        unique (token)
);

INSERT INTO members.members (id, token, username, password, create_time, update_time) VALUES (1, '8e488ab4-7f1f-46d6-bd27-ece5f0673be8', 'xiaominig', 'e10adc3949ba59abbe56e057f20f883e', '2020-04-09 05:11:41', '2019-04-09 05:11:41');
INSERT INTO members.members (id, token, username, password, create_time, update_time) VALUES (2, 'b0be0b91-9719-4361-9b59-8f3ff4d35d55', 'vector', 'e10adc3949ba59abbe56e057f20f883e', '2019-04-09 05:20:31', '2019-04-09 05:20:31');