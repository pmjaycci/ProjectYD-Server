CREATE TABLE rank_time_attack (
    uid VARCHAR(45) NOT NULL,
    clear_time TIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(uid),
    INDEX(clear_time),
);


CREATE TABLE inventory ( 
    id int NOT NULL AUTO_INCREMENT,
    uid VARCHAR(45) NOT NULL,
    item_id int NOT NULL,
    item_count int NOT NULL DEFAULT 1,
    enchant_level int NOT NULL DEFAULT 0,
    PRIMARY KEY(id),
    INDEX(uid),
    INDEX(item_id),
);

-- Redis
CREATE TABLE weapon_slot (
    uid VARCHAR(45) NOT NULL,
    slot_1_id int DEFAULT -1,
    slot_2_id int DEFAULT -1,
    PRIMARY KEY(uid),
);
