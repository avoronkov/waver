let
  pkgs = import <nixpkgs> {};
in
pkgs.buildGoModule rec {
	pname = "avoronkov-waver";

	version = "2.8.0";

	src = ./.;

	vendorHash = "sha256-ZQiye9a7pa8rOhqO/kVHkdOoYZwtiV15Yp2jo/P4hBk=";

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
