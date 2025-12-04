# a crappy key-value store server

mildly inspired by [Skate](https://github.com/charmbracelet/skate).

Why not just use [Skate](https://github.com/charmbracelet/skate)? This project has a self-hosted server

<sub>(unlike Skate, which decided to just go local-only when sun-setting the cloud service)</sub>

---

# TODO (feature list)
- [x] http server
- [x] database (using the experimental [gomn](https://github.com/Supraboy981322/gomn)-as-a-binary)
- [x] fetching value using key
- [x] creating a key-value pair
- [ ] database compression support
- [x] deleting key-value pair
- [ ] fetching full database
  - [ ] standard gomn (maybe, not sure yet)
  - [x] binary gomn
  - [x] key-value pair form
- [ ] client
  - [ ] prototype
  - [ ] stable
  - [ ] downloading client from server
  - [ ] pipe stdin
  - [ ] changing server conf from client
  - [ ] file input
  - [ ] input dir (streamed tarball to server)
  - [ ] flag to extract tarball
  - [ ] flag to compress data (specify compression libs)
