CREATE SCHEMA `wallet` ;
CREATE TABLE `wallet`.`users` (
  `alias` VARCHAR(45) NOT NULL,
  `first_name` VARCHAR(45) NOT NULL,
  `last_name` VARCHAR(45) NOT NULL,
  `email` VARCHAR(45) NOT NULL,
  `password` VARCHAR(45) NOT NULL,
  PRIMARY KEY (`alias`),
  UNIQUE INDEX `email_UNIQUE` (`email` ASC));

CREATE TABLE `wallet`.`movements_btc` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `mov_type` ENUM("deposit", "send","receive","init") NOT NULL, 
  `currency_name` VARCHAR(20) NOT NULL DEFAULT 'BTC',
  `date_created` DATETIME NOT NULL DEFAULT current_timestamp,
  `tx_amount` DECIMAL(18,8) ZEROFILL NOT NULL,
  `total_amount` DECIMAL(18,8) ZEROFILL NOT NULL,
  `interaction_alias` VARCHAR NOT NULL,
  `alias` VARCHAR NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `alias_idx` (`alias` ASC),
  CONSTRAINT `fk_btc_alias`
      FOREIGN KEY (`alias`)
          REFERENCES `wallet`.`users` (`alias`)
          ON DELETE CASCADE
          ON UPDATE CASCADE),
      FOREIGN KEY (`interaction_alias`)
          REFERENCES `wallet`.`users` (`alias`)
          ON DELETE CASCADE
          ON UPDATE CASCADE);

CREATE TABLE `wallet`.`movements_usdt` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `mov_type` ENUM("deposit", "send","receive","init") NOT NULL,
  `currency_name` VARCHAR(20) NOT NULL DEFAULT 'USDT',
  `date_created` DATETIME NOT NULL DEFAULT current_timestamp,
  `tx_amount` DECIMAL(18,2) ZEROFILL NOT NULL,
  `total_amount` DECIMAL(18,2) ZEROFILL NOT NULL,
  `interaction_alias` VARCHAR NOT NULL,
  `alias` VARCHAR NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `user_id_idx` (`alias` ASC),
  CONSTRAINT `fk_usdt_alias`
      FOREIGN KEY (`alias`)
          REFERENCES `wallet`.`users` (`alias`)
          ON DELETE CASCADE
          ON UPDATE CASCADE),  
      FOREIGN KEY (`interaction_alias`)
          REFERENCES `wallet`.`users` (`alias`)
          ON DELETE CASCADE
          ON UPDATE CASCADE);

CREATE TABLE `wallet`.`movements_ars` (
   `id` BIGINT NOT NULL AUTO_INCREMENT,
   `mov_type` ENUM("deposit", "send","receive","init") NOT NULL,
   `currency_name` VARCHAR(20) NOT NULL DEFAULT 'ARS',
   `date_created` DATETIME NOT NULL DEFAULT current_timestamp,
   `tx_amount` DECIMAL(18,2) ZEROFILL NOT NULL,
   `total_amount` DECIMAL(18,2) ZEROFILL NOT NULL,
   `interaction_alias` VARCHAR NOT NULL,
   `alias` VARCHAR NOT NULL,
   PRIMARY KEY (`id`),
   INDEX `user_id_idx` (`alias` ASC),
   CONSTRAINT `fk_ars_alias`
       FOREIGN KEY (`alias`)
           REFERENCES `wallet`.`users` (`alias`)
           ON DELETE CASCADE
           ON UPDATE CASCADE),
       FOREIGN KEY (`interaction_alias`)
          REFERENCES `wallet`.`users` (`alias`)
          ON DELETE CASCADE
          ON UPDATE CASCADE);


/*Triggers


DROP TRIGGER IF EXISTS `wallet`.`movements_usdt_BEFORE_INSERT`;
DELIMITER $$
USE `wallet`$$
CREATE DEFINER = CURRENT_USER TRIGGER `wallet`.`movements_usdt_BEFORE_INSERT` BEFORE INSERT ON `movements_usdt` FOR EACH ROW
BEGIN
IF NEW.mov_type='deposit' THEN
		SET NEW.total_amount = NEW.tx_amount +
		(SELECT total_amount FROM movements_usdt WHERE id =(SELECT max(id) from movements_usdt WHERE user_id = NEW.user_id));
END IF;
IF NEW.mov_type='extract' THEN
		SET NEW.total_amount = (SELECT total_amount FROM movements_usdt WHERE id =(SELECT max(id) from movements_usdt WHERE user_id = NEW.user_id))
        - NEW.tx_amount;
END IF;
END$$
DELIMITER ;



DROP TRIGGER IF EXISTS `wallet`.`movements_btc_BEFORE_INSERT`;

DELIMITER $$
USE `wallet`$$
CREATE DEFINER = CURRENT_USER TRIGGER `wallet`.`movements_btc_BEFORE_INSERT` BEFORE INSERT ON `movements_btc` FOR EACH ROW
BEGIN
IF NEW.mov_type='deposit' THEN
		SET NEW.total_amount = NEW.tx_amount +
		(SELECT total_amount FROM movements_btc WHERE id =(SELECT max(id) from movements_btc WHERE user_id = NEW.user_id));
END IF;
IF NEW.mov_type='extract' THEN
		SET NEW.total_amount = (SELECT total_amount FROM movements_btc WHERE id =(SELECT max(id) from movements_btc WHERE user_id = NEW.user_id))
        - NEW.tx_amount;
END IF;
END$$
DELIMITER ;



DROP TRIGGER IF EXISTS `wallet`.`movements_ars_BEFORE_INSERT`;

DELIMITER $$
USE `wallet`$$
CREATE DEFINER = CURRENT_USER TRIGGER `wallet`.`movements_ars_BEFORE_INSERT` BEFORE INSERT ON `movements_ars` FOR EACH ROW
BEGIN
IF NEW.mov_type='deposit' THEN
		SET NEW.total_amount = NEW.tx_amount +
		(SELECT total_amount FROM movements_ars WHERE id =(SELECT max(id) from movements_ars WHERE user_id = NEW.user_id));
END IF;
IF NEW.mov_type='extract' THEN
		SET NEW.total_amount = (SELECT total_amount FROM movements_ars WHERE id =(SELECT max(id) from movements_ars WHERE user_id = NEW.user_id))
        - NEW.tx_amount;
END IF;
END$$
DELIMITER ;
*/