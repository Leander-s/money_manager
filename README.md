# Money manager

## Config
A config file is evaluated from ~/.config/money_manager/config. If this file does not exist when
using the app, it is created automatically with the default values. As of right now, it is the only
way to change settings. All settings are set in the config file.

### CLI to change config

    money config <key>=<value>

## Enter

    money enter <amount>

This enters a bank account state into the data.

## Read current budget

    money budget 

This prints how much money you currently have available for spending according to the amounts you 
entered and your configured saving habit. Default is half you earn goes to available money, half to 
savings.

## Read last balance

    money balance

This prints the amount you entered last.

## Reset

    money reset

This resets all the data saved in the background. Read will now print 0. You have to enter your 
starting amount again to start the saving process again.

## Run server
### Go server
The go server is not fully operational yet. To check out the current state read more [here](./go/README.md)
### Zig server
> [!CAUTION]
> This is will likely change drastically in the future
> The zig server is currently not being developed. It is being replaced by a [go server](./go/README.md)

    money run

This starts a money manager server to access from a money manager client.
> [!WARNING]
> There is no client yet but you can use curl to access the server.
### Curl examples
Get current available money:

    curl <ip>:8080/budget/

Enter an amount:

    curl <ip>:8080/balance/ -d 1000
