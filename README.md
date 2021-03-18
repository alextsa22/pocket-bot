# pocket-bot

Description
-----------

Telegram bot that adds links to [Pocket](https://getpocket.com/).

Setup instructions
------------------

Clone the repository.

Then add the `.env` file to root directory. It must contain two variables:

+ `TOKEN` is a token for accessing a telegram bot. You can create a telegram
  bot and get a token from [BotFather](https://t.me/botfather).
+ `CONSUMER_KEY` is the key to access your Pocket application. 
  [Here](https://getpocket.com/developer) you can create an application and get a key.

You also need to set a link to your bot in the 
[config.yml](./configs/config.yml) file for the `bot_url` value.

To run, you need to install [Docker](https://docs.docker.com/).

Operating Instructions
-----------------------

Run application:

```
make build & make up
```
