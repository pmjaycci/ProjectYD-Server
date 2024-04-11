CREATE TABLE IF NOT EXISTS account (
    uid VARCHAR(45) NOT NULL PRIMARY KEY,
    user_id VARCHAR(45) NOT NULL,
    user_name VARCHAR(20) DEFAULT "",
    money int DEFAULT 0,
    INDEX (user_name)
);

CREATE TABLE IF NOT EXISTS time_attack_rank (
    uid VARCHAR(45) NOT NULL,
    user_name VARCHAR(20),
    record_time float NOT NULL DEFAULT 0,
    PRIMARY KEY (uid),
   	INDEX (user_name),
    FOREIGN KEY (user_name) REFERENCES account(user_name) ON UPDATE CASCADE,
    INDEX(record_time)
);


CREATE TABLE IF NOT EXISTS items (
    id int NOT NULL AUTO_INCREMENT,
    item_name VARCHAR(40),
    item_type int NOT NULL,
    category int NOT NULL DEFAULT 0,
    img_name VARCHAR(45),
    is_stack BOOLEAN NOT NULL DEFAULT false,
    PRIMARY KEY(id),
    INDEX(item_name)
);

CREATE TABLE IF NOT EXISTS item_weapon (
    id int NOT NULL,
    item_name VARCHAR(40),
    dmg int DEFAULT 0,
    atk_spd float DEFAULT 0.0,
    atk_range int DEFAULT 0,
    PRIMARY KEY (id),
    FOREIGN KEY (id) REFERENCES items(id) ON DELETE CASCADE,
    FOREIGN KEY (item_name) REFERENCES items(item_name) ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS item_effect (
    id int NOT NULL,
    item_name VARCHAR(40),
    max_hp int DEFAULT 0, 
    regen_hp int DEFAULT 0,
    short_dmg int DEFAULT 0,
    PRIMARY KEY (id),
    FOREIGN KEY (id) REFERENCES items(id) ON DELETE CASCADE,
    FOREIGN KEY (item_name) REFERENCES items(item_name) ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS shop (
    id int NOT NULL,
    money_type int,
    price int,
    PRIMARY KEY (id),
    FOREIGN KEY (id) REFERENCES items(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS shop_ingame (
    id int NOT NULL,
    price int DEFAULT 0,
    PRIMARY KEY (id),
    FOREIGN KEY (id) REFERENCES items(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS weapon_enchant_probability (
    enchant_level int,
    probability int,
    PRIMARY KEY (enchant_level)
);

