-- NOTE: This would eventually be only for an integration or test setup as it will setup / blow away data
-- We intend to seed the fake blog database with a set of interesting data so that reviewers can 
-- see the impact of their actions but also see what the intended setup would look like
CREATE DATABASE IF NOT EXISTS `mysql_blogforge` CHARACTER SET utf8mb4; -- COLLATE utf8mb4_unicode_ci;
GRANT ALL PRIVILEGES ON mysql_blogforge.* TO `forge-service`@`%`;
USE mysql_blogforge;

DROP TABLE IF EXISTS `bloggers`;

CREATE TABLE `bloggers`(
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `username` VARCHAR(255) NOT NULL,
    `first_name` VARCHAR(255) NOT NULL,
    `last_name`  VARCHAR(255) NOT NULL,
    `pass_hash` VARCHAR(255) NOT NULL,
    PRIMARY KEY (`id`),
    KEY `usr` (`username`),
    CONSTRAINT uniq_user UNIQUE (`username`)
);

-- assumes that the standard encryption key is in effect for pass hash at this time.
START TRANSACTION;
    INSERT INTO `bloggers`( `username`, `pass_hash`, `first_name`, `last_name`)
    VALUES('nate-admin', 'pass1', "n", "f");
    INSERT INTO `bloggers`( `username`, `pass_hash`, `first_name`, `last_name`)
    VALUES('nate-normal', 'pass2', "n", "food");
COMMIT;


DROP TABLE IF EXISTS `entries`;
CREATE TABLE `entries`(
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `owner_id` int(11) NOT NULL,
    `path` VARCHAR(255),
    `raw` LONGTEXT,
    PRIMARY KEY (`id`),
    CONSTRAINT `owner_of_entry` FOREIGN KEY (`owner_id`) REFERENCES `bloggers`(`id`) ON DELETE CASCADE
);


START TRANSACTION;
    INSERT INTO `entries`(`owner_id`, `raw`)
    VALUES(1, "I am a blog entry that is stored in the db rather than afull doc store.\n We should have done this in mongo or implemented at rest data.");

COMMIT;

