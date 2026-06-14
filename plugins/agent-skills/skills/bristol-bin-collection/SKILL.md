---
name: bristol-bin-collection
description: Look up Bristol bin collection schedule via waste.bristol.gov.uk. Returns next-collection dates for general waste, recycling boxes/bags, food waste, and garden waste. Asks for the user's postcode on first run and stores it. Select the first address returned unless the user specifies one.
---

# Bristol Bin Collection

Two POSTs to `waste.bristol.gov.uk` — no browser, no scripts.

## Postcode

The postcode is read from `~/.config/bristol-bin-collection/postcode` and never hard-coded.

### First run (no stored postcode)

If the file doesn't exist, ask the user:

> "What's your postcode? I'll save it so I don't have to ask again."

Then write it (creating the directory if needed):

```bash
mkdir -p ~/.config/bristol-bin-collection
printf '%s\n' "$USER_POSTCODE" > ~/.config/bristol-bin-collection/postcode
```

### Subsequent runs

```bash
POSTCODE=$(cat ~/.config/bristol-bin-collection/postcode)
```

If the user asks about a one-off different postcode, use that for this run only — don't overwrite the stored value unless they explicitly say "save".

## Steps

### 1. Resolve the address to a UPRN

```bash
curl -fsS -X POST --data-urlencode "postcode=$POSTCODE" https://waste.bristol.gov.uk/choose-an-address \
  | grep -oE '<option value=UPRN[0-9]+ title="[^"]+"' \
  | head -1
```

Each match contains both the UPRN and the address, e.g.:

```
<option value=UPRN000000000000 title="Example Address  Bristol  BS0 0AA"
```

Default: pick the first match. If the user names a specific address, drop `| head -1` and search the full option list instead.

### 2. Fetch the collection schedule

```bash
UPRN="UPRN000000000000"   # from step 1
curl -fsS -X POST \
  --data-urlencode "postcode=$POSTCODE" \
  --data-urlencode "gazId=$UPRN" \
  https://waste.bristol.gov.uk/your-collections \
  | grep -oE 'span id="[a-zA-Z0-9]+NextCollectionDate">[^<]+' \
  | sed -E 's|span id="([a-zA-Z0-9]+)NextCollectionDate">|\1: |'
```

A residential address returns six lines, one per bin type. If `curl` gets HTTP 302 with an empty body, the address isn't on Bristol's domestic collection (commercial site, phone mast, out-of-borough). Ask the user to choose a different address from the step 1 list.

## Bin types

| Key | Friendly name |
| --- | --- |
| `blackWheelieBin180Litres` | Black bin (general waste) |
| `blueRecyclingBag` | Blue recycling bag (plastic / cans) |
| `blackRecyclingBox` | Black recycling box (paper / card) |
| `greenRecyclingBox` | Green recycling box (glass) |
| `brownFoodWasteBin` | Food waste |
| `greenGardenWasteBin` | Garden waste |

## Response format

Concise, weekday + date, using the friendly names:

```
Black bin: Thursday 21 May
Blue recycling bag: Thursday 21 May
Black recycling box: Thursday 21 May
Green recycling box: Thursday 21 May
Food waste: Thursday 21 May
Garden waste: Thursday 28 May
```

- One specific bin asked for? Give one line.
- "All" or "next collection"? List every type.
- Any value missing? Say so and offer to retry with a different address.
