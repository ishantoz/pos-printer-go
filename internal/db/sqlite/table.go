package sqlite

// + sqlite-migrate
const BarcodeJobTableStmt = `CREATE TABLE IF NOT EXISTS barcode_jobs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		vid TEXT, pid TEXT,
		sizeX INTEGER, sizeY INTEGER,
		direction INTEGER, topText TEXT,
		barcodeData TEXT, printCount INTEGER,
		status TEXT, attempts INTEGER,
		createdAt DATETIME, updatedAt DATETIME
	);`

// + sqlite-migrate
const ReceiptPDFJobTableStmt = `CREATE TABLE IF NOT EXISTS receipt_pdf_jobs (
		id               INTEGER PRIMARY KEY AUTOINCREMENT,
		file_path        TEXT    NOT NULL,
		print_count      INTEGER DEFAULT 1,
		connection_type  TEXT    NOT NULL CHECK(connection_type IN ('network','usb')),
		printer_ip       TEXT,
		printer_port     INTEGER,
		usb_vendor_id    INTEGER,
		usb_product_id   INTEGER,
		usb_interface    INTEGER DEFAULT 0,
		printer_width    INTEGER DEFAULT 576,
		threshold        INTEGER DEFAULT 100,
		feed_lines       INTEGER DEFAULT 1,
		zoom             REAL    DEFAULT 2.0,
		status           TEXT    DEFAULT 'pending',
		retry_count      INTEGER DEFAULT 0,
		last_error       TEXT,
		created_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP
  );`
