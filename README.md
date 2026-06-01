# baileys-hammer 

## [Live Demo)[https://sports-team-fines.fly.dev/]


## a simple application for tracking games and season fines for a sports team, it can be self-hosted via docker and includes configuration for a fly.io deployment.

![alt text](docs/Kooha-2024-06-27-10-24-42.gif)

Features
- Match, Score & Injury Tracking
- Fine + Player Multi-Select
- Mobile-First for quick entry while discussing fines on the go
- "FineMaster" page for high-level fine/match/season management

Written in go, using templ, a-h/templ and TomSelect (via [hx-tomselect](https://github.com/kiwikid/hx-tomselect))

## Google admin auth

Admin auth can run in two modes:

- Without `GOOGLE_CLIENT_ID`, the existing `PASS`/team admin password flow is used.
- With `GOOGLE_CLIENT_ID`, Google sign-in is required for admin access and password admin login is disabled.

Setup:

1. Create a Google OAuth Web client ID in Google Cloud Console.
2. Add the app origin to the OAuth client, for example `http://localhost:8080` for local dev and the deployed `https://...` origin for production.
3. Set `GOOGLE_CLIENT_ID` to the Web client ID.
4. Set `GOOGLE_CLIENT_SECRET` to a separate Google/admin session secret. Do not reuse `PASS`.
5. Start the app and sign in through the admin entry form.

The first Google admin user created becomes `super-admin`. After that, new admin sign-up is controlled per team by the `Allow admin registration` team flag. Users who sign up through a team get `team-admin` for that team, and users who create a team get `team-admin` for the new team.

To promote an existing Google admin account on startup, set `ADD_SUPER_ADMIN_TO_EMAIL_ON_STARTUP` to that user's email address. This is useful for recovery or first production rollout after users already exist.

## Development
```
(using vscode by default)
nix develop .

# pre-commit/deploy
nix develop .#devShells.build

# deploy
nix develop .#devShells.deploy

# or build docker via nix flake:
nix develop .#devShells.dockerBuild



# or via dockerFines<
$ docker build -t baileys-hammer .
$ docker run -p 8080:8080 baileys-hammer

# or build docker via nix flake:
nix develop .#devShells.dockerBuild

fly deploy
```



TODO: 
- Fix 'ADD' Option sizing (too small)







New Setup

- (Create Fly Volume for sqlitedb)
- create fly-XXX.toml file
 - add DATABASE_URL to litestream.yml
- set fly.io secrets
```
fly secrets -c fly-XXX.toml
    DB_REPLICA_URL
    R2_ACCESS_KEY_ID
    R2_SECRET_ACCESS_KEY
    R2_BUCKET
```
- add 
