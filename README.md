# Static Spoof

`static-spoof` is an excerise for showing how to serve file at arbitrary paths. In this example there are two image that have **different** paths that the backend maps to the **same** image. This is done by setting a mapping different paths to the same image path and then writing the same image bytes to the frontend.


### Running
1. activate the dev environment by running [nix-shell](https://nixos.org/manual/nix/stable/command-ref/nix-shell.html)
2. run the main program with `go run main.go <path to the image to serve>`
3. navigate to `localhost:8080` in browser
