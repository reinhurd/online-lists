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

## How to create valid yandex app for working with disk
1. You need create app strongly by the link https://oauth.yandex.ru/client/new - only there you can add custom right from the official docs https://yandex.ru/dev/disk-api/doc/ru/concepts/quickstart. It's important, because you can't find valid creation link in the docs.
2. Add rights to the app, and authorize through it
3. Copy the token and put it in secret.env