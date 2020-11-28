127.0.0.1:6379> JSON.GET test-room
"{\"roomCode\":\"test-room\",\"status\":\"waiting\",\"firstTeam\":\"red\",\"turn\":\"red\",\"words\":[{\"text\":\"silver surfer\",\"identity\":\"assassin\"},{\"text\":\"hulk\",\"identity\":\"spectator\"}],\"players\":[{\"name\":\"Ron Dyck\",\"role\":\"spymaster\"},{\"name\":\"Ben Ha\",\"role\":\"operative\"}]}"

127.0.0.1:6379> JSON.GET test-room .words
"[{\"text\":\"silver surfer\",\"identity\":\"assassin\"},{\"text\":\"hulk\",\"identity\":\"spectator\"}]"

127.0.0.1:6379> JSON.GET test-room .players
"[{\"name\":\"Ron Dyck\",\"role\":\"spymaster\"},{\"name\":\"Ben Ha\",\"role\":\"operative\"}]"

127.0.0.1:6379> JSON.GET test-room .players[0]
"{\"name\":\"Ron Dyck\",\"role\":\"spymaster\"}"
