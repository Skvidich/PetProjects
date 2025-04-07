CREATE TABLE IF NOT EXISTS raw_reports (
    id SERIAL PRIMARY KEY,
    service_name VARCHAR(100) NOT NULL,
    get_time TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS raw_component_metrics (
    id SERIAL PRIMARY KEY,
    report_id INTEGER NOT NULL,
    component_name VARCHAR(100) NOT NULL,
    status VARCHAR(100) NOT NULL,
    FOREIGN KEY (report_id) REFERENCES raw_reports(id)
);