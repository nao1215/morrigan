# pwcrack
pwscore subcommand checks password strength and the guessability of password using the dictionary. It also outputs advice on how to increase password strength.

# Synopsis
```
morrigan pwscore USERNAME [flags]
```

# Examples
```sh
$ ./morrigan pwscore nao
Enter password: 
morrigan:INFO : [Weak password or not] OK
morrigan:WARN : [Length              ] NG (15 characters or more is recommended)
morrigan:WARN : [Not contain name    ] OK
morrigan:INFO : [Contains number     ] OK
morrigan:WARN : [Contains upper&lower] NG (Better to include upper character)
morrigan:INFO : [Contains symbol     ] OK
morrigan:INFO : [WIP] Calculate entropy as a measure of password strength
```