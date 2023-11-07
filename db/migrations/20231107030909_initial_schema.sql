-- migrate:up
-- Add PostGIS extension if it's not already created
CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create UserPositions Table
CREATE TABLE user_positions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    user_id UUID NOT NULL,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    place_id UUID,
    checked_in TIMESTAMP WITH TIME ZONE,
    checked_out TIMESTAMP WITH TIME ZONE,
    -- PostGIS geographical point for location (latitude and longitude)
    location GEOGRAPHY(Point, 4326)
);
-- Create indexes for user_id and timestamp
CREATE INDEX idx_user_positions_user_id ON user_positions (user_id);
CREATE INDEX idx_user_positions_created_at ON user_positions (created_at);
-- Create index for location
CREATE INDEX idx_user_positions_location ON user_positions USING GIST (location);
-- Create PhoneMetadata Table
CREATE TABLE phone_metadata (
    user_position_id UUID PRIMARY KEY,
    device_id VARCHAR(255) NOT NULL,
    model VARCHAR(255) NOT NULL,
    os_version VARCHAR(255) NOT NULL,
    carrier VARCHAR(255) NOT NULL,
    corporate_id VARCHAR(255) NOT NULL,
    FOREIGN KEY (user_position_id) REFERENCES user_positions(id)
);
-- Create index for device_id
CREATE INDEX idx_phone_metadata_device_id ON phone_metadata (device_id);
-- migrate:down
-- Drop the tables and indexes
DROP TABLE IF EXISTS phone_metadata;
DROP TABLE IF EXISTS user_positions;
DROP EXTENSION IF EXISTS postgis;