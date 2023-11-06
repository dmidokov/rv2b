package db

var TablesStructSQL = `
	CREATE TABLE IF NOT EXISTS remonttiv2.users(
		user_id SERIAL NOT NULL,
		organization_id integer NOT NULL,
		user_name character varying(50) NOT NULL,
		user_password character varying(350) NOT NULL,
		actions_code integer,
		rights_1 integer NOT NULL,
		create_time integer NOT NULL,
		update_time integer NOT NULL,
		PRIMARY KEY (user_id),
		UNIQUE (organization_id, user_name)
		)
	    TABLESPACE pg_default;

	CREATE TABLE IF NOT EXISTS remonttiv2.organizations(
	    organization_id SERIAL NOT NULL,
	    organization_name character varying(50) NOT NULL,
	    host character varying(100) NOT NULL,
		create_time integer NOT NULL,
		update_time integer NOT NULL,
		creator integer NOT NULL                           
	)
	TABLESPACE pg_default;

	CREATE TABLE IF NOT EXISTS remonttiv2.navigation(
	    navigation_id  SERIAL NOT NULL,
	    title character varying(50) NOT NULL,
	    tooltip_text character varying(50) NOT NULL,
	    navigation_group integer,
	    icon character varying(30),
	    link character varying(200)
	)
	TABLESPACE pg_default;

	CREATE TABLE IF NOT EXISTS remonttiv2.right_category_ids(
	    category_title character varying(50) NOT NULL,
	    category_id integer NOT NULL
	)
	TABLESPACE pg_default;

	CREATE TABLE IF NOT EXISTS remonttiv2.rights(
	    user_id  integer NOT NULL,
	    entity_id integer NOT NULL,
	    entity_group integer NOT NULL
	)
	TABLESPACE pg_default;

	CREATE TABLE IF NOT EXISTS remonttiv2.rights_names(
	    right_id SERIAL NOT NULL,
		name character varying(30) NOT NULL,
		value integer                                      
	)

	TABLESPACE pg_default;

	CREATE TABLE IF NOT EXISTS remonttiv2.branches(
		branch_id SERIAL NOT NULL,
		organization_id integer NOT NULL,
		branch_name character varying(50) NOT NULL,
		address character varying(350) NOT NULL,
		phone character varying(20) NOT NULL,
		work_time character varying(20) NOT NULL,
		create_time integer NOT NULL,
		update_time integer NOT NULL
		)
	    TABLESPACE pg_default;
`

var TablesDataSQL = `
	
	INSERT INTO remonttiv2.users
	    (user_id, organization_id, user_name, user_password, actions_code, rights_1, create_time, update_time) 
	VALUES 
	    (1,0,'admin','$2a$14$adZQlMqeE3qgAgGv.25PhuREomuM.zjCVIrLdoEUCpruv5g6DKEUi',0,0,0,0),
	    (2,0,'admin1','$2a$14$adZQlMqeE3qgAgGv.25PhuREomuM.zjCVIrLdoEUCpruv5g6DKEUi',0,0,0,0);


	/* ORGANIZATIONS  */

	INSERT INTO remonttiv2.organizations
		(organization_id, organization_name, host, create_time, update_time, creator) 
	VALUES
		(0, 'control', 'control.remontti.site', 0, 0); 

	/* NAVIGATION */
	
	INSERT INTO remonttiv2.navigation
	    (navigation_id,title,tooltip_text,navigation_group,icon)
	VALUES 
		(0, 'organizations', 'organization_tooltip', 1, '/icons/organization.svg'),
		(1, 'branch', 'branch_tooltip', 1, '/icons/branch.svg');
-- 	INSERT INTO remonttiv2.navigation
-- 	    (navigation_id,title,tooltip_text,navigation_group,icon)
-- 	VALUES 
-- 		(1, 'organizations', 'organization_tooltip', 1, '');
	
`
var DropSchemaSQL = `DROP SCHEMA remonttiv2 CASCADE; CREATE SCHEMA remonttiv2;`
