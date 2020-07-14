CREATE TABLE IF NOT EXISTS follower
(
      f_id INT NOT NULL UNIQUE AUTO_INCREMENT,
      user_id INT NOT NULL,
      follower_id INT NOT NULL,
      active BOOLEAN NOT NULL DEFAULT TRUE,
      PRIMARY KEY (f_id),
      FOREIGN KEY (user_id)
        REFERENCES user(u_id)
        ON DELETE CASCADE,
      FOREIGN KEY (follower_id)
            REFERENCES user(u_id)
            ON DELETE CASCADE
);