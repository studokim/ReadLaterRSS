# ReadLaterRSS

Your *read later* list is now an RSS feed!

Just like [Pocket](https://getpocket.com/). But simpler.

## Motivation

There are times when you want to save internet articles for later, but having a separate file or an app to read them feels uncomfortable. **ReadLaterRSS** helps to integrate your list into any RSS reader you already use.

## Usage

### Starting the server

```bash
git clone https://github.com/studokim/ReadLaterRSS.git
cd ReadLaterRSS
go get github.com/studokim/ReadLaterRSS
go build
./ReadLaterRSS --listen <port>
```

Now subscribe to the new feed using your RSS reader: the address is `localhost:port/rss`.

### Adding articles

Go to `localhost:port/add`, paste the url of the article and click `Add!`. The article will be converted to an RSS item. If any error occurs, it will be shown immediately.

All the articles are saved into `history.yml` with the timestamp, so you may restart the server anytime preserving the feed state.
