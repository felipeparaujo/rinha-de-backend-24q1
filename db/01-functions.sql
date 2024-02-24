CREATE OR REPLACE FUNCTION process_transaction(
    cliente_id INT,
    tipo CHAR(1),
    descricao TEXT,
    valor INT
)
RETURNS TABLE (new_saldo INT, client_limite INT)
LANGUAGE plpgsql
AS $$
BEGIN

  -- Lock the clientes table row
  -- PERFORM * FROM clientes WHERE id = cliente_id FOR UPDATE;

  -- valor must have the correct sign depending on debit or credit
  IF tipo = 'c' THEN
    SELECT saldo + valor INTO new_saldo FROM clientes WHERE id = cliente_id;
  ELSIF tipo = 'd' THEN
      SELECT saldo - valor INTO new_saldo FROM clientes WHERE id = cliente_id;
  END IF;

  -- Check if the new saldo is smaller than the limite
  SELECT limite INTO client_limite FROM clientes WHERE id = cliente_id;
  IF new_saldo < -1 * client_limite THEN
      RETURN QUERY SELECT 0, -1;
      RETURN;
  END IF;

  -- Update the saldo in the clientes table
  UPDATE clientes SET saldo = new_saldo WHERE id = cliente_id;

  -- Insert a new entry in the transacoes table
  INSERT INTO transacoes (cliente_id, tipo, descricao, valor)
  VALUES (cliente_id, tipo, descricao, valor);

  RETURN QUERY SELECT new_saldo, client_limite;
END
$$;
