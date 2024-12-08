CREATE DATABASE IF NOT EXISTS omcc;
USE omcc;
-- 创建客户表
CREATE TABLE IF NOT EXISTS customers (
    id VARCHAR(36) PRIMARY KEY NOT NULL,
    username VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;
--
CREATE TABLE IF NOT EXISTS social_platforms (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR (50) NOT NULL UNIQUE,
    is_active BOOLEAN DEFAULT TRUE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;
--
INSERT INTO social_platforms (name)
VALUES ('TELEGRAM');
INSERT INTO social_platforms (name)
VALUES ('LINE');
--
CREATE TABLE IF NOT EXISTS trading_platforms (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    is_active BOOLEAN DEFAULT TRUE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;
--
INSERT INTO trading_platforms (name)
VALUES ('BITGET');
INSERT INTO trading_platforms (name)
VALUES ('BINGX');
--
CREATE TABLE IF NOT EXISTS customer_social_bindings (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    customer_id VARCHAR(36) NOT NULL,
    social_id INT NOT NULL,
    user_id VARCHAR(50) NOT NULL,
    username VARCHAR(50),
    firstname VARCHAR(50),
    lastname VARCHAR(50),
    is_active BOOLEAN DEFAULT TRUE,
    deactivated_at TIMESTAMP NULL,
    member_status ENUM('creator', 'administrator', 'member', 'restricted', 'left', 'kicked'),
    status ENUM('normal', 'whitelisted', 'blacklisted') NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_customer (customer_id),
    INDEX idx_social_uid (social_id),
    INDEX idx_user_uid (user_id),
    FOREIGN KEY (customer_id) REFERENCES customers(id),
    FOREIGN KEY (social_id) REFERENCES social_platforms(id)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;
--
CREATE TABLE customer_trading_bindings (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    customer_id VARCHAR(36) NOT NULL,
    trading_id INT NOT NULL,
    uid VARCHAR(50) NOT NULL,
    register_time TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_trading_uid (trading_id, uid),
    INDEX idx_uid (uid),
    FOREIGN KEY (customer_id) REFERENCES customers (id),
    FOREIGN KEY (trading_id) REFERENCES trading_platforms (id)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;
--
CREATE TABLE IF NOT EXISTS trading_histories (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    binding_id BIGINT NOT NULL,
    volume DECIMAL(16, 2) NOT NULL,
    time_period ENUM('daily', 'weekly', 'monthly'),
    trading_date TIMESTAMP NOT NULL,
    FOREIGN KEY (binding_id) REFERENCES customer_trading_bindings (id)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;