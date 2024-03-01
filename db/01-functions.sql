CREATE OR REPLACE FUNCTION process_transaction(
    cliente_id INT,
    descricao TEXT,
    valor INT
)
RETURNS TABLE (new_saldo INT, client_limite INT)
LANGUAGE plpgsql
AS $$
BEGIN 
    SELECT saldo + valor, limite INTO new_saldo, client_limite FROM clientes WHERE id = cliente_id;

    IF new_saldo < client_limite THEN
        RAISE EXCEPTION 'New saldo %s is below the limit %s', new_saldo, client_limite;
    END IF;

    UPDATE clientes SET saldo = new_saldo WHERE id = cliente_id;
    INSERT INTO transacoes (cliente_id, descricao, valor) VALUES (cliente_id, descricao, valor);

    RETURN QUERY SELECT new_saldo, client_limite;
END
$$;
