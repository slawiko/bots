# Russian-Belarusian translation bot

Bot simplifies transition to belarussian language. Don't know how to say %WORD% in belarussian? Just ask this %WORD%, [Жэўжык](https://be.wikipedia.org/wiki/%D0%96%D1%8D%D1%9E%D0%B6%D1%8B%D0%BA_(%D0%BF%D0%B5%D1%80%D1%81%D0%B0%D0%BD%D0%B0%D0%B6)) will answer you with the translation from https://www.skarnik.by/ .

This bot could be added in group. To use it just write `як будзе %WORD%` and Жэўжык will reply with translation.

### Run

```bash
cd bot
go build
bot "TELEGRAM_BOT_API_TOKEN"
```

or using docker

```bash
docker build -t jeujik_bot .
docker run -i jeujik_bot "TELEGRAM_BOT_API_TOKEN"
```

### Attention
To allow your self-hosted bot read group message you have to *disable* [privacy mode](https://core.telegram.org/bots#privacy-mode) **before** adding bot to the group. See how to do this: https://teleme.io/articles/group_privacy_mode_of_telegram_bots
