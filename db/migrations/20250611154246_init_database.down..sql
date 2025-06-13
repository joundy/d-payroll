BEGIN;

DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS user_infos;
DROP TABLE IF EXISTS user_attendances;
DROP TABLE IF EXISTS user_overtimes;
DROP TABLE IF EXISTS user_reimbursements;
DROP TABLE IF EXISTS user_payslips_summary;
DROP TABLE IF EXISTS payrolls;

DROP TYPE IF EXISTS attendance;
DROP TYPE IF EXISTS user_role;

COMMIT;