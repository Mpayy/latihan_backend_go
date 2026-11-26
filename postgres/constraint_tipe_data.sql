-- Topik 5: Constraint & tipe data

-- Konsep: Constraint memindahkan aturan data ke database, jadi data tetap valid meski ada bug di aplikasi atau ada proses lain yang menulis: PRIMARY KEY (unik + NOT NULL), FOREIGN KEY (relasi harus ada induknya), UNIQUE, CHECK (kondisi bebas), NOT NULL. Tipe data menentukan ketepatan: BIGINT untuk ID/hitungan besar, NUMERIC(p,s) untuk uang (eksak), TIMESTAMPTZ untuk waktu (disimpan sebagai UTC).

-- Contoh:

-- go
-- const ddl = `
-- CREATE TABLE wallets (
--   id         BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
--   user_id    BIGINT NOT NULL REFERENCES users(id),
--   balance    NUMERIC(18,2) NOT NULL DEFAULT 0 CHECK (balance >= 0),
--   created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
--   UNIQUE (user_id)
-- )`

-- Latihan:

-- Rancang tabel orders(id, user_id, total, status, created_at). Tentukan tipe, NOT NULL, FK, dan satu CHECK untuk total. Tambahkan CHECK agar status hanya boleh 'pending', 'paid', atau 'cancelled'.

CREATE TABLE IF NOT EXISTS orders (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    total NUMERIC(18,2) NOT NULL DEFAULT 0 CHECK (total >= 0),
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'paid', 'cancelled'))
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
);

CREATE INDEX idx_orders_user ON orders(user_id);

-- Kenapa TIMESTAMP (tanpa TZ) berisiko dibanding TIMESTAMPTZ?
-- Karena TIMESTAMP tidak menyimpan zona waktu, jadi waktu di tulis apa adanya di db. jadi bisa terjadi abigu saat kita menarik datanya. TIMESTAMPTZ menyimpan zona waktu, jadi saat menyimpan akan di convert ke UTC 

-- Di wallet-mu, CHECK (balance >= 0) ada. Kenapa itu tetap berguna meski aplikasi sudah mengecek saldo sebelum debit?
-- untuk menjaga konsisten data level database, jika terdapat bug atau race condition di aplikasi
-- jika dua request debit bersamaan sama-sama lolos cek saldo di aplikasi (race), lalu keduanya menulis. CHECK menjadi jaring pengaman terakhir: salah satu UPDATE gagal, bukan saldo jadi negatif. Tapi CHECK bukan pengganti locking. Ia hanya mencegah data rusak, tidak mengatur urutan.