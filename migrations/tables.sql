-- DROP SCHEMA IF EXISTS remonttiv2 CASCADE;
CREATE SCHEMA IF NOT EXISTS RV2_DB_NAME;

-- user_type
--     0 - Сервисная/Управляющая учетка
--     1 - Учетка сотрудников
CREATE TABLE IF NOT EXISTS RV2_DB_NAME.users
(
    user_id         SERIAL                 NOT NULL,
    organization_id integer                NOT NULL,
    user_name       character varying(50)  NOT NULL,
    user_password   character varying(350) NOT NULL,
    actions_code    integer,
    rights_1        integer                NOT NULL,
    create_time     integer                NOT NULL,
    update_time     integer                NOT NULL,
    account_icon    character varying(100) DEFAULT '/icons/account.svg',
    user_type       integer                DEFAULT 0,
    start_page      character varying(100) DEFAULT '/',
    PRIMARY KEY (user_id),
    UNIQUE (organization_id, user_name)
)
    TABLESPACE pg_default;

CREATE TABLE IF NOT EXISTS RV2_DB_NAME.organizations
(
    organization_id   SERIAL                 NOT NULL,
    organization_name character varying(50)  NOT NULL,
    host              character varying(100) NOT NULL,
    create_time       integer                NOT NULL,
    update_time       integer                NOT NULL,
    creator           integer                NOT NULL
)
    TABLESPACE pg_default;

CREATE TABLE IF NOT EXISTS RV2_DB_NAME.navigation
(
    navigation_id    SERIAL                NOT NULL,
    title            character varying(50) NOT NULL,
    tooltip_text     character varying(50) NOT NULL,
    navigation_group integer,
    icon             character varying(30),
    link             character varying(200)
)
    TABLESPACE pg_default;

CREATE TABLE IF NOT EXISTS RV2_DB_NAME.right_category_ids
(
    category_title character varying(50) NOT NULL,
    category_id    integer               NOT NULL
)
    TABLESPACE pg_default;

CREATE TABLE IF NOT EXISTS RV2_DB_NAME.rights
(
    user_id      integer NOT NULL,
    entity_id    integer NOT NULL,
    entity_group integer NOT NULL
)
    TABLESPACE pg_default;

CREATE TABLE IF NOT EXISTS RV2_DB_NAME.rights_names
(
    right_id SERIAL                NOT NULL,
    name     character varying(30) NOT NULL,
    value    integer
)
    TABLESPACE pg_default;

CREATE TABLE IF NOT EXISTS RV2_DB_NAME.branches
(
    branch_id       SERIAL                 NOT NULL,
    organization_id integer                NOT NULL,
    branch_name     character varying(50)  NOT NULL,
    address         character varying(350) NOT NULL,
    phone           character varying(20)  NOT NULL,
    work_time       character varying(20)  NOT NULL,
    create_time     integer                NOT NULL,
    update_time     integer                NOT NULL
)
    TABLESPACE pg_default;

CREATE TABLE IF NOT EXISTS RV2_DB_NAME.user_branches
(
    user_id   integer NOT NULL,
    branch_id integer NOT NULL
)
    TABLESPACE pg_default;

CREATE TABLE IF NOT EXISTS RV2_DB_NAME.users_create_relations
(
    creator_id integer NOT NULL,
    created_id integer NOT NULL
)
    TABLESPACE pg_default;

CREATE TABLE IF NOT EXISTS RV2_DB_NAME.entity_group_to_entity_name
(
    group_id    integer NOT NULL,
    entity_name character varying(40)
)
    TABLESPACE pg_default;

CREATE TABLE IF NOT EXISTS RV2_DB_NAME.hot_switch_relations
(
    from_user integer NOT NULL,
    to_user   integer NOT NULL
)
    TABLESPACE pg_default;








INSERT INTO RV2_DB_NAME.users (user_id, organization_id, user_name, user_password, actions_code, rights_1, create_time, update_time, account_icon, user_type, start_page) VALUES (1, 1, 'admin', '$2a$14$adZQlMqeE3qgAgGv.25PhuREomuM.zjCVIrLdoEUCpruv5g6DKEUi', 0, 2147483647, 0, 0, '', 0, '/');
INSERT INTO RV2_DB_NAME.users (user_id, organization_id, user_name, user_password, actions_code, rights_1, create_time, update_time, account_icon, user_type, start_page) VALUES (3, 2, 'employee', '$2a$14$NorbQj63.CseGzjcLs4p9.8zFIDusG/ZIAMMua9j3RV1aSvDjfj26', 0, 128, 1704540255, 1704540255, '/icons/img.png', 0, '#/branchselector');
INSERT INTO RV2_DB_NAME.users (user_id, organization_id, user_name, user_password, actions_code, rights_1, create_time, update_time, account_icon, user_type, start_page) VALUES (4, 2, 'l.markova', '$2a$14$Exm9En.WXP7.n1n2dLPP...OftXwsH8kDNzPVKV4Rte207a/u9rcG', 0, 0, 1706374247, 1706374247, '/icons/account.svg', 1, '/');
INSERT INTO RV2_DB_NAME.users (user_id, organization_id, user_name, user_password, actions_code, rights_1, create_time, update_time, account_icon, user_type, start_page) VALUES (5, 2, 't.bespalova', '$2a$14$l.Yh7el.4SkW5DrI7RJjCOAVnTwK5y43t54jSrWyXGp/O3mrjayzq', 0, 0, 1706374258, 1706374258, '/icons/account.svg', 1, '/');
INSERT INTO RV2_DB_NAME.users (user_id, organization_id, user_name, user_password, actions_code, rights_1, create_time, update_time, account_icon, user_type, start_page) VALUES (6, 2, 'r.tagirova', '$2a$14$qEdu6tI.1TAHdWttlnxhCeF9VQsyw1ACA5iPn5Z2NUq3hM32lYARS', 0, 0, 1706374271, 1706374271, '/icons/account.svg', 1, '/');
INSERT INTO RV2_DB_NAME.users (user_id, organization_id, user_name, user_password, actions_code, rights_1, create_time, update_time, account_icon, user_type, start_page) VALUES (2, 2, 'remontti', '$2a$14$EA3./8raO12dFE6tj/6C4evQIig3AlVRDkFuVsQJiJsAjWX7PAw2.', 0, 8085, 1697057352, 1706374309, '/icons/upload-2287900351.png', 0, '/');







INSERT INTO RV2_DB_NAME.organizations
(organization_id, organization_name, host, create_time, update_time, creator)
VALUES (1, 'control', 'control.remontti.site', 0, 0, 1),
       (2, 'remontti', 'work.remontti.site', 1697057352, 1697057352, 1),
       (3, 'test', 'test.RV2_DOMAIN_NAME', 0, 0, 1);





INSERT INTO RV2_DB_NAME.navigation
    (navigation_id, title, tooltip_text, navigation_group, icon, link)
VALUES (1, 'organizations', 'organization_tooltip', 1, '/icons/organization.svg', '#/organizations'),
       (2, 'branch', 'branch_tooltip', 1, '/icons/branch-two.svg', '#/branches'),
       (3, 'settings', 'branch_tooltip', 1, '/icons/settings.svg', '#/settings'),
       (4, 'account', 'branch_tooltip', 1, '/icons/account.svg', '#/account'),
       (5, 'users', 'branch_tooltip', 1, '/icons/users.svg', '#/users'),
       (6, 'money', 'branch_tooltip', 1, '/icons/wallet.svg', '#/money');







INSERT INTO RV2_DB_NAME.right_category_ids
    (category_title, category_id)
VALUES ('navigation', 1);








INSERT INTO RV2_DB_NAME.rights (user_id, entity_id, entity_group) VALUES (1, 1, 1);
INSERT INTO RV2_DB_NAME.rights (user_id, entity_id, entity_group) VALUES (1, 3, 1);
INSERT INTO RV2_DB_NAME.rights (user_id, entity_id, entity_group) VALUES (1, 5, 1);
INSERT INTO RV2_DB_NAME.rights (user_id, entity_id, entity_group) VALUES (1, 6, 1);
INSERT INTO RV2_DB_NAME.rights (user_id, entity_id, entity_group) VALUES (1, 4, 1);
INSERT INTO RV2_DB_NAME.rights (user_id, entity_id, entity_group) VALUES (1, 2, 1);
INSERT INTO RV2_DB_NAME.rights (user_id, entity_id, entity_group) VALUES (2, 2, 1);
INSERT INTO RV2_DB_NAME.rights (user_id, entity_id, entity_group) VALUES (2, 4, 1);
INSERT INTO RV2_DB_NAME.rights (user_id, entity_id, entity_group) VALUES (2, 5, 1);
INSERT INTO RV2_DB_NAME.rights (user_id, entity_id, entity_group) VALUES (2, 3, 1);







INSERT INTO RV2_DB_NAME.rights_names
    (name, value)
VALUES ('ADD_USER', pow(2, 0)),
       ('EDIT_USER', pow(2, 1)),
       ('DELETE_USER', pow(2, 2)),
       ('ADD_ORGANIZATION', pow(2, 3)),
       ('EDIT_ORGANIZATION', pow(2, 4)),
       ('DELETE_ORGANIZATION', pow(2, 5)),
       ('VIEW_ORGANIZATION_LIST', pow(2, 6)),
       ('VIEW_BRANCH_LIST', pow(2, 7)),
       ('CREATE_BRANCH_LIST', pow(2, 8)),
       ('DELETE_BRANCH_LIST', pow(2, 9)),
       ('EDIT_USER_RIGHTS', pow(2, 10)),
       ('EDIT_USER_NAVIGATION', pow(2, 11)),
       ('EDIT_USER_HOT_SWITCH', pow(2, 12));






INSERT INTO RV2_DB_NAME.branches
(branch_id, organization_id, branch_name, address, phone, work_time, create_time, update_time)
VALUES (1, 2, 'plaza', 'lesnoy 47b', '79994565544', '11-19', 0, 0);
INSERT INTO RV2_DB_NAME.branches (branch_id, organization_id, branch_name, address, phone, work_time, create_time,
                                 update_time)
VALUES (2, 2, 'master', 'Ленина 20', '71234567890', '11-19', 1704535350, 1704535350);
INSERT INTO RV2_DB_NAME.branches (branch_id, organization_id, branch_name, address, phone, work_time, create_time,
                                 update_time)
VALUES (3, 2, 'maxi', 'Ленина 14', '78846464466', '10-21', 1704669486, 1704669486);






INSERT INTO RV2_DB_NAME.user_branches (user_id, branch_id)
VALUES (3, 1);
INSERT INTO RV2_DB_NAME.user_branches (user_id, branch_id)
VALUES (3, 2);
INSERT INTO RV2_DB_NAME.user_branches (user_id, branch_id)
VALUES (3, 3);
INSERT INTO RV2_DB_NAME.user_branches (user_id, branch_id)
VALUES (2, 1);
INSERT INTO RV2_DB_NAME.user_branches (user_id, branch_id)
VALUES (2, 2);
INSERT INTO RV2_DB_NAME.user_branches (user_id, branch_id)
VALUES (2, 3);





INSERT INTO RV2_DB_NAME.users_create_relations (creator_id, created_id) VALUES (1, 2);
INSERT INTO RV2_DB_NAME.users_create_relations (creator_id, created_id) VALUES (2, 3);
INSERT INTO RV2_DB_NAME.users_create_relations (creator_id, created_id) VALUES (0, 1);
INSERT INTO RV2_DB_NAME.users_create_relations (creator_id, created_id) VALUES (2, 4);
INSERT INTO RV2_DB_NAME.users_create_relations (creator_id, created_id) VALUES (2, 5);
INSERT INTO RV2_DB_NAME.users_create_relations (creator_id, created_id) VALUES (2, 6);



INSERT INTO RV2_DB_NAME.hot_switch_relations (from_user, to_user) VALUES (3, 4);
INSERT INTO RV2_DB_NAME.hot_switch_relations (from_user, to_user) VALUES (3, 5);



SELECT setval('RV2_DB_NAME.users_user_id_seq', 100);
SELECT setval('RV2_DB_NAME.organizations_organization_id_seq', 100);
SELECT setval('RV2_DB_NAME.navigation_navigation_id_seq', 100);
SELECT setval('RV2_DB_NAME.branches_branch_id_seq', 100);

