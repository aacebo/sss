CREATE TABLE IF NOT EXISTS pages (
    id          	UUID         	PRIMARY KEY,
	domain_id		UUID			NOT NULL REFERENCES domains(id) ON DELETE CASCADE,
	title			TEXT,
	url				TEXT			NOT NULL,
	address			TEXT			NOT NULL,
	size			BIGINT			NOT NULL,
	elapse_ms		BIGINT			NOT NULL,
	link_count		INT				NOT NULL,
    created_at  	TIMESTAMPTZ    	NOT NULL,
    updated_at  	TIMESTAMPTZ    	NOT NULL
);

CREATE UNIQUE INDEX ON pages (url);
