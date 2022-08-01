DROP TABLE IF EXISTS sellers;
CREATE TABLE sellers (
    id              INT NOT NULL AUTO_INCREMENT,
    family_name     VARCHAR(30) NOT NULL,
    first_name      VARCHAR(30) NOT NULL,
    nickname        VARCHAR(10) UNIQUE NOT NULL,
    company         VARCHAR(137) NOT NULL,
    mail            VARCHAR(254) UNIQUE NOT NULL,
    password        VARCHAR(128) NOT NULL,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

DROP TABLE IF EXISTS consumers;
CREATE TABLE consumers (
    id               INT NOT NULL AUTO_INCREMENT,
    family_name      VARCHAR(30) NOT NULL,
    first_name       VARCHAR(30) NOT NULL,
    nickname         VARCHAR(10) UNIQUE NOT NULL,
    company          VARCHAR(137) NOT NULL,
    lottery_units    TINYINT NOT NULL,
    mail             VARCHAR(254) UNIQUE NOT NULL,
    password         VARCHAR(128) NOT NULL,
    created_at       DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

DROP TABLE IF EXISTS images;
CREATE TABLE images (
    id            INT NOT NULL AUTO_INCREMENT,
    file_name     VARCHAR(30) UNIQUE NOT NULL,
    created_at    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

DROP TABLE IF EXISTS upload;
CREATE TABLE upload (
    id             INT NOT NULL AUTO_INCREMENT,
    image_id       INT NOT NULL,
    seller_id      INT NOT NULL,
    created_at     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),

    FOREIGN KEY (image_id)
        REFERENCES images(id),

    FOREIGN KEY (seller_id)
        REFERENCES sellers(id)
);

DROP TABLE IF EXISTS lottery;
CREATE TABLE lottery (
    id            INT NOT NULL AUTO_INCREMENT,
    image_id      INT NOT NULL,
    consumer_id   INT NOT NULL,
    created_at    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),

    FOREIGN KEY (image_id)
        REFERENCES images(id),

    FOREIGN KEY (consumer_id)
        REFERENCES consumers(id)
 );

DROP TABLE IF EXISTS portion;
CREATE TABLE portion (
    lottery_id INT NOT NULL,
    portion TINYINT NOT NULL,

    FOREIGN KEY (lottery_id)
        REFERENCES  lottery(id)
);
