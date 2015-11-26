CREATE DATABASE `std` DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;
GRANT ALL PRIVILEGES ON std.* TO 'std'@'%' IDENTIFIED BY 'std';
GRANT ALL PRIVILEGES ON std.* TO 'std'@'localhost' IDENTIFIED BY 'std';
FLUSH PRIVILEGES;

USE std;


-- -----------------------------------------------------
-- Table `std`.`user`
-- -----------------------------------------------------
CREATE  TABLE IF NOT EXISTS `std`.`user` (
  `user_id` INT NOT NULL AUTO_INCREMENT ,
  `username` CHAR(40) NOT NULL ,
  `password` CHAR(20) NOT NULL ,
  `created` INT NOT NULL COMMENT '创建时间' ,
  PRIMARY KEY (`user_id`) ,
  UNIQUE INDEX `username_UNIQUE` (`username` ASC) )
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `std`.`topic`
-- -----------------------------------------------------
CREATE  TABLE IF NOT EXISTS `std`.`topic` (
  `topic_id` INT NOT NULL AUTO_INCREMENT  PRIMARY KEY,
  `user_id` INT NOT NULL ,
  `content` TEXT NOT NULL ,
  `created` INT NOT NULL COMMENT '创建时间' ,
  INDEX `fk_topic_user1_idx` (`user_id` ASC))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `std`.`label`
-- -----------------------------------------------------
CREATE  TABLE IF NOT EXISTS `std`.`label` (
  `label_id` INT NOT NULL AUTO_INCREMENT ,
  `name` CHAR(50) NOT NULL ,
  `created` INT NOT NULL ,
  PRIMARY KEY (`label_id`) )
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `std`.`keyword`
-- -----------------------------------------------------
CREATE  TABLE IF NOT EXISTS `std`.`keyword` (
  `keyword_id` INT NOT NULL AUTO_INCREMENT ,
  `keyword` CHAR(50) NOT NULL ,
  `created` INT NOT NULL ,
  PRIMARY KEY (`keyword_id`) )
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `std`.`comment`
-- -----------------------------------------------------
CREATE  TABLE IF NOT EXISTS `std`.`comment` (
  `comment_id` INT NOT NULL AUTO_INCREMENT ,
  `topic_id` INT NOT NULL ,
  `creator` INT NOT NULL COMMENT '创建者' ,
  `reply_to` INT NOT NULL DEFAULT 0 COMMENT '回复某个用户，在content中头部将会自动记录： @xxx' ,
  `content` CHAR(250) NOT NULL ,
  `created` INT NOT NULL COMMENT '创建时间' ,
  PRIMARY KEY (`comment_id`) ,
  INDEX `fk_comment_topic1_idx` (`topic_id` ASC) ,
  INDEX `fk_comment_user1_idx` (`creator` ASC) ,
  INDEX `fk_comment_user2_idx` (`reply_to` ASC))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `std`.`like`
-- -----------------------------------------------------
CREATE  TABLE IF NOT EXISTS `std`.`like` (
  `like_id` INT NOT NULL AUTO_INCREMENT ,
  `user_id` INT NOT NULL ,
  `type` TINYINT(1) NOT NULL COMMENT '喜欢的类型， 1: topic, 2:comment' ,
  `type_id` INT NOT NULL COMMENT '具体喜欢类型对应的数据主Id.' ,
  `created` INT NOT NULL COMMENT '喜欢的时间' ,
  PRIMARY KEY (`like_id`) ,
  INDEX `fk_like_user1_idx` (`user_id` ASC))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `std`.`profile`
-- -----------------------------------------------------
CREATE  TABLE IF NOT EXISTS `std`.`profile` (
  `user_id` INT NOT NULL ,
  `gender` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '1: man, 2: woman' ,
  `age` TINYINT(2) NOT NULL DEFAULT 0 ,
  `iiiness` TINYINT(2) NOT NULL COMMENT '病症Id' ,
  `avator` VARCHAR(100) NULL COMMENT '头像url地址' ,
  PRIMARY KEY (`user_id`))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `std`.`notification`
-- -----------------------------------------------------
CREATE  TABLE IF NOT EXISTS `std`.`notification` (
  `notification_id` INT NOT NULL AUTO_INCREMENT COMMENT '通知谁' ,
  `notice_to` INT NOT NULL ,
  `message` CHAR(250) NOT NULL COMMENT '通知内容' ,
  `created` INT NOT NULL ,
  PRIMARY KEY (`notification_id`) ,
  INDEX `fk_notification_user1_idx` (`notice_to` ASC))
ENGINE = InnoDB;



-- -----------------------------------------------------
-- Table `std`.`user_loaded_max_topic`
-- -----------------------------------------------------
CREATE  TABLE IF NOT EXISTS `std`.`user_loaded_max_topic` (
  `user_id` INT NOT NULL ,
  `max_topic_id` INT NOT NULL ,
  PRIMARY KEY (`user_id`, `max_topic_id`) ,
  INDEX `fk_user_loaded_max_topic_topic1_idx` (`max_topic_id` ASC))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `std`.`user_like_keyword_rate`
-- -----------------------------------------------------
CREATE  TABLE IF NOT EXISTS `std`.`user_like_keyword_rate` (
  `user_id` INT NOT NULL ,
  `keyword_id` INT NOT NULL ,
  `rate` SMALLINT(6) NOT NULL COMMENT '初始值为100,对大值为1000,表示1000%的喜欢，当然也可以不限制;最小值为-100, 表示100%的讨厌, 默认值为0' ,
  PRIMARY KEY (`user_id`, `keyword_id`) ,
  INDEX `fk_user_has_keyword_keyword1_idx` (`keyword_id` ASC) ,
  INDEX `fk_user_has_keyword_user1_idx` (`user_id` ASC))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `std`.`label_has_keyword`
-- -----------------------------------------------------
CREATE  TABLE IF NOT EXISTS `std`.`label_has_keyword` (
  `label_has_keyword_id` INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `label_label_id` INT NOT NULL ,
  `keyword_keyword_id` INT NOT NULL ,
  INDEX `fk_label_has_keyword_keyword1_idx` (`keyword_keyword_id` ASC) ,
  INDEX `fk_label_has_keyword_label1_idx` (`label_label_id` ASC))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `std`.`topic_has_keyword`
-- -----------------------------------------------------
CREATE  TABLE IF NOT EXISTS `std`.`topic_has_keyword` (
  `topic_id` INT NOT NULL ,
  `keyword_id` INT NOT NULL ,
  PRIMARY KEY (`topic_id`, `keyword_id`) ,
  INDEX `fk_topic_has_keyword_keyword1_idx` (`keyword_id` ASC) ,
  INDEX `fk_topic_has_keyword_topic1_idx` (`topic_id` ASC))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `std`.`notification_setting`
-- -----------------------------------------------------
CREATE  TABLE IF NOT EXISTS `std`.`notification_setting` (
  `notification_setting_id` INT NOT NULL AUTO_INCREMENT ,
  `user_id` INT NOT NULL ,
  `type` TINYINT(1) NOT NULL ,
  `is_send` TINYINT(1) NOT NULL ,
  PRIMARY KEY (`notification_setting_id`, `user_id`) ,
  INDEX `fk_notification_setting_user1_idx` (`user_id` ASC))
ENGINE = InnoDB;


CREATE TABLE IF NOT EXISTS `std`.`oauth_clients` (
  `oauth_clients_id` INT NOT NULL AUTO_INCREMENT,
  `app_key` CHAR(50) NOT NULL,
  `secret_key` CHAR(50) NOT NULL,
  `descr` CHAR(250) NOT NULL DEFAULT '',
  `created` INT NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`oauth_clients_id`))
ENGINE = InnoDB;


CREATE  TABLE IF NOT EXISTS `std`.`oauth_access_tokens` (
  `oauth_access_tokens_id` INT NOT NULL AUTO_INCREMENT ,
  `oauth_clients_id` INT NOT NULL ,
  `access_token` CHAR(32) NOT NULL ,
  `expires_in` SMALLINT(5) NOT NULL ,
  `created` INT NOT NULL ,
  PRIMARY KEY (`oauth_access_tokens_id`) ,
  INDEX `fk_oauth_access_tokens_oauth_clients1_idx` (`oauth_clients_id` ASC))
ENGINE = InnoDB;

insert into oauth_clients value('', 'e48189a13f8a2c8a2c185b03e9346ec29982277f', '0362919850b479fd8aeff1d9813dd31231ef999d', 'mobile admin', 1410426160);

CREATE  TABLE IF NOT EXISTS `std`.`user_session` (
  `session_id` INT NOT NULL AUTO_INCREMENT ,
  `session_key` CHAR(32) NOT NULL ,
  `user_id` INT NOT NULL ,
  `session_val` CHAR(255) NOT NULL DEFAULT '' ,
  `created` INT NOT NULL ,
  `expried` INT NOT NULL ,
  PRIMARY KEY (`session_id`) )
ENGINE = InnoDB
COMMENT = '用户Session表';

create index user_session_session_key on user_session(session_key);
create index user_session_expried on user_session(expried);
create index user_session_user_id on user_session(user_id);
