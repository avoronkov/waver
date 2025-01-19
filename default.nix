let
  pkgs = import <nixpkgs> {};
in
pkgs.buildGoModule rec {
	pname = "avoronkov-waver";

	version = "2.8.0";

	src = ./.;

	vendorHash = "sha256-Qz/rOIbpKoDglVYlZl9tcR1hnE18QIfgywgxhAJitSk=";

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
