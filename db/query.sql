-- name: CreateRecipe :one
INSERT INTO recipes (title, servings, servings_type, country_code, host_url, original_url)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, created_at, updated_at;

-- name: CreateIngredient :one
INSERT INTO ingredients (name)
VALUES ($1)
RETURNING id;

-- name: CreateRecipeIngredient :one
INSERT INTO recipe_ingredients (recipe_id, ingredient_id, amount, unit)
VALUES ($1, $2, $3, $4)
RETURNING recipe_id, ingredient_id;

-- name: CreateInstruction :one
INSERT INTO instructions (recipe_id, step_number, instruction_text)
VALUES ($1, $2, $3)
RETURNING id;

