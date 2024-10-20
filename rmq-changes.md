# 3-1
https://www.boot.dev/lessons/38ff4121-f511-467d-8ce8-3469a6c31120
- Add exchange
  - Name: peril_direct
  - Type: direct

# 3-4
https://www.boot.dev/lessons/dacf0de0-ef47-4343-80cc-7eca3e1c4a4e
- Temporarily add queue
  - Name: pause_test
  - Durability: Durable
  - Add binding
    - From exchange: peril_direct
    - Routing key: pause
- Delete this queue