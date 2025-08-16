


-- name: CreateGig :one
INSERT INTO gigs (
  user_id, title,description, pricing_mode, category_id, image_url, status, updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, NOW()
)
RETURNING *;

-- name: CreateGigPackage :one
INSERT INTO gig_packages (
      gig_id, tier, price, delivery_time, options
) VALUES (
      $1, $2, $3, $4, $5
)
RETURNING *;

-- name: CreateGigRequirement :one
INSERT INTO gig_requirements (
      gig_id, question, required
) VALUES (
      $1, $2, $3
    )
    RETURNING *;

-- name: UpdateService :one
UPDATE gigs
SET
    title = $2,
    description = $3,
    category_id = $4,
    image_url = $5,
    pricing_mode = $6,
    status = $7,
    updated_at = NOW()
WHERE id = $1
RETURNING *;



-- name: DeleteService :exec
DELETE FROM gigs WHERE id = $1;


-- name: GetService :one
SELECT
    g.id,
    g.user_id,
    g.title,
    g.description,
    g.category_id,
    g.image_url,
    g.pricing_mode,
    g.status,
    g.created_at,
    g.updated_at,
    array_remove(array_agg(c.name), NULL)::text[] AS category_name,
    u.names AS user_name,
    u.profile_pic AS user_profile_pic
FROM gigs g
JOIN users u ON g.user_id = u.id
JOIN categories c ON c.id = ANY(g.category_id)
WHERE g.id = $1
GROUP BY g.id, g.user_id, g.title, g.description, g.category_id,
         g.image_url, g.pricing_mode, g.status, g.created_at,
         g.updated_at, u.names, u.profile_pic;



-- name: GetGigPackagesByGigID :many
SELECT
        id,
        gig_id,
        tier,
        price,
        delivery_time,
        options
FROM
        gig_packages
WHERE
        gig_id = $1;


-- name: GetGigRequirementsByGigID :many
SELECT
        id,
        gig_id,
        question,
        required
FROM
        gig_requirements
WHERE
        gig_id = $1;




-- name: SearchGigs :many
SELECT
    g.id,
    g.title,
    g.description,
    g.image_url,
    g.pricing_mode,
    g.created_at,
    g.updated_at,
    COALESCE(bp.price, 0) AS basic_price
FROM gigs g
LEFT JOIN gig_packages bp 
    ON bp.gig_id = g.id 
    AND bp.tier = 'basic'
WHERE
    (
        LOWER(g.title) LIKE LOWER('%' || sqlc.arg(search_term)::text || '%') OR
        LOWER(g.description) LIKE LOWER('%' || sqlc.arg(search_term)::text || '%') OR
        sqlc.arg(search_term)::text IS NULL
    )
    AND (
        (sqlc.arg(min_price)::float8 = 0 AND sqlc.arg(max_price)::float8 = 0) OR -- Thêm điều kiện này
        (bp.price >= sqlc.arg(min_price)::float8 OR sqlc.arg(min_price)::float8 IS NULL)
    )
    AND (
        (sqlc.arg(max_price)::float8 = 0 AND sqlc.arg(min_price)::float8 = 0) OR -- Thêm điều kiện này
        (bp.price <= sqlc.arg(max_price)::float8 OR sqlc.arg(max_price)::float8 IS NULL)
    )
    AND (
        array_length(sqlc.arg(category_ids)::int[], 1) IS NULL OR sqlc.arg(category_ids)::int[] IS NULL OR g.category_id && sqlc.arg(category_ids)::int[]
    )
    AND (
        sqlc.arg(last_gig_id)::uuid = '00000000-0000-0000-0000-000000000000' OR sqlc.arg(last_gig_id)::uuid IS NULL OR g.id < sqlc.arg(last_gig_id)::uuid
    )
ORDER BY g.id DESC
LIMIT 10;




-- name: GetCategories :many
select 
c1.id as parent_id,
c1.name as parent_name,
c2.id as children_id,
c2.name as children_name
from categories c1
left join categories c2 on c1.id = c2.parent_id;

-- name: GetGigAndPackagesForOrder :one
SELECT
    g.id,
    g.user_id,
    json_agg(
        json_build_object(
            'tier', gp.tier,
            'price', gp.price
        )
    ) FILTER (WHERE gp.id IS NOT NULL) AS gig_packages
FROM gigs g
LEFT JOIN gig_packages gp ON g.id = gp.gig_id
WHERE g.id = $1
GROUP BY g.id;
