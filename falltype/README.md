# falltype

`falltype` is a terminal touch-typing trainer for Finnish/Swedish QWERTY letter layout.
Letters fall down fixed keyboard columns and must be hit before they reach the bottom line.

## Run

```bash
go run ./cmd/falltype
```

Optional config:

```bash
go run ./cmd/falltype -config config.json
```

Example config:

```json
{
  "initial_speed_ms": 800,
  "speed_decrease_per_lesson": 50,
  "max_simultaneous_limit": 5
}
```

## Controls

- Type matching letters to remove falling runes.
- `q` or `Esc` exits.

## Lessons

1. `g h`
2. `f g h j`
3. Home-row core with 2 simultaneous letters
4. Full home-row with 3 simultaneous letters
5. Full alphabet with up to 5 simultaneous letters
