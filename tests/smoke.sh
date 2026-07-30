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

echo "$RESP" | grep -q '"name":"a"'
echo "$RESP" | grep -q '"normalized_type":"int"'
echo "$RESP" | grep -q '"name":"id"'
echo "$RESP" | grep -qi 'INT'
echo "$RESP" | grep -q '"estimated_bytes":4'
echo "$RESP" | grep -q '"estimated_row_bytes":28'   # 23 + 1 + 4
echo "$RESP" | grep -q '"assumed_row_count":1000'
echo "$RESP" | grep -q '"estimated_table_bytes":28000'   # 28 * 1000
echo "$RESP" | grep -q '"estimated_total_bytes":28000'
echo "$RESP" | grep -q '"column_payload_bytes":4'
echo "$RESP" | grep -q '"row_header_bytes":23'
echo "$RESP" | grep -q '"null_bitmap_bytes":1'

echo "== get $ID =="
GET_RESP=$(curl -sf "$BASE/analysis-sessions/$ID")
echo "$GET_RESP" | grep -q "$ID"
echo "$GET_RESP" | grep -q '"name":"a"'
echo "$GET_RESP" | grep -q '"column_count":1'
echo "$GET_RESP" | grep -q '"estimated_row_bytes":28'
echo "$GET_RESP" | grep -q '"assumed_row_count":1000'
echo "$GET_RESP" | grep -q '"estimated_table_bytes":28000'
echo "$GET_RESP" | grep -q '"name":"id"'
echo "$GET_RESP" | grep -q '"estimated_bytes":4'

echo "== patch row count =="
PATCH_RESP=$(curl -sf -X PATCH "$BASE/analysis-sessions/$ID/tables/a" \
  -H 'Content-Type: application/json' \
  -d '{"assumed_row_count":5000}')
echo "$PATCH_RESP"
echo "$PATCH_RESP" | grep -q '"assumed_row_count":5000'
echo "$PATCH_RESP" | grep -q '"estimated_table_bytes":140000'
echo "$PATCH_RESP" | grep -q '"estimated_total_bytes":140000'
GET_RESP2=$(curl -sf "$BASE/analysis-sessions/$ID")
echo "$GET_RESP2" | grep -q '"assumed_row_count":5000'
echo "$GET_RESP2" | grep -q '"estimated_total_bytes":140000'

echo "== list =="
LIST_RESP=$(curl -sf "$BASE/analysis-sessions")
echo "$LIST_RESP" | grep -q "$ID"
echo "$LIST_RESP" | grep -q '"table_count":1'

echo "== reject bad engine =="
CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE/analysis-sessions" \
  -H 'Content-Type: application/json' \
  -d '{"engine":"mysql","ddl_text":"CREATE TABLE t (id INT);"}')
test "$CODE" = "400"

echo "== growth projection =="
GROWTH_RESP=$(curl -sf -X PATCH "$BASE/analysis-sessions/$ID/tables/a/growth" \
  -H 'Content-Type: application/json' \
  -d '{"rows_per_period":100,"period":"day","horizon":30}')
echo "$GROWTH_RESP"
echo "$GROWTH_RESP" | grep -q '"projected_row_count":8000'
echo "$GROWTH_RESP" | grep -q '"projected_table_bytes":224000'
echo "$GROWTH_RESP" | grep -q '"projected_total_bytes":224000'

echo "== delete =="
CODE=$(curl -s -o /dev/null -w "%{http_code}" -X DELETE "$BASE/analysis-sessions/$ID")
test "$CODE" = "204"

echo "== get after delete =="
CODE=$(curl -s -o /dev/null -w "%{http_code}" "$BASE/analysis-sessions/$ID")
test "$CODE" = "404"

echo "== normalize varchar =="
RESP3=$(curl -sf -X POST "$BASE/analysis-sessions" \
  -H 'Content-Type: application/json' \
  -d '{"engine":"postgres","source_name":"v.sql","ddl_text":"CREATE TABLE t (name VARCHAR(100));"}')
echo "$RESP3"
echo "$RESP3" | grep -q '"normalized_type":"varchar"'
echo "$RESP3" | grep -q '"length":100'
echo "$RESP3" | grep -q '"assumed_avg_length":50'
echo "$RESP3" | grep -q '"estimated_bytes":51'
ID3=$(echo "$RESP3" | sed -n 's/.*"id":"\([^"]*\)".*/\1/p')
test -n "$ID3"
curl -s -o /dev/null -X DELETE "$BASE/analysis-sessions/$ID3"

echo "== normalize sqlserver nvarchar =="
RESP4=$(curl -sf -X POST "$BASE/analysis-sessions" \
  -H 'Content-Type: application/json' \
  -d '{"engine":"sqlserver","source_name":"s.sql","ddl_text":"CREATE TABLE t (name NVARCHAR(50));"}')
echo "$RESP4"
echo "$RESP4" | grep -q '"normalized_type":"nvarchar"'
echo "$RESP4" | grep -q '"length":50'
echo "$RESP4" | grep -q '"assumed_avg_length":50'   # 50 < 64 → use full length
echo "$RESP4" | grep -q '"estimated_bytes":102'      # 50*2+2
ID4=$(echo "$RESP4" | sed -n 's/.*"id":"\([^"]*\)".*/\1/p')
test -n "$ID4"
curl -s -o /dev/null -X DELETE "$BASE/analysis-sessions/$ID4"

echo "== parse two tables =="
RESP2=$(curl -sf -X POST "$BASE/analysis-sessions" \
  -H 'Content-Type: application/json' \
  -d '{"engine":"postgres","source_name":"two.sql","ddl_text":"CREATE TABLE users (id INT); CREATE TABLE orders (id INT);"}')
echo "$RESP2"
echo "$RESP2" | grep -q '"name":"users"'
echo "$RESP2" | grep -q '"name":"orders"'
ID2=$(echo "$RESP2" | sed -n 's/.*"id":"\([^"]*\)".*/\1/p')
echo "$RESP2" | grep -q '"estimated_total_bytes":56000'
GET_RESP2=$(curl -sf "$BASE/analysis-sessions/$ID2")
echo "$GET_RESP2" | grep -q '"assumed_row_count":1000'
echo "$GET_RESP2" | grep -q '"estimated_total_bytes":56000'

# clean up second session so list stays tidy
curl -s -o /dev/null -X DELETE "$BASE/analysis-sessions/$ID2"

echo "== structural warnings =="
RESP_W=$(curl -sf -X POST "$BASE/analysis-sessions" \
  -H 'Content-Type: application/json' \
  -d '{"engine":"postgres","source_name":"w.sql","ddl_text":"CREATE TABLE t (name VARCHAR(500));"}')
echo "$RESP_W"
echo "$RESP_W" | grep -q '"code":"WIDE_VARCHAR"'
echo "$RESP_W" | grep -q '"warning_count":1'
echo "$RESP_W" | grep -q '"code":"WIDE_ROW"'
echo "$RESP_W" | grep -q '"signal_count":1'

echo "== parse indexes =="
RESP_IX=$(curl -sf -X POST "$BASE/analysis-sessions" \
  -H 'Content-Type: application/json' \
  -d '{"engine":"postgres","source_name":"ix.sql","ddl_text":"CREATE TABLE users (id INT PRIMARY KEY, email VARCHAR(100)); CREATE INDEX idx_users_email ON users (email);"}')
echo "$RESP_IX"
echo "$RESP_IX" | grep -q '"kind":"primary_key"'
echo "$RESP_IX" | grep -q '"kind":"index"'
echo "$RESP_IX" | grep -q '"name":"idx_users_email"'
ID_IX=$(echo "$RESP_IX" | sed -n 's/.*"id":"\([^"]*\)".*/\1/p')
test -n "$ID_IX"
curl -s -o /dev/null -X DELETE "$BASE/analysis-sessions/$ID_IX"
echo "== ui =="
curl -sf "http://127.0.0.1:8787/" | grep -qi 'Tabyte'
echo "OK: all smoke checks passed"