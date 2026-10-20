# Resource and Memory Budget — go-anilist

`go-anilist` is designed to operate as a core dependency for `aniSync`,
which runs under a strict memory budget (≤ 256 MiB idle, ≤ 512 MiB peak).
Therefore, the library must remain highly efficient and predictable.

---

## 1. Binary Footprint

- **No Bloated GraphQL Clients**: Avoid third-party libraries that drag entire schemas,
  AST parsers, or runtime introspection.
- **Prefer Standard Library**: Keep direct dependencies to the absolute minimum (only Go stdlib).
- **Measured Impact [RE]**: Linking `go-anilist` into an application binary adds less than
  ~50 KiB of compiled code and < 30 KiB stripped (`-s -w`).

---

## 2. Memory Allocation and Streaming

- **Streaming Deserialization**: Use `json.NewDecoder(resp.Body)` directly on the response stream,
  always ensuring `defer resp.Body.Close()`.
- **Avoid Duplicating Buffers**: Do not buffer full HTTP responses in memory except when
  needed for error diagnostic messages.
- **Compact Types**: Design structs containing only the fields requested in the GraphQL query.
  By requesting only necessary fields, payload size across network and memory is minimized.

---

## 3. Pagination and List Limits

- **AniList Pagination Limits**: AniList caps the `perPage` parameter to a maximum of 50 items **[DOC-ANILIST]**.
- **No Forced In-Memory Buffering**: The library must not force loading all user lists into memory at once.
- **Paginated Iteration**: Provide paginated access that allows consumers to process page by page.