CREATE TABLE users (
	id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
	user_id VARCHAR(50) NOT NULL,
	name VARCHAR(50) NOT NULL,
	connecting boolean DEFAULT false,
	location geography(POINT, 4326),
	bio VARCHAR(250) DEFAULT ''
);
