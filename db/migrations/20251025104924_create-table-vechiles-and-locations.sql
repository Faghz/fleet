-- migrate:up

-- Create vehicles table to store vehicle information
CREATE TABLE public.vehicle (
    entity_id uuid NOT NULL,
    vehicle_id TEXT NOT NULL,
    vehicle_type TEXT,
    brand TEXT,
    model TEXT,
    year integer,
    status TEXT DEFAULT 'active' NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_by TEXT NOT NULL,
    updated_at TIMESTAMPTZ,
    updated_by TEXT,
    deleted_at TIMESTAMPTZ,
    deleted_by TEXT,
    CONSTRAINT vehicle_pkey PRIMARY KEY (entity_id),
    CONSTRAINT vehicle_id_unique UNIQUE (vehicle_id)
);

-- Create vehicle_location table to store current and historical locations
CREATE TABLE public.vehicle_location (
    entity_id uuid NOT NULL,
    vehicle_entity_id uuid NOT NULL UNIQUE,
    latitude numeric(10, 8) NOT NULL,
    longitude numeric(11, 8) NOT NULL,
    timestamp bigint NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ,
    CONSTRAINT vehicle_location_pkey PRIMARY KEY (entity_id),
    CONSTRAINT vehicle_location_vehicle_entity_id_fkey FOREIGN KEY (vehicle_entity_id) REFERENCES public.vehicle(entity_id) ON DELETE CASCADE
);

CREATE TABLE public.vehicle_location_history (
    entity_id uuid NOT NULL,
    vehicle_entity_id uuid NOT NULL,
    latitude numeric(10, 8) NOT NULL,
    longitude numeric(11, 8) NOT NULL,
    timestamp bigint NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ,
    CONSTRAINT vehicle_location_history_pkey PRIMARY KEY (entity_id),
    CONSTRAINT vehicle_location_history_vehicle_entity_id_fkey FOREIGN KEY (vehicle_entity_id) REFERENCES public.vehicle(entity_id) ON DELETE CASCADE
);

-- Create index on vehicle_entity_id for faster lookups
CREATE INDEX idx_vehicle_location_vehicle_entity_id ON public.vehicle_location(vehicle_entity_id);

-- Create index on timestamp for historical queries
CREATE INDEX idx_vehicle_location_timestamp ON public.vehicle_location(timestamp DESC);

-- Create composite index for finding latest location per vehicle
CREATE INDEX idx_vehicle_location_vehicle_timestamp ON public.vehicle_location(vehicle_entity_id, timestamp DESC);

-- Create index on vehicle_entity_id for faster lookups on history
CREATE INDEX idx_vehicle_location_history_vehicle_entity_id ON public.vehicle_location_history(vehicle_entity_id);

-- Create index on timestamp for historical queries on history
CREATE INDEX idx_vehicle_location_history_timestamp ON public.vehicle_location_history(timestamp DESC);

-- Create composite index for finding latest history per vehicle
CREATE INDEX idx_vehicle_location_history_vehicle_timestamp ON public.vehicle_location_history(vehicle_entity_id, timestamp DESC);

-- Add trigger to auto-update updated_at on vehicle table
CREATE TRIGGER update_vehicle_updated_at
    BEFORE UPDATE ON public.vehicle
    FOR EACH ROW
    EXECUTE FUNCTION public.update_updated_at_column();

CREATE TRIGGER update_vehicle_location_updated_at
    BEFORE UPDATE ON public.vehicle_location
    FOR EACH ROW
    EXECUTE FUNCTION public.update_updated_at_column();   

-- migrate:down

DROP TRIGGER IF EXISTS update_vehicle_updated_at ON public.vehicle;
DROP TRIGGER IF EXISTS update_vehicle_location_updated_at ON public.vehicle_location;
DROP INDEX IF EXISTS idx_vehicle_location_vehicle_timestamp;
DROP INDEX IF EXISTS idx_vehicle_location_timestamp;
DROP INDEX IF EXISTS idx_vehicle_location_vehicle_id;
DROP INDEX IF EXISTS idx_vehicle_location_history_vehicle_timestamp;
DROP INDEX IF EXISTS idx_vehicle_location_history_timestamp;
DROP INDEX IF EXISTS idx_vehicle_location_history_vehicle_entity_id;
DROP TABLE IF EXISTS public.vehicle_location_history;
DROP TABLE IF EXISTS public.vehicle_location;
DROP TABLE IF EXISTS public.vehicle;
