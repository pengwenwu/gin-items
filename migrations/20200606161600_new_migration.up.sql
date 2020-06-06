# Dump of table item_parameters
# ------------------------------------------------------------

DROP TABLE IF EXISTS `item_parameters`;

CREATE TABLE `item_parameters` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `item_id` int(11) NOT NULL COMMENT '商品ID',
  `parameters` varchar(30) NOT NULL COMMENT '属性名',
  `value` varchar(150) NOT NULL COMMENT '属性值',
  `state` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态(0:已删除，1：正常)',
  `sort` tinyint(4) NOT NULL COMMENT '排序，正序',
  `last_dated` datetime NOT NULL,
  `dated` datetime NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `item_id` (`item_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品的属性表';



# Dump of table item_photos
# ------------------------------------------------------------

DROP TABLE IF EXISTS `item_photos`;

CREATE TABLE `item_photos` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `item_id` int(11) NOT NULL COMMENT '商品id',
  `photo` varchar(255) NOT NULL COMMENT '图片地址',
  `sort` int(11) NOT NULL COMMENT '排序',
  `state` tinyint(4) NOT NULL COMMENT '状态(0：已删除，1：正常)',
  `last_dated` datetime NOT NULL COMMENT '最后更新时间',
  `dated` datetime NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `item_id` (`item_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;



# Dump of table item_prop_values
# ------------------------------------------------------------

DROP TABLE IF EXISTS `item_prop_values`;

CREATE TABLE `item_prop_values` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `item_id` int(11) NOT NULL COMMENT '商品id',
  `prop_name` varchar(255) NOT NULL COMMENT '规格名称',
  `prop_value_name` varchar(255) NOT NULL COMMENT '属性名称',
  `sort` int(11) NOT NULL COMMENT '排序',
  `prop_photo` varchar(255) NOT NULL COMMENT '规格图片',
  `prop_desc` varchar(255) NOT NULL COMMENT '规格介绍',
  `assisted_num` int(11) NOT NULL COMMENT '辅助数量',
  `state` tinyint(4) NOT NULL COMMENT '状态(0：已删除，1：正常)',
  `last_dated` datetime NOT NULL COMMENT '最后更新时间',
  `dated` datetime NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `item_id` (`item_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品规格对应的属性';



# Dump of table item_props
# ------------------------------------------------------------

DROP TABLE IF EXISTS `item_props`;

CREATE TABLE `item_props` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键，自增',
  `item_id` int(11) NOT NULL COMMENT '商品id',
  `prop_name` varchar(255) NOT NULL COMMENT '规格名',
  `sort` int(11) NOT NULL COMMENT '排序',
  `have_photo` tinyint(4) NOT NULL COMMENT '是否有规格图片（0：无，1：有）',
  `prop_desc` varchar(255) NOT NULL COMMENT '规格描述',
  `state` tinyint(4) NOT NULL COMMENT '状态(0：已删除，1：正常)',
  `last_dated` datetime NOT NULL COMMENT '更新时间',
  `dated` datetime NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `item_id` (`item_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品规格表';



# Dump of table item_searches
# ------------------------------------------------------------

DROP TABLE IF EXISTS `item_searches`;

CREATE TABLE `item_searches` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `appkey` char(32) NOT NULL COMMENT '所属应用',
  `channel` int(11) NOT NULL COMMENT '所属channel',
  `item_id` int(11) NOT NULL COMMENT '商品id',
  `sku_id` int(11) NOT NULL COMMENT 'sku id',
  `sku_name` varchar(255) NOT NULL COMMENT 'sku名称',
  `bar_code` varchar(20) NOT NULL COMMENT '条码',
  `sku_code` varchar(20) NOT NULL COMMENT 'sku速记码',
  `item_state` tinyint(4) NOT NULL COMMENT '商品状态 （0：删除(放入回收站) 1：正常  2：彻底删除）',
  `sku_state` tinyint(4) NOT NULL COMMENT 'sku状态 （0：删除(放入回收站) 1：正常  2：彻底删除）',
  `last_dated` datetime NOT NULL COMMENT '最后更新时间',
  `dated` datetime NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `item_id` (`item_id`),
  KEY `sku_id` (`sku_id`),
  KEY `bar_code` (`bar_code`),
  KEY `sku_code` (`sku_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;



# Dump of table item_skus
# ------------------------------------------------------------

DROP TABLE IF EXISTS `item_skus`;

CREATE TABLE `item_skus` (
  `sku_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '商品sku_id',
  `item_id` int(11) NOT NULL COMMENT '商品ID',
  `appkey` char(32) NOT NULL COMMENT '所属应用',
  `channel` int(11) NOT NULL COMMENT '所属channel',
  `item_name` varchar(255) NOT NULL COMMENT '商品名称',
  `sku_name` varchar(255) NOT NULL COMMENT 'sku名称',
  `sku_photo` varchar(512) NOT NULL COMMENT 'sku图片',
  `sku_code` varchar(50) NOT NULL COMMENT 'sku速记码',
  `bar_code` varchar(50) NOT NULL COMMENT '条形码',
  `properties` varchar(512) NOT NULL COMMENT '规格属性字符串',
  `state` tinyint(4) NOT NULL COMMENT '状态（0：删除(放入回收站) 1：正常  2：彻底删除，3: sku本身的删除）',
  `last_dated` datetime NOT NULL COMMENT '最后更新时间',
  `dated` datetime NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`sku_id`),
  KEY `item_id` (`item_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品sku';



# Dump of table items
# ------------------------------------------------------------

DROP TABLE IF EXISTS `items`;

CREATE TABLE `items` (
  `item_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '商品ID',
  `appkey` char(32) NOT NULL COMMENT '应用appkey',
  `channel` int(11) NOT NULL COMMENT '所属channel',
  `name` varchar(255) NOT NULL COMMENT '商品名称',
  `photo` varchar(512) NOT NULL COMMENT '商品主图',
  `detail` text NOT NULL COMMENT '默认详情',
  `state` tinyint(4) NOT NULL COMMENT '状态（0：删除(放入回收站) 1：正常  2：彻底删除）',
  `last_dated` datetime NOT NULL COMMENT '最后更新时间',
  `dated` datetime NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`item_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品表';