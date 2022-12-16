CREATE TABLE public.player (
	id serial NOT NULL,
	"name" varchar NULL,
	team_id int8 NULL,
	CONSTRAINT player_pk PRIMARY KEY (id)
);
