CREATE TABLE IF NOT EXISTS reports (
    id SERIAL PRIMARY KEY,
    service_name VARCHAR(100) NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS component_metrics (
    id SERIAL PRIMARY KEY,
    report_id INTEGER NOT NULL,
    component_name VARCHAR(100) NOT NULL,
    status VARCHAR(100) NOT NULL,
    status_count SMALLINT NOT NULL,
    FOREIGN KEY (report_id) REFERENCES reports(id)
);

CREATE TABLE IF NOT EXISTS incidents (
    id SERIAL PRIMARY KEY,
    service_name VARCHAR(100) NOT NULL,
    component_name VARCHAR(100) NOT NULL,
    status VARCHAR(100) NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL
);
