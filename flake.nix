{
  description = "A flake for building, running, and deploying a Go program with devShells";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-23.05";
    nixpkgs-unstable.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = { self, nixpkgs, nixpkgs-unstable }: 
    let 
      system = "x86_64-linux";
      pkgs = import nixpkgs { inherit system; };
      unstablePkgs = import nixpkgs-unstable { inherit system; };
      go = unstablePkgs.go_1_21;
      templ = unstablePkgs.templ;
      air = unstablePkgs.air;
      docker = unstablePkgs.docker;
      flyctl = unstablePkgs.flyctl;
    in
    {
      devShells = {
        build = pkgs.mkShell {
          buildInputs = [
            go
            templ
            flyctl
          ];
          shellHook = ''
            echo "Building the Go project..."
            export DATABASE_URL=./tmp/data/dev.db
            templ generate && go build -o ./tmp/main .
          '';
        };

        deploy = pkgs.mkShell {
          buildInputs = [
            go
            templ
            flyctl
          ];
          shellHook = ''
            echo "Building the Go project..."
            export DATABASE_URL=./tmp/data/dev.db
            templ generate
            echo "Deploying the application..."
            fly deploy
          '';
        };

        tdeploy = pkgs.mkShell {
          buildInputs = [
            go
            templ
            flyctl
          ];
          shellHook = ''
            echo "Building the Go project..."
            export DATABASE_URL=./tmp/data/dev.db
            templ generate
            echo "Deploying the application..."
            fly deploy -c fly-dev.toml
          '';
        };

        dev = pkgs.mkShell {
          buildInputs = [ 
            air
            templ
            go
            pkgs.tmux
            flyctl
          ];
          shellHook = ''
            export DATABASE_URL=./tmp/data/devTEST.db
            export PASS=pass
            code .
            tmux kill-session -t devSession
            tmux new-session -d -s devSession \; \
              split-window -h \; \
              send-keys -t 0 'templ generate --watch' C-m \; \
              send-keys -t 1 'air' C-m \; \
            attach-session -t devSession
          '';
          shellExit = ''
            tmux kill-session -t devSession
          '';
        };

        shell = pkgs.mkShell {
          buildInputs = [ 
            air
            templ
            go
            pkgs.tmux
            flyctl
          ];
          shellHook = ''
            export DATABASE_URL=./tmp/data/dev.db
            export PASS=pass
          '';
        };

        dockerBuild = pkgs.mkShell {
          buildInputs = [ 
            docker
            flyctl
          ];
          shellHook = ''
           if [ -f .env ]; then
              export $(cat .env | xargs)
            fi
            echo "Building Docker image..."
            docker build -t baileys-hammer . --load
            echo "Built Docker image..."
            docker images
            docker run -p 8080:8080 \
              -e DATABASE_URL="$DATABASE_URL" \
              -e PASS="$PASS" \
              -v "$(pwd)/tmp/data:/tmp/data" \
              baileys-hammer
          '';
        };


      };
      defaultPackage.x86_64-linux = self.devShells.dev;

    };
}
