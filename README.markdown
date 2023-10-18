# `apron-bus`

Takes you to your aircraft.

> Passengers will be transferred from the airport […] gate to the aircraft using an […] apron bus.
>
> https://en.wikipedia.org/wiki/apron_bus

`apron-bus` selects the `flyX.Y.Z` binary that matches the exact version the Concourse target server requires. It can be aliased as `fly` and transparently invokes the correct binary.

# Installation

```command
$ go install github.com/suhlig/apron-bus@latest
```
