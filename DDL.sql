DROP DATABASE IF EXISTS `leetcode_scheduler`;
CREATE DATABASE `leetcode_scheduler`;
USE `leetcode_scheduler`;

SET NAMES utf8;
set character_set_client = utf8mb4;

CREATE TABLE `problem_level` (
	`id` INT AUTO_INCREMENT,
    `level` CHAR(10) UNIQUE NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=INNODB;

CREATE TABLE `leetcode_problem` (
	`id` INT UNIQUE NOT NULL,
    `level_id` INT,
    `problem_name` CHAR(200) UNIQUE NOT NULL,
    `url` TEXT NOT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`level_id`)
        REFERENCES `problem_level` (`id`)
) ENGINE=INNODB;



INSERT INTO `problem_level`(`level`) VALUES("EASY");
INSERT INTO `problem_level`(`level`) VALUES("MEDIUM");
INSERT INTO `problem_level`(`level`) VALUES("HARD");


