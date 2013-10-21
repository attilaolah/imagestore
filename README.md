# Google App Engine image store

`gaeimagestore` implements a very simple image storage service on top of [Google
App Engine][1].

[1]: https://developers.google.com/appengine/

* Images are stored in [Blobstore][2].
* Images are served using the [Image Hosting API][3].

[2]: https://developers.google.com/appengine/docs/go/blobstore/
[3]: https://developers.google.com/appengine/docs/go/images/

The idea is that this app can be deployed in a few minutes and be ready for use
to test other apps. A fair amount of images can easily fit in the free quota.

## Usage

* Make a `GET` request to create an upload session:

```bash
URL="$(
    curl http://localhost:8080/ -s -o /dev/null --dump-header - | \
	awk '/^[Ll]ocation:/ { print $2 }' | \
	sed -e 's/\(.*\)\r/\1/g'
)"
```

* Make a `POST` request to upload the images:

```bash
curl "$URL" --form file=@92cfceb39d57d914ed8b14d0e37643de0797ae56.jpg
```

An arbitrary number of images can be uploaded at once by repeating the `--form`
argument.

Only filenames of the form `(?i)^[\da-f]{40}\.jpe?g$` are supported. Filenames
will be converted to the form that matches `^[\da-f]\.jpg$`.

## TODO

* Batch API calls
* Memcache support
