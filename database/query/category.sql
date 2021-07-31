/*
    table rows
    id | parentID | name
    1 |           | Main
    2 |     1     | SubCategory1
    3 |     1     | SubCategory2

    json object
    {
        id : "1",
        name: "Main",
        categories : [
            {id : "123", name: "SubCategory1"},
            {id : "234", name: "SubCategory2"}
        ]
    }

*/

--name: CreateCategory :one
INSERT INTO categories(parent_id, user_id, name)
VALUES ($1, $2, $3)
RETURNING *;

--name: UpdateCategory :one
UPDATE categories SET name = $3
AND parent_id = $2
WHERE category_id = $1 AND deleted_at = '0001-01-01 00:00:00Z'
RETURNING *;

--name: GetCategoryByID :one
SELECT * FROM categories
WHERE category_id = $1 AND deleted_at = '0001-01-01 00:00:00Z'
LIMIT 1;

--name: ListCategories :many
SELECT * FROM categories
WHERE user_id = $1 AND deleted_at = '0001-01-01 00:00:00Z'
ORDER BY category_id
LIMIT $2
OFFSET $3;

--name: DeleteCategory :one
UPDATE categories SET deleted_at = now()
WHERE category_id = $1
AND deleted_at = '0001-01-01 00:00:00Z'
RETURNING deleted_at;