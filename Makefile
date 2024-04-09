.PHONY: run
run:
	go run . price fetchlatest

.PHONY: test-adapters
test-adapters:
	go test ./internal/adapter/... -count=1

list:
	@sqlite3 -table fund.db "SELECT p.date, p.value, p.net_assets, p.net_assets / p.value 'nums (e4)', substr(f.name, 0, 20) name FROM prices p JOIN funds f ON f.id = p.id WHERE p.date >= date('now', '-7 day') ORDER BY p.id, p.date"
