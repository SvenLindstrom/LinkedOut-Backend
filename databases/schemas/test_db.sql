CREATE TABLE users (
	id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
	user_id VARCHAR(50) NOT NULL,
	name VARCHAR(50) NOT NULL,
	connecting boolean DEFAULT false,
	location geography(POINT, 4326),
	bio VARCHAR(250),
	profession VARCHAR(50),
	deviceCode VARCHAR(250)
);

CREATE TABLE interests (
	id UUID PRIMARY KEY NOT NULL,
	name VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE users_interests (
	user_id UUID NOT NULL,
	interest_id UUID NOT NULL,
	PRIMARY KEY (user_id, interest_id),
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
	FOREIGN KEY (interest_id) REFERENCES interests(id) ON DELETE CASCADE
);

-- MOCK INTERESTS

INSERT INTO interests (id, name) VALUES
('11111111-1111-1111-1111-111111111111', 'Software Development'),
('22222222-2222-2222-2222-222222222222', 'Data Analysis'),
('33333333-3333-3333-3333-333333333333', 'Cybersecurity'),
('44444444-4444-4444-4444-444444444444', 'Cloud Computing'),
('55555555-5555-5555-5555-555555555555', 'UI/UX Design'),
('66666666-6666-6666-6666-666666666666', 'Digital Marketing'),
('77777777-7777-7777-7777-777777777777', 'Project Management'),
('88888888-8888-8888-8888-888888888888', 'Financial Analysis'),
('99999999-9999-9999-9999-999999999999', 'Psychology'),
('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'Public Relations');

-- MOCK USERS

INSERT INTO users (id, user_id, name, bio, profession, deviceCode)
VALUES
('aaaa1111-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'user_001', 'Alice Example', 'Building scalable web systems.', 'Backend Developer', 'devCodeA'),
('bbbb2222-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'user_002', 'Bob Test', 'Designing engaging UI/UX.', 'UX Designer', 'devCodeB'),
('cccc3333-cccc-cccc-cccc-cccccccccccc', 'user_003', 'Charlie Demo', 'Analyzing business trends.', 'Data Analyst', 'devCodeC'),
('dddd4444-dddd-dddd-dddd-dddddddddddd', 'user_004', 'Dana Sample', 'Managing software projects.', 'Project Manager', 'devCodeD'),
('eeee5555-eeee-eeee-eeee-eeeeeeeeeeee', 'user_005', 'Eve Mock', 'Marketing and branding expert.', 'Digital Marketer', 'devCodeE');

-- MOCK INTERESTS PER USER

-- Alice Example
INSERT INTO users_interests (user_id, interest_id) VALUES
('aaaa1111-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '11111111-1111-1111-1111-111111111111'), -- Software Dev
('aaaa1111-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '44444444-4444-4444-4444-444444444444'), -- Cloud Computing
('aaaa1111-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '55555555-5555-5555-5555-555555555555'), -- UI/UX
('aaaa1111-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '77777777-7777-7777-7777-777777777777'), -- Project Mgmt
('aaaa1111-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '99999999-9999-9999-9999-999999999999'); -- Psychology

-- Bob Test
INSERT INTO users_interests (user_id, interest_id) VALUES
('bbbb2222-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '55555555-5555-5555-5555-555555555555'),
('bbbb2222-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '66666666-6666-6666-6666-666666666666'),
('bbbb2222-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '77777777-7777-7777-7777-777777777777'),
('bbbb2222-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '88888888-8888-8888-8888-888888888888'),
('bbbb2222-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa');

-- Charlie Demo
INSERT INTO users_interests (user_id, interest_id) VALUES
('cccc3333-cccc-cccc-cccc-cccccccccccc', '22222222-2222-2222-2222-222222222222'),
('cccc3333-cccc-cccc-cccc-cccccccccccc', '44444444-4444-4444-4444-444444444444'),
('cccc3333-cccc-cccc-cccc-cccccccccccc', '88888888-8888-8888-8888-888888888888'),
('cccc3333-cccc-cccc-cccc-cccccccccccc', '33333333-3333-3333-3333-333333333333'),
('cccc3333-cccc-cccc-cccc-cccccccccccc', '11111111-1111-1111-1111-111111111111');

-- Dana Sample
INSERT INTO users_interests (user_id, interest_id) VALUES
('dddd4444-dddd-dddd-dddd-dddddddddddd', '77777777-7777-7777-7777-777777777777'),
('dddd4444-dddd-dddd-dddd-dddddddddddd', '33333333-3333-3333-3333-333333333333'),
('dddd4444-dddd-dddd-dddd-dddddddddddd', '88888888-8888-8888-8888-888888888888'),
('dddd4444-dddd-dddd-dddd-dddddddddddd', '66666666-6666-6666-6666-666666666666'),
('dddd4444-dddd-dddd-dddd-dddddddddddd', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa');

-- Eve Mock
INSERT INTO users_interests (user_id, interest_id) VALUES
('eeee5555-eeee-eeee-eeee-eeeeeeeeeeee', '66666666-6666-6666-6666-666666666666'),
('eeee5555-eeee-eeee-eeee-eeeeeeeeeeee', '22222222-2222-2222-2222-222222222222'),
('eeee5555-eeee-eeee-eeee-eeeeeeeeeeee', '99999999-9999-9999-9999-999999999999'),
('eeee5555-eeee-eeee-eeee-eeeeeeeeeeee', '11111111-1111-1111-1111-111111111111'),
('eeee5555-eeee-eeee-eeee-eeeeeeeeeeee', '77777777-7777-7777-7777-777777777777');