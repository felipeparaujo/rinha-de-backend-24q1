CREATE TABLE clientes (
  id SERIAL PRIMARY KEY,
  nome VARCHAR(255) NOT NULL,
  limite INTEGER NOT NULL,
  saldo_inicial INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE transacoes (
  id SERIAL PRIMARY KEY,
  cliente_id INTEGER NOT NULL,
  valor VARCHAR(255) NOT NULL,
  tipo CHAR NOT NULL,
  descricao VARCHAR(10) NOT NULL,
  realizada_em TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_cliente FOREIGN KEY(cliente_id) REFERENCES clientes(id)
);

INSERT INTO
  clientes (nome, limite)
VALUES
  ('o barato sai caro', 1000 * 100),
  ('zan corp ltda', 800 * 100),
  ('les cruders', 10000 * 100),
  ('padaria joia de cocaia', 100000 * 100),
  ('kid mais', 5000 * 100);

END;
