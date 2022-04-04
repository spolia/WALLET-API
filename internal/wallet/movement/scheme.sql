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

