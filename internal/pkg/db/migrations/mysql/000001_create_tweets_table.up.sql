CREATE TABLE IF NOT EXISTS tweet
(
      t_id INT NOT NULL UNIQUE AUTO_INCREMENT,
      u_id INT NOT NULL,
      created_on TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
      content VARCHAR(280) NOT NULL,
      PRIMARY KEY (t_id),
      FOREIGN KEY (u_id)
        REFERENCES user(u_id)
        ON DELETE CASCADE
);