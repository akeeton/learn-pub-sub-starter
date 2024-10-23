TODO: Turn off gopls settings if they're bad

<!-- # 3-1
https://www.boot.dev/lessons/38ff4121-f511-467d-8ce8-3469a6c31120
- Add exchange (obsoleted by my code that creates it)
  - Name: peril_direct
  - Type: direct

# 3-10
https://www.boot.dev/lessons/98b655ca-255b-4a09-b066-d86f54a56495
- Add exchange (obsoleted by my code that creates it)
  - Name: peril_topic
  - Type: topic -->

# 5-1
https://www.boot.dev/lessons/1a81f164-aa4f-4f66-a481-fac6a99e1eba
- Add exchange
  - Name: peril_dlx
  - Type: fanout
- Add queue
  - Name: peril_dlq
  - Binding: peril_dlx
  - Key: none

<!-- # 6-2
https://www.boot.dev/lessons/275899a4-5787-4297-8507-5685789a5517
- Add queue (obsolete)
  - Name: game_logs
  - Binding: peril_topic
  - Key: game_logs.* -->