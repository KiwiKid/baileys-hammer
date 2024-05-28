baileys-hammer is an application for tracking seasonal fines for a sports team, it can be self-hosted via docker and includes configuration for a fly.io deployment



## Development
```
nix develop .#devShells.dev

# pre-commit/deploy
nix develop .#devShells.build

# deploy
nix develop .#devShells.deploy



# or via docker
$ docker build -t baileys-hammer .
$ docker run -p 8080:8080 baileys-hammer


fly deploy
```





