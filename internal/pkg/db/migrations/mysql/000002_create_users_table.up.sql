CREATE TABLE IF NOT EXISTS user
(
      u_id INT NOT NULL UNIQUE AUTO_INCREMENT,
      f_name VARCHAR(150) NOT NULL,
      l_name VARCHAR(150),
      email VARCHAR(200) NOT NULL,
      password VARCHAR(200) NOT NULL,
      user_name VARCHAR(50) NOT NULL,
      deleted BOOLEAN NOT NULL DEFAULT FALSE,
      PRIMARY KEY (u_id)
);