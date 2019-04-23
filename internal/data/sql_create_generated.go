package data

const SQLCreate = `
PRAGMA foreign_keys = ON;
PRAGMA encoding = 'UTF-8';
-- %MYAPPDATA%\vimec\vimec.sqlite

CREATE TABLE IF NOT EXISTS act
(
    act_id          INTEGER PRIMARY KEY NOT NULL,
    created_at      TIMESTAMP           NOT NULL DEFAULT (datetime('now')),
    act_number      INTEGER             NOT NULL,
    doc_code        INTEGER             NOT NULL CHECK ( doc_code IN (2, 4, 6) ),
    route_sheet     TEXT                NOT NULL CHECK ( route_sheet != '' ),
    products_count INTEGER             NOT NULL CHECK ( products_count > 0 )
);`
