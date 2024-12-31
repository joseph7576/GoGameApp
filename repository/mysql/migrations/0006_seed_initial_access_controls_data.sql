-- +migrate Up
INSERT INTO 
    `access_controls`(`id`, `actor_type`, `actor_id`, `permission_id`) 
VALUES
    (1, "role", 2, 2),
    (2, "role", 2, 1);

-- +migrate Down
DELETE FROM
    `access_controls`
WHERE
    `id` IN (1, 2);