# Pattern — Cursor pagination

Tech Stack T3 in practice. Cursors for unbounded collections; offset only for small bounded ones.

## Why cursor over offset

Offset pagination (`?page=42&size=20`) breaks when:

- The collection mutates while a client walks pages (rows shift, duplicates and skips happen).
- The collection is large; `OFFSET 1000000 LIMIT 20` forces the database to scan and discard.

Cursors solve both by encoding a position in the result set, not a count.

## Request shape

```http
GET /v2/customers/acme/invoices?cursor=eyJpZCI6IjEyMyJ9&limit=50
```

- `cursor` — opaque, URL-safe, server-issued. Omitted for the first page.
- `limit` — page size, capped at a server-defined max (default 50, max 200).
- Other filters (e.g. `period=2026-05`) compose with the cursor — the cursor is opaque to the consumer but ties to the same filter set.

## Response shape

```json
{
  "data": [
    { "id": "abc-1", "amount": 1000 },
    { "id": "abc-2", "amount": 2000 }
  ],
  "next_cursor": "eyJpZCI6ImFiYy0yIn0",
  "has_more": true
}
```

- `next_cursor` is present iff `has_more` is true. Omit (or `null`) on the last page.
- `data` is always an array, even when empty.
- No total count — counts are expensive on large tables and rarely needed. If a client truly needs one, expose a separate `/count` endpoint with a cached value.

## Cursor encoding

The cursor is **opaque to the consumer**. Internally it's a base64url-encoded JSON object with the fields needed to resume:

```json
{ "id": "abc-2", "sort": "created_at:2026-05-22T10:00:00Z" }
```

- Always include enough to disambiguate ties — the primary sort key plus the row id.
- Sign or HMAC the cursor if it might be tampered with to access otherwise-restricted data. (Usually not needed when the cursor encodes only what the user could already filter on.)

## SQL pattern (Postgres example)

For a collection sorted by `(created_at DESC, id DESC)`:

```sql
SELECT id, amount, created_at
FROM invoices
WHERE tenant = $1
  AND (created_at, id) < ($2, $3)   -- cursor decoded into (created_at, id)
ORDER BY created_at DESC, id DESC
LIMIT $4 + 1;                       -- one extra row tells us has_more
```

The `LIMIT n + 1` trick: if the extra row comes back, set `has_more = true`, drop it, and build `next_cursor` from the last returned row.

## Common mistakes

- **Sorting by a non-unique field** (e.g. `created_at` alone). Two rows with the same timestamp produce duplicates and skips. Always include a unique tiebreaker.
- **Letting the consumer modify the cursor**. The cursor is opaque. If the API documents its internals, consumers will couple to them.
- **Returning total counts on huge tables**. Don't.
- **Using offset under the hood while pretending to be a cursor**. Same scan cost; defeats the purpose.

## When offset is OK

- Small bounded collections where a "page 3 of 5" UI is genuinely useful.
- Admin tools, internal dashboards.
- Anywhere the maximum collection size is known to be ≤ a few thousand.

## References

- Tech Stack T3 (API design and specs)
- Constitution P11 (non-breaking changes — cursor format can evolve without bumping the API version), P5 (observability — log cursor in/out for debugging)
