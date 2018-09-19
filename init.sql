CREATE TABLE `tax-calc`.`order_detail` (
  `order_dtl_id` INT NOT NULL AUTO_INCREMENT,
  `item_name` VARCHAR(100) NULL,
  `item_tax_code` INT NULL,
  `item_amount` INT NULL,
  PRIMARY KEY (`order_dtl_id`));
