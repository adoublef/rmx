CREATE TABLE `rooms` (
    -- xid format
    `id` TEXT,
    `title` TEXT NOT NULL,
    `capacity` INT NOT NULL,
    -- xid format 
    -- not a reference but needs to be treated as such
    `owner` TEXT NOT NULL,
    -- constraints
    CONSTRAINT title_bound_len CHECK (title <> '' AND length(title) <= 60),
    CONSTRAINT cap_bound_size CHECK (capacity > 0 AND capacity <= 8),
    PRIMARY KEY (id)
);

-- handle deletes by using tombstone 