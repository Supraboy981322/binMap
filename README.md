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
- [x] fetching full database
  - [x] standard gomn (maybe, not sure yet)
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

---

# Usage

>[!NOTE]
>The client is not-yet written, so this shows `curl` usage for now

<sub>the client will be much simpler and less of a pain</sub> 

- Actions:
  - `set`:
    Create a key-value pair
  - `get`:
    Retrieve a value
  - `del`:
    Delete a value
  - `db`:
    Download the full database

- For pretty much all actions, use you can use the following args
  - Key (header): 
    - `-H "k:your key name"`
    - `-H "key:your key name"`
  - Value (header):
    - `-H "v:your value"`
    - `-H "val:your value"`
    - `-H "value:your value"`

- Create a key-value pair (replace `[::1]:4780` with your server address)
  - Using the header for the value
    ```sh
    curl [::1]:4780/set -H "k:foo" -H "v:bar"
    ```
  - Using the request body for the value
    ```sh
    curl [::1]:4780/set -H "k:foo" -d "bar"
    ```
  - Piping from stdin (replace `tar -cf - *` with your command)
    ```sh
    tar -cf * | curl [::1]:4780/set -H "k:home directory" --binary-data @-
    ```
  - Sending a file (replace `image.png` with your file)
    ```sh
    curl [::1]:4780/set -H "k:picture" --data-binary "@image.png"
    ```

- Get a value (replace `[::1]:4780` with your server address)
  - Using the header
    ```sh
    curl [::1]:4780/get -H "k:foo
    ```
  - Using the request body
    ```sh
    curl [::1]:4780/get -d "foo"
    ```
  - Saving to a file
    ```sh
    curl [::1]:4780/get -o home.tar -H "k:home directory"
    ```
  - Saving to a file using stdout
    ```sh
    curl [::1]:4780/get -H "k:picture" > image.png
    ```

- Delete a value (replace `[::1]:4780` with your server address)
  - Using the header
    ```sh
    curl [::1]:4780/del -H "k:foo"
    ```
  - Using the request body
    ```sh
    curl [::1]:4780/del -d "foo"
    ```

- Downloading the database (replace `[::1]:4780` with your address)
  - The following args arg valid to specify format (replace `[type]`)
    - `-H "t:[type]"`
    - `-H "typ:[type"`
    - `-H "type:[type"`
    - `-d "[type]"`

  - As gaas ([gomn](https://github.com/Supraboy981322/gomn)-as-a-binary)
    - Since `binMap` stores it's data base in gaas by default, the type doesn't need to be specified if it's the format you want:
      ```sh
      curl [::1]:4780/db -o db.gaas
      ```
    - Or, if you do want to include the type, the following are valid:
      - `b`
      - `r`
      - `bin`
      - `raw`
      - `gaas`
      - `binary`

  - As standard gomn
    ```sh
    curl [::1]:4780/db -H "t:gomn" -o db.gomn
    ```
    - The following types are valid:
      - `g`
      - `s`
      - `std`
      - `gomn`
      - `standard`

  - As a basic key-value pair
    ```sh
    curl [::1]:4780/db -H "t:k-v" -o db.txt
    ```
    - The following types are valid:
      - `t`
      - `p`
      - `kv`
      - `k-v`
      - `text`
      - `pair`
      - `pairs`
      - `key_val`
      - `key-val`
      - `key value`
      - `key-value`

---

# Notes

- When piping a file to the server using `curl`, it's best to use `--binary-data` instead of `-d` or `--data`, because `curl` strips certain data from stdin before sending it to the server.
