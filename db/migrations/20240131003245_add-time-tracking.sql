-- migrate:up
CREATE TABLE IF NOT EXISTS time_tracking_sessions (
    session_id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    status INT NOT NULL,
    checked_in_time TIMESTAMP WITH TIME ZONE NOT NULL,
    checked_out_time TIMESTAMP WITH TIME ZONE,
    total_work_time INTERVAL,
    total_break_time INTERVAL,
    last_known_location GEOMETRY(Point, 4326)
);
CREATE TABLE IF NOT EXISTS time_tracking_breaks (
    break_id UUID PRIMARY KEY,
    session_id UUID NOT NULL,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE,
    CONSTRAINT fk_session FOREIGN KEY(session_id) REFERENCES time_tracking_sessions(session_id)
);
-- migrate:down
DROP TABLE IF EXISTS time_tracking_breaks;
DROP TABLE IF EXISTS time_tracking_sessions;