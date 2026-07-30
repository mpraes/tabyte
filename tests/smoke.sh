#!/usr/bin/env bash
set -euo pipefail

BASE="${BASE_URL:-http://127.0.0.1:8787/api/v1}"

echo "== health =="
curl -sf "$BASE/health" | grep -q '"status":"ok"'

echo "== info =="
curl -sf "$BASE/info" | grep -q '"app":"tabyte"'

echo "== create =="
RESP=$(curl -sf -X POST "$BASE/analysis-sessions" \
  -H 'Content-Type: application/json' \
  -d '{"engine":"postgres","source_name":"a.sql","ddl_text":"CREATE TABLE a (id INT);"}')
echo "$RESP"
ID=$(echo "$RESP" | sed -n 's/.*"id":"\([^"]*\)".*/\1/p')
test -n "$ID"

echo "== get $ID =="
GET_RESP=$(curl -sf "$BASE/analysis-sessions/$ID")
echo "$GET_RESP" | grep -q "$ID"
echo "$GET_RESP" | grep -q '"name":"a"'

echo "== list =="
LIST_RESP=$(curl -sf "$BASE/analysis-sessions")
echo "$LIST_RESP" | grep -q "$ID"
echo "$LIST_RESP" | grep -q '"table_count":1'

echo "== reject bad engine =="
CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE/analysis-sessions" \
  -H 'Content-Type: application/json' \
  -d '{"engine":"mysql","ddl_text":"CREATE TABLE t (id INT);"}')
test "$CODE" = "400"

echo "== delete =="
CODE=$(curl -s -o /dev/null -w "%{http_code}" -X DELETE "$BASE/analysis-sessions/$ID")
test "$CODE" = "204"

echo "== get after delete =="
CODE=$(curl -s -o /dev/null -w "%{http_code}" "$BASE/analysis-sessions/$ID")
test "$CODE" = "404"

echo "== parse two tables =="
RESP2=$(curl -sf -X POST "$BASE/analysis-sessions" \
  -H 'Content-Type: application/json' \
  -d '{"engine":"postgres","source_name":"two.sql","ddl_text":"CREATE TABLE users (id INT); CREATE TABLE orders (id INT);"}')
echo "$RESP2"
echo "$RESP2" | grep -q '"name":"users"'
echo "$RESP2" | grep -q '"name":"orders"'
ID2=$(echo "$RESP2" | sed -n 's/.*"id":"\([^"]*\)".*/\1/p')

# clean up second session so list stays tidy
curl -s -o /dev/null -X DELETE "$BASE/analysis-sessions/$ID2"

echo "OK: all smoke checks passed"