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
