CREATE SCHEMA dbmigration;

-- MIGRATION
CREATE TABLE dbmigration.migration (
  id_migration      TEXT NOT NULL,
  "user"            TEXT DEFAULT user,
  executed_at       TIMESTAMP DEFAULT NOW(),
  CONSTRAINT dbmigration_id__pkey PRIMARY KEY (id_migration)
);


-- HISTORY
CREATE TABLE dbmigration.migration_history (LIKE dbmigration.migration);
ALTER TABLE dbmigration.migration_history ADD COLUMN _operation TEXT NOT NULL;
ALTER TABLE dbmigration.migration_history ADD COLUMN "_user" TEXT NOT NULL;
ALTER TABLE dbmigration.migration_history ADD COLUMN "_operation_at" TIMESTAMP DEFAULT NOW();

CREATE OR REPLACE FUNCTION function_migration_history() RETURNS TRIGGER AS $$
BEGIN
    IF (TG_OP = 'DELETE') THEN
        INSERT INTO dbmigration.migration_history VALUES(OLD.*, 'D', user, now());
        RETURN OLD;
    ELSIF (TG_OP = 'UPDATE') THEN
        INSERT INTO dbmigration.migration_history VALUES(NEW.*, 'U', user, now());
        RETURN NEW;
    ELSIF (TG_OP = 'INSERT') THEN
        INSERT INTO dbmigration.migration_history VALUES(NEW.*, 'I', user, now());
        RETURN NEW;
    END IF;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_migration_history
AFTER INSERT OR UPDATE OR DELETE ON dbmigration.migration
    FOR EACH ROW EXECUTE PROCEDURE function_migration_history();