ALTER SYSTEM SET password_encryption = 'scram-sha-256';
SELECT pg_reload_conf();
