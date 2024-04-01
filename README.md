# online-lists
App for maintaining lists through different clients

## You need to have secret.env with params:
- TG_SECRET_KEY: telegram bot secret key
- YANDEX_TOKEN: token for yandex cloud
- YDFILE: path to the file in yandex cloud

## How to create telegram bot
1. Create a bot with BotFather (https://core.telegram.org/bots#botfather) - @BotFather
2. Follow the instructions to create (/newbot and so on)
3. Get the token and put it in secret.env

#### Work in progress

- Add deployment logic
- Add a clients through messengers
- Add a web client
- Add message to tg when app is closing