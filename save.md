use scamreport ;

CREATE TABLE user (
    id INT AUTO_INCREMENT PRIMARY KEY,
    phone VARCHAR(20) unique,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    password varchar(255),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
CREATE TABLE user_name (
	id INT AUTO_INCREMENT PRIMARY KEY,
    uid INT ,
    email VARCHAR(255) NOT NULL unique,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (uid) REFERENCES user(id)
);

CREATE TABLE report (
	id int AUTO_INCREMENT PRIMARY KEY,
    report_id int ,
    reporter_id int ,
    message text,
    reason text,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (report_id) REFERENCES user(id),
    FOREIGN KEY (reporter_id) REFERENCES user(id)
);

CREATE INDEX idx_id_phone ON user(id, phone);

egwt3EYxR_Uuw1msXy3n