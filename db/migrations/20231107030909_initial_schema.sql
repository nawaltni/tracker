-- migrate:up
-- Add PostGIS extension if it's not already created
CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- Create UserPositions Table
CREATE TABLE user_positions (
    user_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    place_id UUID NULL,
    place_name TEXT NULL,
    checked_in TIMESTAMP WITH TIME ZONE NULL,
    checked_out TIMESTAMP WITH TIME ZONE NULL,
    -- PostGIS geographical point for location (latitude and longitude)
    location GEOMETRY(Point, 4326)
);
-- Create indexes for user_id and timestamp
CREATE INDEX idx_user_positions_created_at ON user_positions (created_at);
-- Create index for location
CREATE INDEX idx_user_positions_location ON user_positions USING GIST (location);
-- Create PhoneMetadata Table
CREATE TABLE phone_metadata (
    user_id UUID PRIMARY KEY REFERENCES user_positions(user_id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    device_id VARCHAR(255) NOT NULL,
    model VARCHAR(255) NOT NULL,
    os_version VARCHAR(255) NOT NULL,
    carrier VARCHAR(255) NOT NULL,
    corporate_id VARCHAR(255) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES user_positions(user_id) ON DELETE CASCADE
);
-- Create index for device_id
CREATE INDEX idx_phone_metadata_device_id ON phone_metadata (device_id);
-- migrate:down
-- Drop the tables and indexes
DROP TABLE IF EXISTS phone_metadata;
DROP TABLE IF EXISTS user_positions;
DROP EXTENSION IF EXISTS postgis;