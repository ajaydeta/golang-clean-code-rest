package cmd

import (
	"gorm.io/gorm"
)

func MigrateAndSeed(db *gorm.DB) error {
	var err error

	if !db.Migrator().HasTable("product") {
		err = db.Exec(`
			create table product
			(
				id         varchar(36)  not null primary key,
				name       varchar(225) not null,
				price      double       not null,
				created_at datetime default current_timestamp
			)
		`).Error
		if err != nil {
			return err
		}

		err = seedProduct(db)
		if err != nil {
			return err
		}
	} else {
		if isEmpty(db, "product") {
			err = seedProduct(db)
			if err != nil {
				return err
			}
		}
	}

	if !db.Migrator().HasTable("category") {
		err = db.Exec(`
			create table category
			(
				id         varchar(36)  not null primary key,
				name       varchar(225) not null,
				created_at datetime default current_timestamp
			)
		`).Error
		if err != nil {
			return err
		}

		err = seedCategory(db)
		if err != nil {
			return err
		}
	} else {
		if isEmpty(db, "category") {
			err = seedCategory(db)
			if err != nil {
				return err
			}
		}
	}

	if !db.Migrator().HasTable("product_category") {
		err = db.Exec(`
			create table product_category
			(
				product_id  varchar(36) not null,
				category_id varchar(36) not null,
			
				primary key (product_id, category_id),
				CONSTRAINT fk_product_id FOREIGN KEY (product_id) REFERENCES product (id),
				CONSTRAINT fk_category_id FOREIGN KEY (category_id) REFERENCES category (id)
			)
		`).Error
		if err != nil {
			return err
		}

		err = seedProductCategory(db)
		if err != nil {
			return err
		}
	} else {
		if isEmpty(db, "product_category") {
			err = seedProductCategory(db)
			if err != nil {
				return err
			}
		}
	}

	if !db.Migrator().HasTable("customer") {
		err = db.Exec(`
			create table customer
			(
				id         varchar(36)         not null primary key,
				name       varchar(225)        not null,
				email      varchar(225) unique not null,
				password   varchar(225)        not null,
				created_at datetime default current_timestamp
			)
		`).Error
		if err != nil {
			return err
		}
	}

	if !db.Migrator().HasTable("transaction") {
		err = db.Exec(`
			create table transaction
			(
				id          varchar(36) not null primary key,
				customer_id varchar(36) not null,
				subtotal    double      not null,
				discount    double   default 0,
				total       double      not null,
				created_at  datetime default current_timestamp,
			
				CONSTRAINT fk_customer_id FOREIGN KEY (customer_id) REFERENCES customer (id)
			)
		`).Error
		if err != nil {
			return err
		}
	}

	if !db.Migrator().HasTable("transaction_item") {
		err = db.Exec(`
			create table transaction_item
			(
				id             varchar(36) not null primary key,
				transaction_id varchar(36) not null,
				product_id     varchar(36) not null,
				notes          text,
				price          double,
				qty            bigint,
				total          double,
				created_at     datetime default current_timestamp,
			
				CONSTRAINT fk_transaction_item_product_id FOREIGN KEY (product_id) REFERENCES product (id),
				CONSTRAINT fk_transaction_id FOREIGN KEY (transaction_id) REFERENCES transaction (id)
			)
		`).Error
		if err != nil {
			return err
		}
	}

	if !db.Migrator().HasTable("transaction_payment") {
		err = db.Exec(`
			create table transaction_payment
			(
				id             varchar(36) not null primary key,
				transaction_id varchar(36) not null,
				payment_type   varchar(36) not null comment '1_transfer_bank, 2_supermarket',
				paid           tinyint,
				created_at     datetime default current_timestamp,
			
				CONSTRAINT fk_transaction_payment_transaction_id FOREIGN KEY (transaction_id) REFERENCES transaction (id)
			)
		`).Error
		if err != nil {
			return err
		}
	}

	if !db.Migrator().HasTable("shopping_cart") {
		err = db.Exec(`
			create table shopping_cart
			(
				id          varchar(36) not null primary key,
				customer_id varchar(36) not null,
				product_id  varchar(36) not null,
				notes       text,
				qty         bigint,
				created_at  datetime default current_timestamp,
			
				CONSTRAINT fk_shopping_cart_customer_id FOREIGN KEY (customer_id) REFERENCES customer (id),
				CONSTRAINT fk_shopping_cart_product_id FOREIGN KEY (product_id) REFERENCES product (id)
			)
		`).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func seedProduct(db *gorm.DB) error {
	err := db.Exec(`
			insert into product (id, name, price, created_at)
			values  ('024b88ca-4767-4650-90ad-7aa5aa50950d', 'Sandal', 17500, default),
					('0edc594c-3af1-468a-9e01-064847444b34', 'Baju Taqwo', 57000, default),
					('3d3b4089-e439-4921-a727-b2a1e1932f35', 'Celana Renang', 250000, default),
					('4b6965d2-dd96-438f-a658-d0b5ddf843ac', 'Kacamata Anti Radiasi', 750000, default),
					('5846fa16-1730-4ddf-a483-a543c4562daa', 'Laptop Gaming', 20000000, default),
					('5beeb633-3e32-415f-baa8-a08d941182b5', 'Sepatu', 150000, default),
					('6ca0c5cc-1dd0-4b02-96b7-36cd3f797ad4', 'Kaos Kaki Lucu', 10000, default),
					('99437a64-d9bf-4039-9b20-c20b92f55067', 'Baju Olahraga', 325000, default),
					('ea659273-9456-4dd3-a4ae-a4faf852b736', 'Jaket Boomber', 900000, default),
					('f90ff503-62e3-42c3-9555-f3190b432a24', 'TWS', 2100000, default)
		`).Error
	if err != nil {
		return err
	}

	return nil
}

func isEmpty(db *gorm.DB, tableName string) bool {
	var count int64
	db.Table(tableName).Count(&count)
	return count == 0
}

func seedCategory(db *gorm.DB) error {
	err := db.Exec(`
			insert into category (id, name, created_at)
			values  ('280227cf-af40-4181-be90-267cf0471b05', 'Pelengkap', default),
					('464e64e5-1c60-4b21-99c8-f7cbcbd40c80', 'Teknologi', default),
					('6272c6d8-3269-453b-86b4-5795efa86c4b', 'Olahraga', default),
					('ff6d8dd3-7165-425a-99da-a6d7abee1dc4', 'Pakaian', default)
		`).Error
	if err != nil {
		return err
	}

	return nil
}

func seedProductCategory(db *gorm.DB) error {
	err := db.Exec(`
			insert into product_category (product_id, category_id)
			values  ('024b88ca-4767-4650-90ad-7aa5aa50950d', '280227cf-af40-4181-be90-267cf0471b05'),
					('4b6965d2-dd96-438f-a658-d0b5ddf843ac', '280227cf-af40-4181-be90-267cf0471b05'),
					('5beeb633-3e32-415f-baa8-a08d941182b5', '280227cf-af40-4181-be90-267cf0471b05'),
					('6ca0c5cc-1dd0-4b02-96b7-36cd3f797ad4', '280227cf-af40-4181-be90-267cf0471b05'),
					('4b6965d2-dd96-438f-a658-d0b5ddf843ac', '464e64e5-1c60-4b21-99c8-f7cbcbd40c80'),
					('5846fa16-1730-4ddf-a483-a543c4562daa', '464e64e5-1c60-4b21-99c8-f7cbcbd40c80'),
					('f90ff503-62e3-42c3-9555-f3190b432a24', '464e64e5-1c60-4b21-99c8-f7cbcbd40c80'),
					('3d3b4089-e439-4921-a727-b2a1e1932f35', '6272c6d8-3269-453b-86b4-5795efa86c4b'),
					('99437a64-d9bf-4039-9b20-c20b92f55067', '6272c6d8-3269-453b-86b4-5795efa86c4b'),
					('0edc594c-3af1-468a-9e01-064847444b34', 'ff6d8dd3-7165-425a-99da-a6d7abee1dc4'),
					('3d3b4089-e439-4921-a727-b2a1e1932f35', 'ff6d8dd3-7165-425a-99da-a6d7abee1dc4'),
					('99437a64-d9bf-4039-9b20-c20b92f55067', 'ff6d8dd3-7165-425a-99da-a6d7abee1dc4'),
					('ea659273-9456-4dd3-a4ae-a4faf852b736', 'ff6d8dd3-7165-425a-99da-a6d7abee1dc4')
		`).Error
	if err != nil {
		return err
	}

	return nil
}
