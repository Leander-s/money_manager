# Money manager

## Enter

    money enter <amount>

This enters a bank account state into the data.

## Read

    money read

This prints how much money you currently have available for spending according to the amounts you 
entered and your configured saving habit. Default is half you earn goes to available money, half to 
savings.

## Reset

    money reset

This resets all the data saved in the background. Read will now print 0. You have to enter your 
starting amount again to start the saving process again.

## Run server
[!CAUTION]
This is will likely change drastically in the future

    money run

This starts a money manager server to access from a money manager client.
[!WARNING]
There is no client yet and the server just echoes whatever one sends to it.
