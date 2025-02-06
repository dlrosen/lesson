-- SQLite
--SELECT school_name||'-'||school_id FROM school ORDER BY school_id;
--select * from school;
--select * from time_period;
--select * from student ORDER BY id LIMIT 975,50;
--select * from student ORDER BY id;
--select * from school order by id;
--select COUNT(*) from student ORDER BY id LIMIT 10;

--select COUNT(*) 
--FROM (select id from student ORDER BY id LIMIT 975,50);

--select case when s2.school_id is NULL then false else true end, s1.school_name||'-'||s1.school_id
--from school s1
--LEFT OUTER JOIN school s2
--ON s1.school_id = s2.school_id
--AND s2.school_id = 10
--;

--DROP TABLE student;
--CREATE TABLE "student" (
--	"student_id"	INTEGER,
--	"first_name"	TEXT,
--	"last_name"	    TEXT,
--	"email"	        TEXT,
--    "school_id"     INTEGER,
--    "active"        BOOLEAN,
--    PRIMARY KEY("student_id" AUTOINCREMENT)
--);

--DROP TABLE school;

--CREATE TABLE "school" (
--	"school_id"	    INTEGER,
--	"school_name"	TEXT,
--  "active"        BOOLEAN,
--	PRIMARY KEY("school_id" AUTOINCREMENT)
--);

--INSERT INTO school (school_id, school_name, active) VALUES (0, 'None', false);

--CREATE TABLE "instructor" (
--	"instructor_id"	    INTEGER,
--	"instructor_name"	TEXT,
--  "active"        BOOLEAN,
--	PRIMARY KEY("instructor_id" AUTOINCREMENT)
--);

--CREATE TABLE "time_period" (
--	"time_period_id" INTEGER,
--	"description"	 TEXT,
--    "start_date"	 DATE,
--    "end_date"  	 DATE,
--	PRIMARY KEY("time_period_id" AUTOINCREMENT)
--)

--CREATE TABLE "instructor_availability" (
--	"instructor_id"	 INTEGER,
--    "time_period_id" INTEGER,
--    "school_id"      INTEGER,
--    "day_of_week"    TEXT,
--    "start_time"     INTEGER,
--    "end_time"       INTEGER,
--	PRIMARY KEY("instructor_id", "time_period_id", "school_id", "day_of_week", "start_time")
--)

--CREATE TABLE "school_availability" (
--    "school_id"      INTEGER,
--    "time_period_id" INTEGER,
--    "day_of_week"    TEXT,
--    "start_time"     INTEGER,
--    "end_time"       INTEGER,
--    PRIMARY KEY("school_id", "time_period_id", "day_of_week", "start_time")
--)