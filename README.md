Anonymous chatting app with no Javascript*

TODO:
- [ ] Figure out if user.websocketReader goroutine are leaking
- [x] improve looks
- [x] avoid overwriting a room if randomly generated id is already in use
- [x] differentiate user messages from peer messages
- [x] add user pseudonym
- [ ] Add loading transition while waiting for websocket to connect
- [ ] Add pages for rejected requests (e.g.: room is full)
- [ ] add online user count
- [ ] close room if left empty for long enough
- [ ] handle new line characters in messages
- [ ] probably more

\* except for that one form reset function call. Also ignore the tailwind config
