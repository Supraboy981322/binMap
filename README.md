# a crappy key-value store server

mildly inspired by [Skate](https://github.com/charmbracelet/skate).

Why not just use [Skate](https://github.com/charmbracelet/skate)? This project has a self-hosted server

<sub>(unlike Skate, which decided to just go local-only when sun-setting the cloud service)</sub>

---

# Features/To-do

See the [features doc](docs/TODO.md)

---

# Usage

See the [usage doc](docs/USAGE.md)

---

# Notes

- When piping a file to the server using `curl`, it's best to use `--binary-data` instead of `-d` or `--data`, because `curl` strips certain data from stdin before sending it to the server.
