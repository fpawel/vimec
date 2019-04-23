PRAGMA foreign_keys = ON;
PRAGMA encoding = 'UTF-8';
-- %MYAPPDATA%\vimec\vimec.sqlite

CREATE TABLE IF NOT EXISTS act
(
    act_id         INTEGER PRIMARY KEY NOT NULL,
    year           INTEGER             NOT NULL CHECK (year > 2000),
    month          INTEGER             NOT NULL CHECK (month BETWEEN 1 and 12),
    day            INTEGER             NOT NULL CHECK (month BETWEEN 1 and 31),
    act_number     INTEGER             NOT NULL,
    doc_code       INTEGER             NOT NULL CHECK ( doc_code IN (2, 4, 6) ),
    route_sheet    TEXT                NOT NULL CHECK ( route_sheet != '' ),
    products_count INTEGER             NOT NULL CHECK ( products_count > 0 ),
    UNIQUE (act_number, year),
    UNIQUE (act_number, route_sheet),
    UNIQUE (route_sheet, year)
);