# baileys-hammer 

## a simple application for tracking games and season fines for a sports team, it can be self-hosted via docker and includes configuration for a fly.io deployment.

![alt text](docs/Kooha-2024-06-27-10-24-42.gif)

Features
- Match, Score & Injury Tracking
- Fine + Player Multi-Select
- Mobile-First for quick entry while discussing fines on the go
- "FineMaster" page for high-level fine/match/season management

Written in go, using templ, a-h/templ and TomSelect (via [hx-tomselect](https://github.com/kiwikid/hx-tomselect))
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





