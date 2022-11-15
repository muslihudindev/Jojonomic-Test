CREATE TABLE IF NOT EXISTS harga (
    reff_id VARCHAR (15) PRIMARY KEY NOT NULL,
    admin_id VARCHAR (15) NOT NULL,
    harga_topup FLOAT DEFAULT 0,
    harga_buyback FLOAT DEFAULT 0,
    date INTEGER
);

CREATE TABLE IF NOT EXISTS rekening (
    reff_id VARCHAR (15)  PRIMARY KEY NOT NULL,
    norek VARCHAR (15) UNIQUE NOT NULL,
    saldo DECIMAL(12, 2) DEFAULT 0,
    date INTEGER
);

CREATE TABLE if NOT EXISTS transaksi (
    reff_id VARCHAR (15)  PRIMARY KEY NOT NULL,
    norek VARCHAR (15),
    type VARCHAR(15),
    gram FLOAT DEFAULT 0,
    harga_topup FLOAT DEFAULT 0,
    harga_buyback FLOAT DEFAULT 0,
    saldo FLOAT DEFAULT 0,
    date INTEGER
);