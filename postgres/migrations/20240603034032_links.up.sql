CREATE TABLE IF NOT EXISTS links (
    from_id         UUID         	NOT NULL REFERENCES pages(id) ON DELETE CASCADE,
	to_id			UUID			NOT NULL REFERENCES pages(id) ON DELETE CASCADE,
    created_at  	TIMESTAMPTZ    	NOT NULL,

	PRIMARY KEY (from_id, to_id)
);
