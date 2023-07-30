<div align="center">
  <h3>Rotabot</h3>
  <p><b>Making handling rotas dead simple</b></p>
  <a href="https://github.com/rotabot-io/rotabot/blob/main/LICENSE"><img src="https://user-images.githubusercontent.com/4412200/201544613-a7197bc4-8b61-4fc5-bf09-68ee10133fd7.svg"/></a>
  <img src="https://github.com/rotabot-io/rotabot/actions/workflows/build.yml/badge.svg"/>
  <br/>
  <b><a target="_blank" href="https://rotabot.io" >Learn more »</a></b>
</div>
<br/>

## Development

Rotabot uses [testcontainers](https://testcontainers.com/) to ensure that the tests have everything they need to run. Make sure you have docker running locally.

In order to run the tests you need to do the following:

```shell
make install # You only need to do this once

make test
```

If you want to run the app locally you can do so by running `make dev`. This will start a PG database for you.

Once you have the app running you can connect to slack's API by using [ngrok](https://ngrok.com/)

## License

Rotabot is licensed under ELv2. Please see the [LICENSE](https://github.com/rotabot-io/rotabot/blob/main/LICENSE) file
for additional information. If you have any licensing questions please email me@kevinrobayna.com.
