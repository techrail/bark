CREATE TABLE app_log
(
    id           BIGSERIAL PRIMARY KEY,
    log_time     TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT (NOW() AT TIME ZONE ('utc')),
    log_level    VARCHAR                     NOT NULL DEFAULT 'info',
    service_name VARCHAR(64)                 NOT NULL DEFAULT 'def_svc',
    session_name VARCHAR(64)                 NOT NULL DEFAULT 'def_sess',
    code         VARCHAR(16)                 NOT NULL DEFAULT '000000',
    msg          TEXT                        NOT NULL DEFAULT '_no_msg_supplied_',
    more_data    jsonb                       NOT NULL DEFAULT '{}'::jsonb
);

COMMENT ON TABLE app_log IS 'Table to store application logs';

COMMENT ON COLUMN app_log.id IS 'Unique log ID, Primary Key, Auto-incrementing number';
COMMENT ON COLUMN app_log.log_time IS 'UTC time when the log occurred. Defaults to "now()"';
COMMENT ON COLUMN app_log.log_level IS 'Severity level - INFO, WARNING, ERROR, PANIC etc.';
COMMENT ON COLUMN app_log.service_name IS 'Name of the service that sent this log (defaults to "def_svc")';
COMMENT ON COLUMN app_log.session_name IS 'Name of the session of the service (or Pod Name) that sent this log (defaults to "def_sess")';
COMMENT ON COLUMN app_log.code IS 'Unique Code of the log (LMID) - to be sent by the caller (defaults to "000000")';
COMMENT ON COLUMN app_log.msg IS 'Actual log message as text';
COMMENT ON COLUMN app_log.more_data IS 'Anything else that needs to be saved alongside the log entry';

