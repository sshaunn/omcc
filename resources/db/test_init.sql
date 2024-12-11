CREATE DATABASE test_db;
CREATE USER 'test'@'localhost' IDENTIFIED BY 'test';
GRANT ALL PRIVILEGES ON test_db.* TO 'test'@'localhost';
FLUSH PRIVILEGES;