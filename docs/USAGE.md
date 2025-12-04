# binMap --> Usage

>[!NOTE]
>The client is not-yet written, so this shows `curl` usage for now
>  (The client will be much simpler and less of a pain)

## Actions

- `set`:
  Create a key-value pair
- `get`:
  Retrieve a value
- `del`:
  Delete a value
- `db`:
  Download the full database

## General args

  For pretty much all actions, use you can use the following args

- Key (header): 
  - `-H "k:your key name"`
  - `-H "key:your key name"`
- Value (header):
  - `-H "v:your value"`
  - `-H "val:your value"`
  - `-H "value:your value"`

## Create a key-value pair

 (replace `[::1]:4780` with your server address)

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

## Get a value
 
 (replace `[::1]:4780` with your server address)

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

## Delete a pair

 (replace `[::1]:4780` with your server address)

- Using the header
  ```sh
  curl [::1]:4780/del -H "k:foo"
  ```
- Using the request body
  ```sh
  curl [::1]:4780/del -d "foo"
  ```

## Downloading the database

 (replace `[::1]:4780` with your address)

- The following args are valid to specify format (replace `[type]`)
  - `-H "t:[type]"`
  - `-H "typ:[type]"`
  - `-H "type:[type]"`
  - `-d "[type]"`

- As gaas ([gomn](https://github.com/Supraboy981322/gomn)-as-a-binary)

  Since `binMap` stores it's data base in gaas by default, the type doesn't need to be specified if it's the format you want:
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
