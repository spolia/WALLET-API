CREATE TABLE `users` (
  `alias` VARCHAR(45) NOT NULL,
  `first_name` VARCHAR(45) NOT NULL,
  `last_name` VARCHAR(45) NOT NULL,
  `email` VARCHAR(45) NOT NULL,
  `password` VARCHAR(45) NOT NULL,
  PRIMARY KEY (`alias`),
  UNIQUE INDEX `email_UNIQUE` (`email` ASC));

CREATE TABLE `movements_btc` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `mov_type` ENUM("deposit", "send","receive","init") NOT NULL, 
  `currency_name` VARCHAR(20) NOT NULL DEFAULT 'BTC',
  `date_created` DATETIME NOT NULL DEFAULT current_timestamp,
  `tx_amount` DECIMAL(18,8) ZEROFILL NOT NULL,
  `total_amount` DECIMAL(18,8) ZEROFILL NOT NULL,
  `interaction_alias` VARCHAR(45) NOT NULL,
  `alias` VARCHAR(45) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `alias_idx` (`alias` ASC),
  INDEX `interaction_alias_idx` (`interaction_alias` ASC),
  CONSTRAINT `fk_btc_alias`
      FOREIGN KEY (`alias`)
          REFERENCES `users` (`alias`)
          ON DELETE CASCADE
          ON UPDATE CASCADE,
  CONSTRAINT `fk_btc_interaction_alias`
      FOREIGN KEY (`interaction_alias`)
          REFERENCES `users` (`alias`)
          ON DELETE CASCADE
          ON UPDATE CASCADE);

CREATE TABLE `movements_usdt` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `mov_type` ENUM("deposit", "send","receive","init") NOT NULL,
  `currency_name` VARCHAR(20) NOT NULL DEFAULT 'USDT',
  `date_created` DATETIME NOT NULL DEFAULT current_timestamp,
  `tx_amount` DECIMAL(18,2) ZEROFILL NOT NULL,
  `total_amount` DECIMAL(18,2) ZEROFILL NOT NULL,
  `interaction_alias` VARCHAR(45) NOT NULL,
  `alias` VARCHAR(45) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `alias_id_idx` (`alias` ASC),
  INDEX `interaction_alias_idx` (`interaction_alias` ASC),
  CONSTRAINT `fk_usdt_alias`
      FOREIGN KEY (`alias`)
          REFERENCES `users` (`alias`)
          ON DELETE CASCADE
          ON UPDATE CASCADE,  
  CONSTRAINT `fk_usdt_interaction_alias`
      FOREIGN KEY (`interaction_alias`)
          REFERENCES `users` (`alias`)
          ON DELETE CASCADE
          ON UPDATE CASCADE);

CREATE TABLE `movements_ars` (
   `id` BIGINT NOT NULL AUTO_INCREMENT,
   `mov_type` ENUM("deposit", "send","receive","init") NOT NULL,
   `currency_name` VARCHAR(20) NOT NULL DEFAULT 'ARS',
   `date_created` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
   `tx_amount` DECIMAL(18,2) ZEROFILL NOT NULL,
   `total_amount` DECIMAL(18,2) ZEROFILL NOT NULL,
   `interaction_alias` VARCHAR(45) NOT NULL,
   `alias` VARCHAR(45) NOT NULL,
   PRIMARY KEY (`id`),
   INDEX `alias_id_idx` (`alias` ASC),
   INDEX `interaction_alias_idx` (`interaction_alias` ASC),
   CONSTRAINT `fk_ars_alias`
       FOREIGN KEY (`alias`)
           REFERENCES `users` (`alias`)
           ON DELETE CASCADE
           ON UPDATE CASCADE,
   CONSTRAINT `fk_ars_interaction_alias`
       FOREIGN KEY (`interaction_alias`)
          REFERENCES `users` (`alias`)
          ON DELETE CASCADE
          ON UPDATE CASCADE);


/*Triggers*/
DROP TRIGGER IF EXISTS `movements_usdt_BEFORE_INSERT`;
DELIMITER $$
CREATE TRIGGER movements_usdt_BEFORE_INSERT BEFORE INSERT ON movements_usdt FOR EACH ROW
BEGIN
IF NEW.mov_type='deposit' || NEW.mov_type='receive' THEN
		SET NEW.total_amount = NEW.tx_amount +
		(SELECT total_amount FROM movements_usdt WHERE id = (SELECT MAX(id)FROM movements_usdt WHERE alias = NEW.alias));
END IF;
IF NEW.mov_type='send' THEN
		SET NEW.total_amount = (SELECT total_amount FROM movements_usdt WHERE id = (SELECT MAX(id)FROM movements_usdt WHERE alias = NEW.alias))
        - NEW.tx_amount;
END IF;
END $$
DELIMITER ;

DROP TRIGGER IF EXISTS `movements_btc_BEFORE_INSERT`;
DELIMITER $$
CREATE TRIGGER movements_btc_BEFORE_INSERT BEFORE INSERT ON movements_btc FOR EACH ROW
BEGIN
IF NEW.mov_type='deposit' || NEW.mov_type='receive' THEN
		SET NEW.total_amount = NEW.tx_amount +
		(SELECT total_amount FROM movements_btc WHERE id = (SELECT MAX(id)FROM movements_btc WHERE alias = NEW.alias));
END IF;
IF NEW.mov_type='send' THEN
		SET NEW.total_amount = (SELECT total_amount FROM movements_btc WHERE id = (SELECT MAX(id)FROM movements_btc WHERE alias = NEW.alias))
        - NEW.tx_amount;
END IF;
END$$
DELIMITER ;

DROP TRIGGER IF EXISTS `movements_ars_BEFORE_INSERT`;
DELIMITER $$
CREATE TRIGGER movements_ars_BEFORE_INSERT BEFORE INSERT ON movements_ars FOR EACH ROW
BEGIN
IF NEW.mov_type='deposit' || NEW.mov_type='receive' THEN
		SET NEW.total_amount = NEW.tx_amount +
		(SELECT total_amount FROM movements_ars WHERE id = (SELECT MAX(id)FROM movements_ars WHERE alias = NEW.alias));
END IF;
IF NEW.mov_type='send' THEN
		SET NEW.total_amount = (SELECT total_amount FROM movements_ars WHERE id = (SELECT MAX(id)FROM movements_ars WHERE alias = NEW.alias))
        - NEW.tx_amount;
END IF;
END$$
DELIMITER ;