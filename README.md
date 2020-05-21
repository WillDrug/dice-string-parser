# GO dice string parser
## About
This repository takes a dice string into the `Parse` function and returns a `RollParse` object.

## Dice string
Dice string consists of dice rolls, plus and minus signs. e.g.: `1d100+120-12d10kl1`. All results are stored and returned separately inside the `RollParse` object.
A roll expression is `1[d1[(kh|kl)1]l1]`
Dice string currently supports the following:
* `d`: number of dice sides. Required for all the other flags.
* `kh`: Keep HIGH. Keeps the set number of dice from all the rolls, prioritizing the high values
* `kl`: Keep LOW. Same as keep high, but priritizes lower values
* `l`: Limit. Whatever the roll result is, it cannot go higher than the provided value

*Edge case*: If you put just a number (as in `1d100+120` it is read like `120d1` and no dice rolls will be made)

## TODO
* `e` and `b` flags for *explode* and *botch*
* `l` flag for several rolls, like `(1d100+120)l200`

## Go MOD
Pushed tag v1.0.0