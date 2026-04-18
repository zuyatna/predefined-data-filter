CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

CREATE TABLE colors (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL
);

CREATE TABLE labels (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL
);

CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    category_id INT REFERENCES categories(id) ON DELETE SET NULL,
    price DECIMAL(12, 2) NOT NULL,
    purchases_count INT DEFAULT 0,
    reviews_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE product_colors (
    product_id INT REFERENCES products(id) ON DELETE CASCADE,
    color_id INT REFERENCES colors(id) ON DELETE CASCADE,
    PRIMARY KEY (product_id, color_id)
);

CREATE TABLE product_labels (
    product_id INT REFERENCES products(id) ON DELETE CASCADE,
    label_id INT REFERENCES labels(id) ON DELETE CASCADE,
    PRIMARY KEY (product_id, label_id)
);

-- SEED DATA
INSERT INTO categories (name) VALUES ('Sepatu'), ('Sandal'), ('Baju'), ('Celana');
INSERT INTO colors (name) VALUES ('Hitam'), ('Putih'), ('Merah'), ('Biru'), ('Hijau');
INSERT INTO labels (name) VALUES ('New'), ('Sale'), ('Best Seller'), ('Limited');

-- Sepatu
INSERT INTO products (name, category_id, price, purchases_count, reviews_count, created_at) VALUES 
('Sneakers X', 1, 500000, 150, 20, '2023-01-10 10:00:00'),
('Running Shoes Pro', 1, 750000, 500, 80, '2023-02-15 14:30:00'),
('Boots Y', 1, 1200000, 50, 5, '2023-05-20 09:15:00');

-- Sandal
INSERT INTO products (name, category_id, price, purchases_count, reviews_count, created_at) VALUES 
('Sandal Jepit Z', 2, 50000, 1000, 150, '2023-06-01 11:00:00'),
('Sandal Gunung A', 2, 250000, 300, 45, '2023-03-12 16:20:00');

-- Baju
INSERT INTO products (name, category_id, price, purchases_count, reviews_count, created_at) VALUES 
('Kaos Polos', 3, 75000, 2000, 350, '2023-01-05 08:00:00'),
('Kemeja Flanel', 3, 200000, 400, 60, '2023-04-10 13:45:00'),
('Jaket Hoodie', 3, 350000, 800, 120, '2023-07-22 10:10:00');

-- Celana
INSERT INTO products (name, category_id, price, purchases_count, reviews_count, created_at) VALUES 
('Celana Jeans', 4, 300000, 1200, 180, '2023-02-28 09:30:00'),
('Celana Chino', 4, 250000, 600, 90, '2023-05-05 14:00:00');

-- Relasi Product Colors
INSERT INTO product_colors (product_id, color_id) VALUES 
(1, 1), (1, 2), -- Sneakers X (Hitam, Putih)
(2, 3), (2, 4), -- Running Shoes Pro (Merah, Biru)
(3, 1),         -- Boots Y (Hitam)
(4, 1), (4, 4), (4, 5), -- Sandal Jepit Z (Hitam, Biru, Hijau)
(5, 1),         -- Sandal Gunung A (Hitam)
(6, 1), (6, 2), (6, 3), (6, 4), -- Kaos Polos (Hitam, Putih, Merah, Biru)
(7, 3),         -- Kemeja Flanel (Merah)
(8, 1), (8, 2), -- Jaket Hoodie (Hitam, Putih)
(9, 4),         -- Celana Jeans (Biru)
(10, 1), (10, 2); -- Celana Chino (Hitam, Putih)

-- Relasi Product Labels
INSERT INTO product_labels (product_id, label_id) VALUES 
(1, 1),         -- Sneakers X (New)
(2, 2), (2, 3), -- Running Shoes Pro (Sale, Best Seller)
(3, 4),         -- Boots Y (Limited)
(4, 3),         -- Sandal Jepit Z (Best Seller)
(5, 2),         -- Sandal Gunung A (Sale)
(6, 3),         -- Kaos Polos (Best Seller)
(7, 1),         -- Kemeja Flanel (New)
(8, 3),         -- Jaket Hoodie (Best Seller)
(9, 3),         -- Celana Jeans (Best Seller)
(10, 2);        -- Celana Chino (Sale)
