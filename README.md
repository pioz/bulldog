# Bulldog

__Bulldog__ is an amazing hunting dog that checks for you a list of URLs and
warns you by email if one of them returns a http code that is not 200.

## Installation

Install it yourself as:

    $ go get github.com/pioz/bulldog
    $ cd $GOPATH/src/github.com/pioz/bulldog
    $ go build
    $ ./bulldog -v

## Usage

Bulldog loads a list of URLs and for each of them executes a GET http request.
If the request fails or returns a different http code than 200, it will send you
an email to alert you. When the list is over it starts to sleep for a while and
then restarts the controls.

Bulldog uses a configuration file with various options:

```ini
[time]
# After checking the entire list of URLs sleep for these seconds.
s=60
# After checking the entire list of URLs and at least a check fail sleep for
# these seconds. Usually this time is greater to not warn you continuously.
se=600
# Http request timeout. If the timeout is reached the check is to be considered
# as failed.
t=10
[logging]
# Log file path. If empty logs on stdout.
logfile=/var/log/bulldog.log
# Disables logs.
quiet=false
[email]
# Gmail account. If this is present send email using the gmail smtp server. Use
# -pass flag to specify the gmail account password. If this flag is empty send
# email using `mail` command line program.
gmail=account@gmail.com
# Gmail account password. Only relevant when using -gmail flag.
pass=pa$$w0rd
# When a check fails send an email on this email address. If is empty the email
# alert is disabled.
to=your@email.com
[urls]
# Comma-separated list of URLs to check.
urls=http://google.com,http://twitter.com/
```

To unleash Bulldog run the follow command:

    $ bulldog -config /path/to/config/file

You can pass also the config file options as command arguments:

    $ bulldog -urls http://google.com -to your@email.com -s 10

To make a check on the list only once and then exit:

    $ bulldog -urls http://google.com,http://twitter.com -1

For the complete list of command line arguments:

    $ bulldog -h

### Systemd init start-stop script

[Instruction here](https://github.com/pioz/bulldog/wiki/Systemd-script-on-Debian).

## Contributing

Bug reports and pull requests are welcome on GitHub at https://github.com/pioz/bulldog.

## License

The package is available as open source under the terms of the [GPL License](https://github.com/pioz/bulldog/blob/master/LICENSE).