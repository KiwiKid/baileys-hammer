baileys-hammer is an application for tracking seasonal fines for a sports team, it can be self-hosted via docker and includes configuration for a fly.io deployment



## Development
```
(using vscode by default)
nix develop .

# pre-commit/deploy
nix develop .#devShells.build

# deploy
nix develop .#devShells.deploy



# or via docker
$ docker build -t baileys-hammer .
$ docker run -p 8080:8080 baileys-hammer

# or build docker via nix flake:
nix develop .#devShells.dockerBuild

fly deploy
```





