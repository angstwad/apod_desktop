# APOD Desktop Setter for OS X

This was my first attempt at some Go.  It's probably suckariffic.

I wanted to make sure my desktop was rotating daily and being set to the [NASA APOD](http://apod.nasa.gov/‎), because, you know, why not?  

### Contents

1. Some Go.
2. `launchd` plist so this runs at boot and every hour.

### Configuration?

I built the package and copied the binary into `/usr/local/bin`, and copied the `ApodDesktop.plist` to `~Library/LaunchAgents/com.angstwad.ApodDesktop.plist`.  

That's actually a pretty stupid way to install this app since my `$GOPATH/bin` should probably be in my $PATH, but it's not.  I could have fixed that, but I didn't.

Enjoy.
