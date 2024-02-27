package api

const Select10LatestTransactionsForUser = `
SELECT
	transaction_value, transaction_type, transaction_description, create_time
FROM
	ledger
WHERE
	client_id = $1 AND is_seed_transaction = false
ORDER BY
  client_transaction_count DESC
LIMIT 10
`

const SelectBalanceAndLimitForUser = `
SELECT
	client_balance, client_limit
FROM
	ledger
WHERE
	client_id = $1
ORDER BY
	client_transaction_count DESC
LIMIT 1`

const CreateTransaction = `
	SELECT out_limit, out_balance, out_updated_row_count FROM create_transaction($1, $2, $3, $4)
`
