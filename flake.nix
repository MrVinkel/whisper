{
  description = "Whisper secrets to your development environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};

        lastTag = "v0.0.1";

        revision =
          if (self ? shortRev)
          then "${self.shortRev}"
          else "${self.dirtyShortRev or "dirty"}";

        # Add the commit to the version string for flake builds
        version = "${lastTag}";

        # Run `make vendor-hash` to update the vendor-hash
        vendorHash =
          if builtins.pathExists ./vendor-hash
          then builtins.readFile ./vendor-hash
          else "";

        buildGoModule = pkgs.buildGo123Module;

      in
      {
        inherit self;
        packages.default = buildGoModule {
          pname = "whisper";
          inherit version vendorHash;

          src = ./.;

          subpackage = [ ./cmd/whisper ];

          ldflags = [
            "-s"
            "-w"
            "-X github.com/mrvinkel/whisper/cmd/whisper/cmd.Version=${version}"
          ];

          # Disable tests if they require network access or are integration tests
          doCheck = false;

          nativeBuildInputs = [ pkgs.installShellFiles ];

          meta = with pkgs.lib; {
            description = "Whisper secrets to your development environment";
            homepage = "https://github.com/mrvinkel/whisper";
            license = licenses.unlicense;
            maintainers = with maintainers; [ mrvinkel ];
          };
        };
      }
    );
}