let
  pkgs = import <nixpkgs> {};
in
pkgs.buildGoModule rec {
	pname = "avoronkov-waver";

	version = "1.2.0";

	src = ./.;

	vendorHash = "sha256-tdrdZA8Uwos1Lfxr+0SubtB23hsAG7LmFLxom5STijk=";

	subPackages = [ "cmd/waver" ];

	meta = {
		description = "Waver. Audio programming language and interpreter";
	};

	buildInputs = with pkgs; [
		alsa-lib
	];

	nativeBuildInputs = with pkgs; [
		pkg-config
	];
}
