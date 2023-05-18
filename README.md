goANS — Automatic Notification System for the Second Life Marketplace written in Go
===================================================================================

## What's this for?

Almost two decades ago, a little-known feature of the [Second Life Marketplace](https://marketplace.secondlife.com) was inherited by Linden Lab when they bought the website from its previous owners: the [Automatic Notification System](https://wiki.secondlife.com/wiki/Direct_Delivery_and_Automatic_Notification_System), a very simple API that contacts a remote web server (chosen by the merchant) to automatically inform them when a new sale has been made via the SL Marketplace.

The most well-known example of this is [Casper](https://casperdns.com/), now a fully-owned subsidiary of Linden Lab, and they're _not_ telling us any of their secrets :-)

## Compiling/Building

Tested under Go 1.20. If you don't have it, nor use it, this is hardly the place to explain everything. Suffice to say that if you have Go properly installed, all you need to do is to type `go install goans` and that should be ok. Alternatively, `go build -o goans` will drop a binary on the current irectory, which you can then place wherever you wish.

## Configuration

Copy `config.ini.sample` to `config.ini`, fill with your details, and you're ready to go.

You will probably need to get the Merchant Salt Code first. Just type your processor URL on the right spot on the Merchant backoffice, and LL will show you yours.

## License

Licensed under a [MIT license](https://gwyneth-llewelyn.mit-license.org/)