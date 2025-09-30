-- name: InsertEvent :one
INSERT INTO outbox_queue (
    key, schema_id, message
) VALUES (
             @key, @schema_id, @message
         ) returning id;


-- name: UpdateMessage :exec
UPDATE outbox_queue
SET message = @message
WHERE id = @id;

