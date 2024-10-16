let
  pkgs = import <nixpkgs> {};
in
pkgs.buildGoModule rec {
	pname = "avoronkov-waver";

	version = "2.8.0";

	src = ./.;

	vendorHash = "sha256-b5ziDhECy3qMuVKdbxXu9AU7D0amcVY0RE5nnw7oYvE=";

	subPackages = [ "cmd/waver" ];

	meta = {
		description = "Waver. Audio programming language and interpreter";
	};

	buildInputs = with pkgs; [
		alsa-lib
		pulseaudio
	];

	nativeBuildInputs = with pkgs; [
		pkg-config
	];
}
