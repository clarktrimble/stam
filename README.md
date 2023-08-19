
# Stam

Yet another implementation of Joe Stam's famous fluid simulation.

<img width="618" alt="fluidsim-screenshot" src="https://github.com/clarktrimble/stam/assets/5055161/edc7ad02-bf28-4a7a-8990-a5ce63eab2c8">

## Why?

For fun!

Traveling this summer across the vast expanses of the Baltic (har), I wanted
to encourage my 15 year old and game-ish coding seemed like a natural.  And it
worked!  After seeing this coming to life, he wrote from scratch a fluid thing in
Rust with ascii display :)

## Try It Out!

    git clone git@github.com:clarktrimble/stam.git
    cd stam

    go test -count 1 ./... ## minimal
    go run github.com/clarktrimble/stam/cmd/fluidsim

## Travelogue

Translating the C from Stam's paper was a doddle.  Wiring it up to the game engine was
more of a challenge for me:

 - I suppose these algorithms can be unit tested, but I threw up my hands :/
 - Density being modeled here is of dye in an incompressible fluid, aha!
 - Adding density and velocity via one of the swappable vectors still feels awkward.
 - Much fiddling with with viscosity, diffusivity, etc. ...

Ebiten was easy to use.  Want to try out the WASM feature!

C to Golang-isms:

 - Does avoiding multi-dimensional arrays make sense in Golang?
 - Getting there with Golang swap, which is awkward given its pointers.

## Performance

I've tried to avoid optimizing so far.  I'd like to dig into this soon!

## Golang (Anti) Idioms

I dig the Golang community, but I might be a touch rouge with:

  - multi-char variable names
  - named return parameters
  - BDD/DSL testing
  - liberal use of vertical space

All in the name of readability, which of course, tends towards the subjective.

## License

This is free and unencumbered software released into the public domain.

Anyone is free to copy, modify, publish, use, compile, sell, or
distribute this software, either in source code form or as a compiled
binary, for any purpose, commercial or non-commercial, and by any
means.

In jurisdictions that recognize copyright laws, the author or authors
of this software dedicate any and all copyright interest in the
software to the public domain. We make this dedication for the benefit
of the public at large and to the detriment of our heirs and
successors. We intend this dedication to be an overt act of
relinquishment in perpetuity of all present and future rights to this
software under copyright law.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS BE LIABLE FOR ANY CLAIM, DAMAGES OR
OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
OTHER DEALINGS IN THE SOFTWARE.

For more information, please refer to <http://unlicense.org/>

