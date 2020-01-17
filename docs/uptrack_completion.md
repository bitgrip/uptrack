## uptrack completion

Generates bash completion scripts

### Synopsis



To configure your bash shell to load completions for each session add to your bashrc

* ~/.bashrc or ~/.profile

        . <(uptrack completion)

* On Mac (with bash completion installed from brew)

        uptrack completion > $(brew --prefix)/etc/bash_completion.d/uptrack

* To load completion run

        . <(uptrack completion)

This will only temporaly activate completion in the current session.


```
uptrack completion [flags]
```

### Options

```
  -h, --help   help for completion
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.uptrack.yaml)
      --log-json        if to log using json format
  -v, --verbosity int   verbosity level to use
```

### SEE ALSO

* [uptrack](uptrack.md)	 - track down your uptime

