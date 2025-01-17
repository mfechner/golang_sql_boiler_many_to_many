BEGIN;

CREATE TABLE IF NOT EXISTS admin
(
    id BIGINT(20) NOT NULL AUTO_INCREMENT,
    username VARCHAR(255) DEFAULT NULL,
    password VARCHAR(255) NOT NULL,
    super TINYINT(1) NOT NULL,
    active TINYINT(1) NOT NULL,
    created DATETIME NOT NULL,
    modified DATETIME DEFAULT NULL,
    PRIMARY KEY (id),
    UNIQUE KEY IX_Username_1 (username)
);

CREATE TABLE IF NOT EXISTS domain (
    id BIGINT(20) NOT NULL AUTO_INCREMENT,
    domain VARCHAR(255) NOT NULL,
    description VARCHAR(255) DEFAULT NULL,
    max_aliases INT(11) NOT NULL,
    alias_count BIGINT(20) NOT NULL,
    max_mailboxes INT(11) NOT NULL,
    mailbox_count BIGINT(20) NOT NULL,
    max_quota BIGINT(20) NOT NULL,
    quota BIGINT(20) NOT NULL,
    transport VARCHAR(255) NOT NULL,
    backupmx TINYINT(1) NOT NULL,
    active TINYINT(1) NOT NULL,
    homedir VARCHAR(255) DEFAULT NULL,
    maildir VARCHAR(255) DEFAULT NULL,
    uid INT(11) DEFAULT NULL,
    gid INT(11) DEFAULT NULL,
    created DATETIME NOT NULL,
    modified DATETIME DEFAULT NULL,
    PRIMARY KEY (id),
    UNIQUE KEY IX_Domain_1 (domain)
);

CREATE TABLE IF NOT EXISTS domain_admins (
    Admin_id BIGINT(20) NOT NULL,
    Domain_id BIGINT(20) NOT NULL,
    PRIMARY KEY (Admin_id,Domain_id),
    KEY IDX_CD8319C69D5DE046 (Admin_id),
    KEY IDX_CD8319C693AE8C46 (Domain_id)
);

ALTER TABLE domain_admins ADD CONSTRAINT FK_CD8319C693AE8C46 FOREIGN KEY IF NOT EXISTS (Domain_id) REFERENCES domain (id),
                          ADD CONSTRAINT FK_CD8319C69D5DE046 FOREIGN KEY IF NOT EXISTS (Admin_id) REFERENCES admin (id);

COMMIT;
