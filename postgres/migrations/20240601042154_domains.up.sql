CREATE TABLE IF NOT EXISTS domains (
    id          	UUID         	PRIMARY KEY,
	name			TEXT			NOT NULL,
	extension		TEXT			NOT NULL,
    created_at  	TIMESTAMPTZ    	NOT NULL,
    updated_at  	TIMESTAMPTZ    	NOT NULL
);

CREATE UNIQUE INDEX ON domains (name, extension);
