INSERT INTO remonttiv2.users
(user_id, organization_id, user_name, user_password, actions_code, rights_1, create_time, update_time)
VALUES (1, 1, 'admin', '$2a$14$adZQlMqeE3qgAgGv.25PhuREomuM.zjCVIrLdoEUCpruv5g6DKEUi', 0, 64 | 4 | 1, 0, 0),
       (2, 2, 'remontti', '$2a$14$EA3./8raO12dFE6tj/6C4evQIig3AlVRDkFuVsQJiJsAjWX7PAw2.', 0, 128 | 256 | 512, 1697057352,
        1697057352);

INSERT INTO remonttiv2.organizations
    (organization_id, organization_name, host, create_time, update_time)
VALUES (1, 'control', 'control.remontti.site', 0, 0),
       (2, 'remontti', 'work.remontti.site', 1697057352, 1697057352),
       (3, 'test', 'test.remontti.site', 0, 0);

INSERT INTO remonttiv2.navigation
    (navigation_id, title, tooltip_text, navigation_group, icon, link)
VALUES (1, 'organizations', 'organization_tooltip', 1, '/icons/organization.svg', '#/organizations'),
       (2, 'branch', 'branch_tooltip', 1, '/icons/branch.svg', '#/branches'),
       (3, 'settings', 'branch_tooltip', 1, '/icons/settings.svg', '#/settings'),
       (4, 'account', 'branch_tooltip', 1, '/icons/vue.svg', '#/account'),
       (5, 'users', 'branch_tooltip', 1, '/icons/users.svg', '#/users'),
       (6, 'money', 'branch_tooltip', 1, '/icons/wallet.svg', '#/money');

INSERT INTO remonttiv2.right_category_ids
    (category_title, category_id)
VALUES ('navigation', 1);


INSERT INTO remonttiv2.rights
    (user_id, entity_id, entity_group)
VALUES (1, 1, 1),
       (1, 3, 1),
       (1, 5, 1),
       (1, 6, 1),
       (2, 2, 1),
       (2, 5, 1),
       (2, 3, 1);


INSERT INTO remonttiv2.rights_names
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
       ('DELETE_BRANCH_LIST', pow(2, 9));

INSERT INTO remonttiv2.branches
(branch_id, organization_id, branch_name, address, phone, work_time, create_time, update_time)
VALUES (1, 8, 'plaza', 'lesnoy 47b', '79994565544', '11-19', 0, 0);


SELECT setval('remonttiv2.users_user_id_seq', 100);
SELECT setval('remonttiv2.organizations_organization_id_seq', 100);
SELECT setval('remonttiv2.navigation_navigation_id_seq', 100);
SELECT setval('remonttiv2.branches_branch_id_seq', 100);