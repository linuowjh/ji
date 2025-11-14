USE yun_nian_memorial;

CREATE TABLE IF NOT EXISTS service_reviews (
  id VARCHAR(36) PRIMARY KEY,
  user_id VARCHAR(36) NOT NULL,
  booking_id VARCHAR(36) NOT NULL UNIQUE,
  rating INT NOT NULL,
  comment TEXT,
  tags JSON,
  is_anonymous BOOLEAN DEFAULT FALSE,
  created_at DATETIME(3),
  INDEX idx_service_reviews_user_id (user_id),
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (booking_id) REFERENCES service_bookings(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS service_staff (
  id VARCHAR(36) PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  role VARCHAR(50) NOT NULL,
  specialties JSON,
  avatar_url VARCHAR(255),
  bio TEXT,
  rating DECIMAL(3,2) DEFAULT 5.0,
  review_count INT DEFAULT 0,
  is_available BOOLEAN DEFAULT TRUE,
  created_at DATETIME(3),
  updated_at DATETIME(3)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
