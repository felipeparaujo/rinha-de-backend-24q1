#!/bin/bash

export DB_CPU=1.2
export LB_CPU=0.1
export API_CPU=0.1

while (($(bc <<<"$DB_CPU >= 0.12"))); do
  docker compose rm -fsv && docker compose up --build >"stack-$DB_CPU" 2>"stack-error-$DB_CPU" &
  sleep 5
  ./executar-teste-local.sh >"teste-$DB_CPU" 2>"teste-error-$DB_CPU"

  DB_CPU=$(bc <<<"$DB_CPU - 0.09")
  LB_CPU=$(bc <<<"$LB_CPU + 0.03")
  API_CPU=$(bc <<<"$API_CPU + 0.03")
done
