ALTER DATABASE rinha SET log_error_verbosity to 'TERSE';

CREATE TABLE clientes (
  id SERIAL PRIMARY KEY,
  nome VARCHAR(255) NOT NULL,
  limite INTEGER NOT NULL,
  saldo INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE transacoes (
  id SERIAL PRIMARY KEY,
  cliente_id INTEGER NOT NULL,
  valor INTEGER NOT NULL,
  tipo CHAR NOT NULL,
  descricao VARCHAR(10) NOT NULL,
  realizada_em TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_cliente FOREIGN KEY (cliente_id) REFERENCES clientes (id)
);

CREATE INDEX idx_realizada_em ON transacoes(realizada_em);
