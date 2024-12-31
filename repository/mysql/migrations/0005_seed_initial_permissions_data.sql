-- +migrate Up
INSERT INTO 
    `permissions`(`id`, `title`) 
VALUES
    (1, "user-list"), 
    (2, "user-delete");

-- +migrate Down
DELETE FROM
    `permissions`
WHERE
    `id` IN (1, 2);