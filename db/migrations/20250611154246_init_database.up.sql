CREATE TYPE user_role AS ENUM ('ADMIN', 'EMPLOYEE');

-- TODO: make index, unique, and foreign key
CREATE TABLE users (
	id SERIAL PRIMARY KEY,
	username VARCHAR(255) NOT NULL,
	password VARCHAR(255) NOT NULL,
	role user_role NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP DEFAULT NULL
);

CREATE TABLE user_infos (
	id SERIAL PRIMARY KEY,
	user_id INT NOT NULL,
	monthly_salary BIGINT DEFAULT NULL
);

-- CREATE TYPE attendance as ENUM ('CHECKIN', 'CHECKOUT');

-- CREATE TABLE user_attendances {
-- 	id SERIAL PRIMARY KEY,
-- 	user_id INT NOT NULL,
-- 	type attendance NOT NULL,
-- 	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
-- 	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
-- 	deleted_at TIMESTAMP DEFAULT NULL
-- }

-- CREATE TABLE user_overtimes {
-- 	id SERIAL PRIMARY KEY,
-- 	user_id INT NOT NULL,
-- 	description TEXT NOT NULL,
-- 	overtime_at TIMESTAMP NOT NULL,
-- 	duration_minutes INT NOT NULL,
-- 	approved_by_user_id INT DEFAULT NULL,
-- 	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
-- 	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
-- 	deleted_at TIMESTAMP DEFAULT NULL
-- };

-- CREATE TABLE user_reimbursements {
-- 	id SERIAL PRIMARY KEY,
-- 	user_id INT NOT NULL,
-- 	description TEXT NOT NULL,
-- 	amount INT NOT NULL,
-- 	approved_by_user_id INT DEFAULT NULL,
-- 	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
-- 	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
-- 	deleted_at TIMESTAMP DEFAULT NULL
-- };

-- CREATE TABLE payrolls (
-- 	id SERIAL PRIMARY KEY,
-- 	name VARCHAR(255) NOT NULL,
-- 	started_at TIMESTAMP NOT NULL,
-- 	ended_at TIMESTAMP NOT NULL,
-- 	rolled_at TIMESTAMP DEFAULT NULL,
-- 	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
-- 	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
-- 	deleted_at TIMESTAMP DEFAULT NULL
-- );

-- CREATE TABLE user_payslips_summary (
-- 	id SERIAL PRIMARY KEY,
-- 	payroll_id INT NOT NULL,
-- 	user_id INT NOT NULL,

-- 	total_take_home_pay INT NOT NULL,

-- 	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
-- 	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
-- 	deleted_at TIMESTAMP DEFAULT NULL
-- );