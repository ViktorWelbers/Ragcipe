CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS countries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL UNIQUE,
    code CHAR(2) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS sources (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    url VARCHAR(2048) NOT NULL,
    UNIQUE(url)
);

CREATE TABLE IF NOT EXISTS recipes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    servings INTEGER NOT NULL,
    servings_type VARCHAR(50) NOT NULL,
    country_id UUID REFERENCES countries (id),
    source_id UUID REFERENCES sources (id),
    original_url TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS ingredients (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS recipe_ingredients (
    recipe_id UUID REFERENCES recipes(id) ON DELETE CASCADE,
    ingredient_id UUID REFERENCES ingredients(id) ON DELETE CASCADE,
    amount DECIMAL(10, 2),
    unit VARCHAR(50),
    PRIMARY KEY (recipe_id, ingredient_id)
);

CREATE TABLE IF NOT EXISTS instructions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    recipe_id UUID REFERENCES recipes(id) ON DELETE CASCADE,
    step_number INTEGER NOT NULL,
    instruction_text TEXT NOT NULL,
    UNIQUE (recipe_id, step_number)
);

CREATE INDEX IF NOT EXISTS idx_recipe_ingredients_recipe_id 
    ON recipe_ingredients(recipe_id);
CREATE INDEX IF NOT EXISTS idx_recipe_ingredients_ingredient_id 
    ON recipe_ingredients(ingredient_id);
CREATE INDEX IF NOT EXISTS idx_instructions_recipe_id 
    ON instructions(recipe_id);
CREATE INDEX IF NOT EXISTS idx_ingredients_name 
    ON ingredients(name);
CREATE INDEX IF NOT EXISTS idx_recipes_country_id 
    ON recipes(country_id);
CREATE INDEX IF NOT EXISTS idx_recipes_source_id 
    ON recipes(source_id);
