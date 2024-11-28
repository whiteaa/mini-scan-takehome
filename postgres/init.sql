CREATE TABLE scan (
    id SERIAL PRIMARY KEY,
    ip VARCHAR,
    port int,
    service VARCHAR,
    scanned_at int,
    data_version int,
    data VARCHAR
);
CREATE UNIQUE INDEX idx_ip_port_service ON scan(ip, port, service);
