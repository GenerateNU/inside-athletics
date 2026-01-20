
## How to use Doppler

# Requirements:

- Doppler CLI: https://docs.doppler.com/docs/install-cli

# How to setup:

First login into doppler using your account. Once you have the doppler cli installed run:
`doppler login` and follow the instructions provided. It should open a window for you to login
and then you need to past the key it gives into that authenticated window. 

Next, cd into `inside-athletics/backend` and run doppler setup. Then select which project (backend) and 
which config (dev)

Now to use this we need to run a command with doppler run --command in order to inject the env variables and use them. This should be handled for you within the makefile but good to know in case you want to run your own command. Here's an example: 

doppler run --comand "atlas migrate apply --env dev"