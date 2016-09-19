
Move funds

curl --data "amount=42.2&from=bank1&to=bank2" -i -X POST "localhost:8787/funds/move"

Get balance

curl -i -v "localhost:8787/balance/bank1"

Get stats

curl -i -v "localhost:8787/stats"

Checksum

curl -i -v "localhost:8787/checksum"
