-- +goose Up
-- +goose StatementBegin
-- Create roles table
CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
CREATE INDEX idx_roles_deleted_at ON roles(deleted_at);
CREATE INDEX idx_role_name ON roles(name);
-- Seed initial roles
INSERT INTO roles (id, name)
VALUES (1, 'User'),
    (2, 'Moderator'),
    (3, 'Admin'),
    (4, 'SuperAdmin') ON CONFLICT (id) DO NOTHING;
-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    email_verifyed BOOLEAN DEFAULT FALSE,
    full_name VARCHAR(255),
    role_id INTEGER REFERENCES roles(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
-- Create images table
CREATE TABLE IF NOT EXISTS images (
    id SERIAL PRIMARY KEY,
    path VARCHAR(500),
    url VARCHAR(500),
    size DOUBLE PRECISION,
    width INTEGER,
    height INTEGER,
    type VARCHAR(100),
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
CREATE INDEX idx_images_deleted_at ON images(deleted_at);
CREATE INDEX idx_images_user_id ON images(user_id);
-- Create profile_images table
CREATE TABLE IF NOT EXISTS profile_images (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    image_id INTEGER REFERENCES images(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
CREATE INDEX idx_profile_images_deleted_at ON profile_images(deleted_at);
CREATE INDEX idx_profile_images_user_id ON profile_images(user_id);
-- Create team_roles table
CREATE TABLE IF NOT EXISTS team_roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
CREATE INDEX idx_team_roles_deleted_at ON team_roles(deleted_at);
-- Seed team roles
INSERT INTO team_roles (id, name)
VALUES (1, 'Owner'),
    (2, 'Admin'),
    (3, 'Member') ON CONFLICT (id) DO NOTHING;
-- Create teams table
CREATE TABLE IF NOT EXISTS teams (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    address TEXT,
    phone VARCHAR(50),
    email VARCHAR(255),
    username VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
CREATE INDEX idx_teams_deleted_at ON teams(deleted_at);
CREATE INDEX idx_teams_username ON teams(username);
-- Create team_members table (composite primary key)
CREATE TABLE IF NOT EXISTS team_members (
    id SERIAL PRIMARY KEY,
    team_id INTEGER NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    team_role_id INTEGER REFERENCES team_roles(id),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    UNIQUE(team_id, user_id)
);
CREATE INDEX idx_team_members_deleted_at ON team_members(deleted_at);
CREATE INDEX idx_team_members_team_id ON team_members(team_id);
CREATE INDEX idx_team_members_user_id ON team_members(user_id);
-- Create products table
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    team_id INTEGER REFERENCES teams(id) ON DELETE CASCADE,
    code VARCHAR(100),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    uom VARCHAR(50),
    price BIGINT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
CREATE INDEX idx_products_deleted_at ON products(deleted_at);
CREATE INDEX idx_products_team_id ON products(team_id);
-- Create product_categories table
CREATE TABLE IF NOT EXISTS product_categories (
    id SERIAL PRIMARY KEY,
    team_id INTEGER REFERENCES teams(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
CREATE INDEX idx_product_categories_deleted_at ON product_categories(deleted_at);
CREATE INDEX idx_product_categories_team_id ON product_categories(team_id);
-- Create product_product_categories junction table
CREATE TABLE IF NOT EXISTS product_product_categories (
    id SERIAL PRIMARY KEY,
    product_id INTEGER REFERENCES products(id) ON DELETE CASCADE,
    product_category_id INTEGER REFERENCES product_categories(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    UNIQUE(product_id, product_category_id)
);
CREATE INDEX idx_product_product_categories_deleted_at ON product_product_categories(deleted_at);
-- Create product_images table
CREATE TABLE IF NOT EXISTS product_images (
    id SERIAL PRIMARY KEY,
    product_id INTEGER REFERENCES products(id) ON DELETE CASCADE,
    image_id INTEGER REFERENCES images(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
CREATE INDEX idx_product_images_deleted_at ON product_images(deleted_at);
CREATE INDEX idx_product_images_product_id ON product_images(product_id);
-- Create auto_increment_sequences table
CREATE TABLE IF NOT EXISTS auto_increment_sequences (
    id SERIAL PRIMARY KEY,
    entity_type VARCHAR(50) NOT NULL,
    team_id INTEGER NOT NULL,
    entity_id INTEGER NOT NULL,
    sequence INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    UNIQUE(entity_type, team_id, entity_id)
);
CREATE INDEX idx_auto_increment_sequences_deleted_at ON auto_increment_sequences(deleted_at);
CREATE INDEX idx_sequence_entity_team ON auto_increment_sequences(entity_type, team_id, entity_id);
-- Create report_json_schema_types table
CREATE TABLE IF NOT EXISTS report_json_schema_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
CREATE INDEX idx_report_json_schema_types_deleted_at ON report_json_schema_types(deleted_at);
-- Seed report json schema types
INSERT INTO report_json_schema_types (id, name)
VALUES (1, 'Common') ON CONFLICT (id) DO NOTHING;
-- Create report_templates table
CREATE TABLE IF NOT EXISTS report_templates (
    id SERIAL PRIMARY KEY,
    code VARCHAR(100),
    name VARCHAR(255) NOT NULL,
    display_name VARCHAR(255),
    icon VARCHAR(255),
    report_json_schema_type_id INTEGER REFERENCES report_json_schema_types(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
CREATE INDEX idx_report_templates_deleted_at ON report_templates(deleted_at);
-- Create table
CREATE TABLE IF NOT EXISTS casbin_rule (
    p_type VARCHAR(32),
    v0 VARCHAR(255),
    v1 VARCHAR(255),
    v2 VARCHAR(255),
    v3 VARCHAR(255),
    v4 VARCHAR(255),
    v5 VARCHAR(255)
);
-- Seed Casbin rule
INSERT INTO casbin_rule (p_type, v0, v1, v2, v3, v4, v5)
VALUES ('p', '1', '*', '*', '', '', ''),
    ('p', '2', '*', '*', '', '', ''),
    ('p', '3', '/teams/*', 'GET', '', '', ''),
    ('p', '3', '/teams', 'POST', '', '', ''),
    ('p', '3', '/teams/join/link', 'POST', '', '', ''),
    (
        'p',
        '3',
        '/teams/{id}/request-join',
        'POST',
        '',
        '',
        ''
    ),
    (
        'p',
        '3',
        '/teams/{id}/products/*',
        'GET',
        '',
        '',
        ''
    ),
    (
        'p',
        '3',
        '/teams/{id}/product_category/*',
        'GET',
        '',
        '',
        ''
    );
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS report_templates;
DROP TABLE IF EXISTS report_json_schema_types;
DROP TABLE IF EXISTS auto_increment_sequences;
DROP TABLE IF EXISTS product_images;
DROP TABLE IF EXISTS product_product_categories;
DROP TABLE IF EXISTS product_categories;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS team_members;
DROP TABLE IF EXISTS teams;
DROP TABLE IF EXISTS team_roles;
DROP TABLE IF EXISTS profile_images;
DROP TABLE IF EXISTS images;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS casbin_rule;
-- +goose StatementEnd