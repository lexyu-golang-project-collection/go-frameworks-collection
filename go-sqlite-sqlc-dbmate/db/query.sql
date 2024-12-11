-- name: GetDistinctCityOrderByCity :many
SELECT DISTINCT city
FROM customers
ORDER BY city;

-- name: GetDistinctCityAndCustomerOrderByCountry :many
SELECT DISTINCT
  city,
  country
FROM
  customers
ORDER BY
  country;

-- name: GetDistinctCompany :many
SELECT DISTINCT
  company
FROM
  customers;